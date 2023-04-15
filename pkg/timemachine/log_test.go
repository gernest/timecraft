package timemachine_test

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stealthrocket/timecraft/pkg/timemachine"
)

func TestReadLogHeader(t *testing.T) {
	b := new(bytes.Buffer)
	w := timemachine.NewLogWriter(b)

	header := &timemachine.LogHeader{
		Runtime: timemachine.Runtime{
			Runtime: "test",
			Version: "dev",
			Functions: []timemachine.Function{
				{Module: "env", Name: "f0"},
				{Module: "env", Name: "f1"},
				{Module: "env", Name: "f2"},
				{Module: "env", Name: "f3"},
				{Module: "env", Name: "f4"},
			},
		},
		Process: timemachine.Process{
			ID:        timemachine.Hash{"sha", "f572d396fae9206628714fb2ce00f72e94f2258f"},
			Image:     timemachine.Hash{"sha", "28935580a9bbb8cc7bcdea62e7dfdcf7e0f31f87"},
			StartTime: time.Now(),
			Args:      os.Args,
			Environ:   os.Environ(),
		},
		Segment:     42,
		Compression: timemachine.Zstd,
	}

	if err := w.WriteLogHeader(header); err != nil {
		t.Fatal(err)
	}

	r0 := bytes.NewReader(b.Bytes())
	r1 := timemachine.NewLogReader(r0)

	for i := 0; i < 10; i++ {
		h, err := r1.ReadLogHeader()
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(h, header); diff != "" {
			t.Fatal(diff)
		}
		r0.Reset(b.Bytes())
		r1.Reset(r0)
	}
}

func BenchmarkLogReader(b *testing.B) {
	buffer := new(bytes.Buffer)
	writer := timemachine.NewLogWriter(buffer)

	header := &timemachine.LogHeader{
		Runtime: timemachine.Runtime{
			Runtime: "test",
			Version: "dev",
			Functions: []timemachine.Function{
				{Module: "env", Name: "f0"},
				{Module: "env", Name: "f1"},
				{Module: "env", Name: "f2"},
				{Module: "env", Name: "f3"},
				{Module: "env", Name: "f4"},
			},
		},
		Process: timemachine.Process{
			ID:        timemachine.Hash{"sha", "f572d396fae9206628714fb2ce00f72e94f2258f"},
			Image:     timemachine.Hash{"sha", "28935580a9bbb8cc7bcdea62e7dfdcf7e0f31f87"},
			StartTime: time.Now(),
			Args:      os.Args,
			Environ:   os.Environ(),
		},
		Segment:     42,
		Compression: timemachine.Zstd,
	}

	if err := writer.WriteLogHeader(header); err != nil {
		b.Fatal(err)
	}

	b.Run("ReadLogHeader", func(b *testing.B) {
		r0 := bytes.NewReader(buffer.Bytes())
		r1 := timemachine.NewLogReader(r0)

		for i := 0; i < b.N; i++ {
			_, err := r1.ReadLogHeader()
			if err != nil {
				b.Fatal(i, err)
			}
			r0.Reset(buffer.Bytes())
			r1.Reset(r0)
		}
	})
}

