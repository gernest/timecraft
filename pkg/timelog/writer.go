package timelog

import (
	"bytes"
	"fmt"
	"io"
	"time"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/stealthrocket/timecraft/pkg/format/logsegment"
	"github.com/stealthrocket/timecraft/pkg/format/types"
)

type LogWriter struct {
	output    io.Writer
	builder   *flatbuffers.Builder
	buffer    *bytes.Buffer
	functions map[Function]uint32
	startTime time.Time
	records   []flatbuffers.UOffsetT
	offsets   []flatbuffers.UOffsetT
}

func NewLogWriter(output io.Writer) *LogWriter {
	return &LogWriter{
		output:  output,
		builder: flatbuffers.NewBuilder(4096),
		buffer:  bytes.NewBuffer(make([]byte, 0, 4096)),
	}
}

func (w *LogWriter) Reset(output io.Writer) {
	w.output = output
	w.builder.Reset()
	w.buffer.Reset()
	w.startTime = time.Time{}
	w.records = w.records[:0]

	for fn := range w.functions {
		delete(w.functions, fn)
	}
}

func (w *LogWriter) WriteLogHeader(header *LogHeader) error {
	w.builder.Reset()

	processID := w.prependHash(header.Process.ID)
	processImage := w.prependHash(header.Process.Image)
	processArguments := w.prependStringVector(header.Process.Args)
	processEnvironment := w.prependStringVector(header.Process.Environ)

	var parentProcessID flatbuffers.UOffsetT
	if header.Process.ParentProcessID.Digest != "" {
		parentProcessID = w.prependHash(header.Process.ParentProcessID)
	}

	logsegment.ProcessStart(w.builder)
	logsegment.ProcessAddId(w.builder, processID)
	logsegment.ProcessAddImage(w.builder, processImage)
	logsegment.ProcessAddUnixStartTime(w.builder, header.Process.StartTime.UnixNano())
	logsegment.ProcessAddArguments(w.builder, processArguments)
	logsegment.ProcessAddEnvironment(w.builder, processEnvironment)
	if parentProcessID != 0 {
		logsegment.ProcessAddParentProcessId(w.builder, parentProcessID)
		logsegment.ProcessAddParentForkOffset(w.builder, header.Process.ParentForkOffset)
	}
	processOffset := logsegment.ProcessEnd(w.builder)

	functionOffsets := make([][2]flatbuffers.UOffsetT, 0, 64)
	if len(header.Runtime.Functions) <= cap(functionOffsets) {
		functionOffsets = functionOffsets[:len(header.Runtime.Functions)]
	} else {
		functionOffsets = make([][2]flatbuffers.UOffsetT, len(header.Runtime.Functions))
	}

	if w.functions == nil {
		w.functions = make(map[Function]uint32, len(header.Runtime.Functions))
	}

	for i, fn := range header.Runtime.Functions {
		w.functions[fn] = uint32(i)
		functionOffsets[i][0] = w.builder.CreateSharedString(fn.Module)
		functionOffsets[i][1] = w.builder.CreateString(fn.Name)
	}

	functions := w.prependObjectVector(len(header.Runtime.Functions),
		func(i int) flatbuffers.UOffsetT {
			logsegment.FunctionStart(w.builder)
			logsegment.FunctionAddModule(w.builder, functionOffsets[i][0])
			logsegment.FunctionAddName(w.builder, functionOffsets[i][1])
			return logsegment.FunctionEnd(w.builder)
		},
	)

	runtime := w.builder.CreateString(header.Runtime.Runtime)
	version := w.builder.CreateString(header.Runtime.Version)
	logsegment.RuntimeStart(w.builder)
	logsegment.RuntimeAddRuntime(w.builder, runtime)
	logsegment.RuntimeAddVersion(w.builder, version)
	logsegment.RuntimeAddFunctions(w.builder, functions)
	runtimeOffset := logsegment.RuntimeEnd(w.builder)

	logsegment.LogHeaderStart(w.builder)
	logsegment.LogHeaderAddRuntime(w.builder, runtimeOffset)
	logsegment.LogHeaderAddProcess(w.builder, processOffset)
	logsegment.LogHeaderAddSegment(w.builder, header.Segment)
	logsegment.LogHeaderAddCompression(w.builder, types.Compression(header.Compression))
	logHeader := logsegment.LogHeaderEnd(w.builder)

	w.builder.FinishSizePrefixedWithFileIdentifier(logHeader, tl0)

	if _, err := w.output.Write(w.builder.FinishedBytes()); err != nil {
		return err
	}
	w.startTime = header.Process.StartTime
	return nil
}

