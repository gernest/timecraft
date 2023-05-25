package timemachine

import (
	"fmt"
	"io"
	"time"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/stealthrocket/timecraft/format/logsegment"
	"github.com/stealthrocket/timecraft/internal/buffer"
	"github.com/stealthrocket/timecraft/internal/stream"
)

var (
	compressedBufferPool   buffer.Pool
	uncompressedBufferPool buffer.Pool
)

// RecordBatch is a read-only batch of records read from a log segment.
//
// The records themselves are compressed and stored separately. To support
// use cases where the user may want to read batch metadata in order to skip
// the processing of records, the record batch is structured such that records
// are read and decompressed lazily.
type RecordBatch struct {
	// Reader for the records data section adjacent to the record batch. When
	// the records are accessed, they're lazily read into the records buffer.
	reader  io.LimitedReader
	records *buffer.Buffer

	batch logsegment.RecordBatch

	// Capture of the log header for the segment that the record batch was read
	// from.
	header *Header

	// When reading records from the batch, this holds the current offset into
	// the records buffer.
	offset uint32
}

// MakeRecordBatch creates a record batch from the specified buffer.
//
// The buffer must live as long as the record batch.
func MakeRecordBatch(header *Header, buffer []byte, reader io.Reader) (rb RecordBatch) {
	rb.Reset(header, buffer, reader)
	return
}

// Reset resets the record batch.
func (b *RecordBatch) Reset(header *Header, buf []byte, reader io.Reader) {
	if b.records != nil {
		recordsBufferPool := &compressedBufferPool
		if b.header.Compression == Uncompressed {
			recordsBufferPool = &uncompressedBufferPool
		}
		buffer.Release(&b.records, recordsBufferPool)
	}
	b.records = nil
	b.header = header
	if len(buf) > 0 {
		b.batch = *logsegment.GetSizePrefixedRootAsRecordBatch(buf, 0)
	} else {
		b.batch = logsegment.RecordBatch{}
	}
	b.offset = 0
	b.reader.R = reader
	b.reader.N = int64(b.RecordsSize())
}

// RecordsSize is the size of the adjacent record data.
func (b *RecordBatch) RecordsSize() (size int) {
	if b.Compression() == Uncompressed {
		size = int(b.UncompressedSize())
	} else {
		size = int(b.CompressedSize())
	}
	return
}

// Compression returns the compression algorithm used to encode the record
// batch data section.
func (b *RecordBatch) Compression() Compression {
	return b.header.Compression
}

// FirstOffset returns the logical offset of the first record in the batch.
func (b *RecordBatch) FirstOffset() int64 {
	return b.batch.FirstOffset()
}

// FirstTimestamp returns the time at which the first record was produced.
func (b *RecordBatch) FirstTimestamp() time.Time {
	return b.header.Process.StartTime.Add(time.Duration(b.batch.FirstTimestamp()))
}

// CompressedSize returns the size of the record batch data section in the log
// segment.
func (b *RecordBatch) CompressedSize() int64 {
	return int64(b.batch.CompressedSize())
}

// UncompressedSize returns the size of the record batch data section after
// being uncompressed.
func (b *RecordBatch) UncompressedSize() int64 {
	return int64(b.batch.UncompressedSize())
}

// NumRecords returns the number of records in the batch.
func (b *RecordBatch) NumRecords() int {
	return int(b.batch.NumRecords())
}

// Read reads records from the batch.
//
// The record values share memory buffers with the record batch, they remain
// valid until the next call to ReadRecordBatch on the parent LogReader.
func (b *RecordBatch) Read(records []Record) (int, error) {
	batch, err := b.readRecords()
	if err != nil {
		return 0, err
	}
	startTime := b.header.Process.StartTime
	functions := b.header.Runtime.Functions
	var n int
	for n = range records {
		if b.offset == uint32(len(batch)) {
			return n, io.EOF
		}
		if b.offset+4 > uint32(len(batch)) {
			return n, fmt.Errorf("cannot read record at offset %d as records buffer is length %d: %w", b.offset, len(batch), io.ErrUnexpectedEOF)
		}
		size := flatbuffers.GetSizePrefix(batch, flatbuffers.UOffsetT(b.offset))
		if b.offset+size < b.offset || b.offset+size > uint32(len(batch)) {
			return n, fmt.Errorf("cannot read record at [%d:%d+%d] as records buffer is length %d: %w", b.offset, b.offset, size, len(batch), io.ErrUnexpectedEOF)
		}
		records[n] = MakeRecord(startTime, functions, batch[b.offset:b.offset+size+4])
		b.offset += size + 4
	}
	return n, nil
}