func BenchmarkLogWriter(b *testing.B) {
	b.Run("WriteLogHeader", func(b *testing.B) {
		tests := []struct {
			scenario string
			header   *timemachine.LogHeader
		}{
			{
				scenario: "common log header",
				header: &timemachine.LogHeader{
					Runtime: timemachine.Runtime{
						Runtime: "test",
						Version: "dev",
						Functions: []timemachine.Function{
							{Module: "env", Name: "f0"},
							{Module: "env", Name: "f1"},
							{Module: "env", Name: "f2"},
							{Module: "env", Name: "f3"},
							{Module: "env", Name: "f4"},
						},
					},
					Process: timemachine.Process{
						ID:        timemachine.Hash{"sha", "f572d396fae9206628714fb2ce00f72e94f2258f"},
						Image:     timemachine.Hash{"sha", "28935580a9bbb8cc7bcdea62e7dfdcf7e0f31f87"},
						StartTime: time.Now(),
						Args:      os.Args,
						Environ:   os.Environ(),
					},
					Segment:     42,
					Compression: timemachine.Zstd,
				},
			},
		}

		for _, test := range tests {
			b.Run(test.scenario, func(b *testing.B) {
				benchmarkLogWriterWriteLogHeader(b, test.header)
			})
		}
	})

	b.Run("WriteRecordBatch", func(b *testing.B) {
		header := &timemachine.LogHeader{
			Runtime: timemachine.Runtime{
				Runtime: "test",
				Version: "dev",
				Functions: []timemachine.Function{
					{Module: "env", Name: "f0"},
					{Module: "env", Name: "f1"},
					{Module: "env", Name: "f2"},
					{Module: "env", Name: "f3"},
					{Module: "env", Name: "f4"},
				},
			},
			Process: timemachine.Process{
				ID:        timemachine.Hash{"sha", "f572d396fae9206628714fb2ce00f72e94f2258f"},
				Image:     timemachine.Hash{"sha", "28935580a9bbb8cc7bcdea62e7dfdcf7e0f31f87"},
				StartTime: time.Now(),
				Args:      os.Args,
				Environ:   os.Environ(),
			},
			Segment:     42,
			Compression: timemachine.Zstd,
		}

		tests := []struct {
			scenario string
			batch    []timemachine.Record
		}{
			{
				scenario: "zero records",
			},

			{
				scenario: "one record",
				batch: []timemachine.Record{
					{
						Timestamp: header.Process.StartTime.Add(1 * time.Millisecond),
						Function:  0,
						Params:    []uint64{1},
						Results:   []uint64{42},
						MemoryAccess: []timemachine.MemoryAccess{
							{Memory: []byte("hello world!"), Offset: 1234, Access: timemachine.MemoryRead},
						},
					},
				},
			},

			{
				scenario: "five records",
				batch: []timemachine.Record{
					{
						Timestamp: header.Process.StartTime.Add(1 * time.Millisecond),
						Function:  0,
						Params:    []uint64{1},
						Results:   []uint64{42},
						MemoryAccess: []timemachine.MemoryAccess{
							{Memory: []byte("hello world!"), Offset: 1234, Access: timemachine.MemoryRead},
						},
					},
					{
						Timestamp: header.Process.StartTime.Add(2 * time.Millisecond),
						Function:  1,
						Params:    []uint64{1, 2},
						Results:   []uint64{42},
					},
					{
						Timestamp: header.Process.StartTime.Add(3 * time.Millisecond),
						Function:  2,
						Params:    []uint64{1, 2, 3},
						Results:   []uint64{42},
					},
					{
						Timestamp: header.Process.StartTime.Add(4 * time.Millisecond),
						Function:  3,
						MemoryAccess: []timemachine.MemoryAccess{
							{Memory: []byte("A"), Offset: 1, Access: timemachine.MemoryRead},
							{Memory: []byte("B"), Offset: 2, Access: timemachine.MemoryRead},
							{Memory: []byte("C"), Offset: 3, Access: timemachine.MemoryRead},
							{Memory: []byte("D"), Offset: 4, Access: timemachine.MemoryRead},
						},
					},
					{
						Timestamp: header.Process.StartTime.Add(5 * time.Millisecond),
						Function:  4,
						Params:    []uint64{1},
						Results:   []uint64{42},
						MemoryAccess: []timemachine.MemoryAccess{
							{Memory: []byte("hello world!"), Offset: 1234, Access: timemachine.MemoryRead},
							{Memory: make([]byte, 10e3), Offset: 1234567, Access: timemachine.MemoryWrite},
						},
					},
				},
			},
		}

		for _, test := range tests {
			b.Run(test.scenario, func(b *testing.B) {
				benchmarkLogWriterWriteRecordBatch(b, header, test.batch)
			})
		}
	})
}

func benchmarkLogWriterWriteLogHeader(b *testing.B, header *timemachine.LogHeader) {
	w := timemachine.NewLogWriter(io.Discard)

	for i := 0; i < b.N; i++ {
		if err := w.WriteLogHeader(header); err != nil {
			b.Fatal(err)
		}
		w.Reset(io.Discard)
	}
}

func benchmarkLogWriterWriteRecordBatch(b *testing.B, header *timemachine.LogHeader, batch []timemachine.Record) {
	w := timemachine.NewLogWriter(io.Discard)
	w.WriteLogHeader(header)

	for i := 0; i < b.N; i++ {
		if _, err := w.WriteRecordBatch(batch); err != nil {
			b.Fatal(err)
		}
	}
}