func (w *LogWriter) WriteRecordBatch(batch []Record) error {
	w.builder.Reset()
	w.buffer.Reset()
	w.records = w.records[:0]

	uncompressedSize := uint32(0)

	for _, record := range batch {
		offset := uncompressedSize
		length := uint32(0)

		for _, access := range record.MemoryAccess {
			uncompressedSize += uint32(len(access.Memory))
			length += uint32(len(access.Memory))

			if _, err := w.buffer.Write(access.Memory); err != nil {
				return err
			}
		}

		logsegment.RecordStartMemoryAccessVector(w.builder, len(record.MemoryAccess))
		for i := len(record.MemoryAccess) - 1; i >= 0; i-- {
			a := record.MemoryAccess[i]
			logsegment.CreateMemoryAccess(w.builder,
				a.Offset,
				offset,
				uint32(len(a.Memory)),
				logsegment.MemoryAccessType(a.Access))
		}

		memory := w.builder.EndVector(len(record.MemoryAccess))
		params := w.prependUint64Vector(record.Params)
		results := w.prependUint64Vector(record.Results)
		timestamp := int64(record.Timestamp.Sub(w.startTime))

		function, ok := w.functions[record.Function]
		if !ok {
			return fmt.Errorf("record for unknown function: %s %s",
				record.Function.Module,
				record.Function.Name)
		}

		logsegment.RecordStart(w.builder)
		logsegment.RecordAddTimestamp(w.builder, timestamp)
		logsegment.RecordAddFunction(w.builder, function)
		logsegment.RecordAddParams(w.builder, params)
		logsegment.RecordAddResults(w.builder, results)
		logsegment.RecordAddOffset(w.builder, offset)
		logsegment.RecordAddLength(w.builder, length)
		logsegment.RecordAddMemoryAccess(w.builder, memory)
		w.records = append(w.records, logsegment.RecordEnd(w.builder))
	}

	records := w.prependOffsetVector(w.records)
	compressedSize := uint32(w.buffer.Len())

	logsegment.RecordBatchStart(w.builder)
	logsegment.RecordBatchAddCompressedSize(w.builder, compressedSize)
	logsegment.RecordBatchAddUncompressedSize(w.builder, uncompressedSize)
	logsegment.RecordBatchAddChecksum(w.builder, 0)
	logsegment.RecordBatchAddRecords(w.builder, records)
	w.builder.FinishSizePrefixed(logsegment.RecordBatchEnd(w.builder))

	_, err := w.output.Write(w.builder.FinishedBytes())
	if err != nil {
		return err
	}
	_, err = w.output.Write(w.buffer.Bytes())
	return err
}

func (w *LogWriter) prependHash(hash Hash) flatbuffers.UOffsetT {
	algorithm := w.builder.CreateSharedString(hash.Algorithm)
	digest := w.builder.CreateString(hash.Digest)
	types.HashStart(w.builder)
	types.HashAddAlgorithm(w.builder, algorithm)
	types.HashAddDigest(w.builder, digest)
	return types.HashEnd(w.builder)
}

func (w *LogWriter) prependStringVector(values []string) flatbuffers.UOffsetT {
	return w.prependObjectVector(len(values), func(i int) flatbuffers.UOffsetT {
		return w.builder.CreateString(values[i])
	})
}

func (w *LogWriter) prependUint64Vector(values []uint64) flatbuffers.UOffsetT {
	w.builder.StartVector(8, len(values), 8)
	for i := len(values) - 1; i >= 0; i-- {
		w.builder.PlaceUint64(values[i])
	}
	return w.builder.EndVector(len(values))
}

func (w *LogWriter) prependObjectVector(numElems int, create func(int) flatbuffers.UOffsetT) flatbuffers.UOffsetT {
	if numElems <= cap(w.offsets) {
		w.offsets = w.offsets[:numElems]
	} else {
		w.offsets = make([]flatbuffers.UOffsetT, numElems)
	}
	for i := range w.offsets {
		w.offsets[i] = create(i)
	}
	return w.prependOffsetVector(w.offsets)
}

func (w *LogWriter) prependOffsetVector(offsets []flatbuffers.UOffsetT) flatbuffers.UOffsetT {
	w.builder.StartVector(4, len(offsets), 4)
	for i := len(offsets) - 1; i >= 0; i-- {
		w.builder.PlaceUOffsetT(offsets[i])
	}
	return w.builder.EndVector(len(offsets))
}