func (b *RecordBatch) readRecords() ([]byte, error) {
	if b.records != nil {
		return b.records.Data, nil
	}

	recordsBufferPool := &compressedBufferPool
	if b.header.Compression == Uncompressed {
		recordsBufferPool = &uncompressedBufferPool
	}
	recordsBuffer := recordsBufferPool.Get(b.RecordsSize())

	_, err := io.ReadFull(&b.reader, recordsBuffer.Data)
	if err != nil {
		recordsBufferPool.Put(recordsBuffer)
		return nil, err
	}

	if c := checksum(recordsBuffer.Data); c != b.batch.Checksum() {
		return nil, fmt.Errorf("bad record data: expect checksum %#x, got %#x", b.batch.Checksum(), c)
	}

	if b.header.Compression == Uncompressed {
		b.records = recordsBuffer
		return recordsBuffer.Data, nil
	}
	defer recordsBufferPool.Put(recordsBuffer)

	b.records = uncompressedBufferPool.Get(int(b.UncompressedSize()))
	return decompress(b.records.Data, recordsBuffer.Data, b.header.Compression)
}

var (
	_ stream.Reader[Record] = (*RecordBatch)(nil)
)

// RecordBatchBuilder is a builder for record batches.
type RecordBatchBuilder struct {
	builder        *flatbuffers.Builder
	compression    Compression
	firstOffset    int64
	firstTimestamp int64
	recordCount    uint32
	uncompressed   []byte
	compressed     []byte
	records        []byte
	result         []byte
	finished       bool
	concatenated   bool
}

// Reset resets the builder.
func (b *RecordBatchBuilder) Reset(compression Compression, firstOffset int64) {
	if b.builder == nil {
		b.builder = flatbuffers.NewBuilder(buffer.DefaultSize)
	} else {
		b.builder.Reset()
	}
	b.compression = compression
	b.uncompressed = b.uncompressed[:0]
	b.compressed = b.compressed[:0]
	b.result = b.result[:0]
	b.records = nil
	b.firstOffset = firstOffset
	b.firstTimestamp = 0
	b.finished = false
	b.concatenated = false
	b.recordCount = 0
}

// AddRecord adds a record to the batch.
//
// The record is consumed immediately and can be reused safely when the
// call returns.
func (b *RecordBatchBuilder) AddRecord(record *RecordBuilder) {
	if b.finished {
		panic("builder must be reset before records can be added")
	}
	b.uncompressed = append(b.uncompressed, record.Bytes()...)
	if b.recordCount == 0 {
		b.firstTimestamp = record.timestamp
	}
	b.recordCount++
}

// Bytes returns the serialized representation of the record batch.
//
// Since the batch is made up of two components – the batch metadata
// and then the compressed records – additional buffering is required
// here to merge the two together. If efficiency is required, Write
// should be used instead.
func (b *RecordBatchBuilder) Bytes() []byte {
	if !b.finished {
		b.build()
		b.finished = true
	}
	if !b.concatenated {
		b.result = append(b.result, b.builder.FinishedBytes()...)
		b.result = append(b.result, b.records...)
		b.concatenated = true
	}
	return b.result
}

// Write writes the serialized representation of the record batch
// to the specified writer.
func (b *RecordBatchBuilder) Write(w io.Writer) (int, error) {
	if !b.finished {
		b.build()
		b.finished = true
	}
	n1, err := w.Write(b.builder.FinishedBytes())
	if err != nil {
		return n1, err
	}
	n2, err := w.Write(b.records)
	return n1 + n2, err
}

func (b *RecordBatchBuilder) build() {
	if b.builder == nil {
		b.builder = flatbuffers.NewBuilder(buffer.DefaultSize)
	}
	b.builder.Reset()

	b.records = b.uncompressed
	if b.compression != Uncompressed {
		b.compressed = compress(b.compressed[:cap(b.compressed)], b.uncompressed, b.compression)
		b.records = b.compressed
	}

	logsegment.RecordBatchStart(b.builder)
	logsegment.RecordBatchAddFirstOffset(b.builder, b.firstOffset)
	logsegment.RecordBatchAddFirstTimestamp(b.builder, b.firstTimestamp)
	logsegment.RecordBatchAddCompressedSize(b.builder, uint32(len(b.compressed)))
	logsegment.RecordBatchAddUncompressedSize(b.builder, uint32(len(b.uncompressed)))
	logsegment.RecordBatchAddChecksum(b.builder, checksum(b.records))
	logsegment.RecordBatchAddNumRecords(b.builder, b.recordCount)
	logsegment.FinishSizePrefixedRecordBatchBuffer(b.builder, logsegment.RecordBatchEnd(b.builder))
}
