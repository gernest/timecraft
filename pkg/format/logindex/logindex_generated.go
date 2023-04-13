// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package logindex

import (
	flatbuffers "github.com/google/flatbuffers/go"

	types "github.com/stealthrocket/timecraft/pkg/format/types"
)

type RecordIndex struct {
	_tab flatbuffers.Table
}

func GetRootAsRecordIndex(buf []byte, offset flatbuffers.UOffsetT) *RecordIndex {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &RecordIndex{}
	x.Init(buf, n+offset)
	return x
}

func GetSizePrefixedRootAsRecordIndex(buf []byte, offset flatbuffers.UOffsetT) *RecordIndex {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &RecordIndex{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func (rcv *RecordIndex) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *RecordIndex) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *RecordIndex) ProcessId(obj *types.Hash) *types.Hash {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(types.Hash)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *RecordIndex) Segment() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *RecordIndex) MutateSegment(n uint32) bool {
	return rcv._tab.MutateUint32Slot(6, n)
}

func (rcv *RecordIndex) Index(obj *Entry, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 16
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *RecordIndex) IndexLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func RecordIndexStart(builder *flatbuffers.Builder) {
	builder.StartObject(3)
}
func RecordIndexAddProcessId(builder *flatbuffers.Builder, processId flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(processId), 0)
}
func RecordIndexAddSegment(builder *flatbuffers.Builder, segment uint32) {
	builder.PrependUint32Slot(1, segment, 0)
}
func RecordIndexAddIndex(builder *flatbuffers.Builder, index flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(index), 0)
}
func RecordIndexStartIndexVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(16, numElems, 8)
}
func RecordIndexEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}

type Entry struct {
	_tab flatbuffers.Struct
}

func (rcv *Entry) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Entry) Table() flatbuffers.Table {
	return rcv._tab.Table
}

func (rcv *Entry) Key() uint64 {
	return rcv._tab.GetUint64(rcv._tab.Pos + flatbuffers.UOffsetT(0))
}
func (rcv *Entry) MutateKey(n uint64) bool {
	return rcv._tab.MutateUint64(rcv._tab.Pos+flatbuffers.UOffsetT(0), n)
}

func (rcv *Entry) Value() uint64 {
	return rcv._tab.GetUint64(rcv._tab.Pos + flatbuffers.UOffsetT(8))
}
func (rcv *Entry) MutateValue(n uint64) bool {
	return rcv._tab.MutateUint64(rcv._tab.Pos+flatbuffers.UOffsetT(8), n)
}

func CreateEntry(builder *flatbuffers.Builder, key uint64, value uint64) flatbuffers.UOffsetT {
	builder.Prep(8, 16)
	builder.PrependUint64(value)
	builder.PrependUint64(key)
	return builder.Offset()
}
