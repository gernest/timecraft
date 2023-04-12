// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package format

import (
	"strconv"
	flatbuffers "github.com/google/flatbuffers/go"
)

type Compression int8

const (
	CompressionNone   Compression = 0
	CompressionSnappy Compression = 1
	CompressionZstd   Compression = 2
)

var EnumNamesCompression = map[Compression]string{
	CompressionNone:   "None",
	CompressionSnappy: "Snappy",
	CompressionZstd:   "Zstd",
}

var EnumValuesCompression = map[string]Compression{
	"None":   CompressionNone,
	"Snappy": CompressionSnappy,
	"Zstd":   CompressionZstd,
}

func (v Compression) String() string {
	if s, ok := EnumNamesCompression[v]; ok {
		return s
	}
	return "Compression(" + strconv.FormatInt(int64(v), 10) + ")"
}

type MemoryAccessType uint32

const (
	MemoryAccessTypeMemoryRead  MemoryAccessType = 0
	MemoryAccessTypeMemoryWrite MemoryAccessType = 1
)

var EnumNamesMemoryAccessType = map[MemoryAccessType]string{
	MemoryAccessTypeMemoryRead:  "MemoryRead",
	MemoryAccessTypeMemoryWrite: "MemoryWrite",
}

var EnumValuesMemoryAccessType = map[string]MemoryAccessType{
	"MemoryRead":  MemoryAccessTypeMemoryRead,
	"MemoryWrite": MemoryAccessTypeMemoryWrite,
}

func (v MemoryAccessType) String() string {
	if s, ok := EnumNamesMemoryAccessType[v]; ok {
		return s
	}
	return "MemoryAccessType(" + strconv.FormatInt(int64(v), 10) + ")"
}

type UUID struct {
	_tab flatbuffers.Struct
}

func (rcv *UUID) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *UUID) Table() flatbuffers.Table {
	return rcv._tab.Table
}

func (rcv *UUID) Lo() uint64 {
	return rcv._tab.GetUint64(rcv._tab.Pos + flatbuffers.UOffsetT(0))
}
func (rcv *UUID) MutateLo(n uint64) bool {
	return rcv._tab.MutateUint64(rcv._tab.Pos+flatbuffers.UOffsetT(0), n)
}

func (rcv *UUID) Hi() uint64 {
	return rcv._tab.GetUint64(rcv._tab.Pos + flatbuffers.UOffsetT(8))
}
func (rcv *UUID) MutateHi(n uint64) bool {
	return rcv._tab.MutateUint64(rcv._tab.Pos+flatbuffers.UOffsetT(8), n)
}

func CreateUUID(builder *flatbuffers.Builder, lo uint64, hi uint64) flatbuffers.UOffsetT {
	builder.Prep(8, 16)
	builder.PrependUint64(hi)
	builder.PrependUint64(lo)
	return builder.Offset()
}
type Runtime struct {
	_tab flatbuffers.Table
}

func GetRootAsRuntime(buf []byte, offset flatbuffers.UOffsetT) *Runtime {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Runtime{}
	x.Init(buf, n+offset)
	return x
}

func GetSizePrefixedRootAsRuntime(buf []byte, offset flatbuffers.UOffsetT) *Runtime {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &Runtime{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func (rcv *Runtime) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Runtime) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Runtime) Version() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Runtime) Runtime() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Runtime) Modules(obj *Module, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *Runtime) ModulesLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func RuntimeStart(builder *flatbuffers.Builder) {
	builder.StartObject(3)
}
func RuntimeAddVersion(builder *flatbuffers.Builder, version flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(version), 0)
}
func RuntimeAddRuntime(builder *flatbuffers.Builder, runtime flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(runtime), 0)
}
func RuntimeAddModules(builder *flatbuffers.Builder, modules flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(modules), 0)
}
func RuntimeStartModulesVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func RuntimeEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
type Module struct {
	_tab flatbuffers.Table
}

func GetRootAsModule(buf []byte, offset flatbuffers.UOffsetT) *Module {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Module{}
	x.Init(buf, n+offset)
	return x
}

func GetSizePrefixedRootAsModule(buf []byte, offset flatbuffers.UOffsetT) *Module {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &Module{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func (rcv *Module) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Module) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Module) Module() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Module) Functions(j int) []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.ByteVector(a + flatbuffers.UOffsetT(j*4))
	}
	return nil
}

func (rcv *Module) FunctionsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func ModuleStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func ModuleAddModule(builder *flatbuffers.Builder, module flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(module), 0)
}
func ModuleAddFunctions(builder *flatbuffers.Builder, functions flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(functions), 0)
}
func ModuleStartFunctionsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func ModuleEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
type Process struct {
	_tab flatbuffers.Table
}

func GetRootAsProcess(buf []byte, offset flatbuffers.UOffsetT) *Process {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Process{}
	x.Init(buf, n+offset)
	return x
}

func GetSizePrefixedRootAsProcess(buf []byte, offset flatbuffers.UOffsetT) *Process {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &Process{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func (rcv *Process) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Process) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Process) Id(obj *UUID) *UUID {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := o + rcv._tab.Pos
		if obj == nil {
			obj = new(UUID)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *Process) UnixStartTime() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Process) MutateUnixStartTime(n int64) bool {
	return rcv._tab.MutateInt64Slot(6, n)
}

func (rcv *Process) Module() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Process) Arguments(j int) []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.ByteVector(a + flatbuffers.UOffsetT(j*4))
	}
	return nil
}

func (rcv *Process) ArgumentsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *Process) Environment(j int) []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.ByteVector(a + flatbuffers.UOffsetT(j*4))
	}
	return nil
}

func (rcv *Process) EnvironmentLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *Process) ParentProcessId(obj *UUID) *UUID {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		x := o + rcv._tab.Pos
		if obj == nil {
			obj = new(UUID)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *Process) ParentForkOffset() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Process) MutateParentForkOffset(n int64) bool {
	return rcv._tab.MutateInt64Slot(16, n)
}

func ProcessStart(builder *flatbuffers.Builder) {
	builder.StartObject(7)
}
func ProcessAddId(builder *flatbuffers.Builder, id flatbuffers.UOffsetT) {
	builder.PrependStructSlot(0, flatbuffers.UOffsetT(id), 0)
}
func ProcessAddUnixStartTime(builder *flatbuffers.Builder, unixStartTime int64) {
	builder.PrependInt64Slot(1, unixStartTime, 0)
}
func ProcessAddModule(builder *flatbuffers.Builder, module flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(module), 0)
}
func ProcessAddArguments(builder *flatbuffers.Builder, arguments flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(arguments), 0)
}
func ProcessStartArgumentsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func ProcessAddEnvironment(builder *flatbuffers.Builder, environment flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(environment), 0)
}
func ProcessStartEnvironmentVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func ProcessAddParentProcessId(builder *flatbuffers.Builder, parentProcessId flatbuffers.UOffsetT) {
	builder.PrependStructSlot(5, flatbuffers.UOffsetT(parentProcessId), 0)
}
func ProcessAddParentForkOffset(builder *flatbuffers.Builder, parentForkOffset int64) {
	builder.PrependInt64Slot(6, parentForkOffset, 0)
}
func ProcessEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
type LogHeader struct {
	_tab flatbuffers.Table
}

func GetRootAsLogHeader(buf []byte, offset flatbuffers.UOffsetT) *LogHeader {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &LogHeader{}
	x.Init(buf, n+offset)
	return x
}

func GetSizePrefixedRootAsLogHeader(buf []byte, offset flatbuffers.UOffsetT) *LogHeader {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &LogHeader{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func (rcv *LogHeader) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *LogHeader) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *LogHeader) Runtime(obj *Runtime) *Runtime {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(Runtime)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *LogHeader) Process(obj *Process) *Process {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(Process)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *LogHeader) Segment() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *LogHeader) MutateSegment(n uint32) bool {
	return rcv._tab.MutateUint32Slot(8, n)
}

func LogHeaderStart(builder *flatbuffers.Builder) {
	builder.StartObject(3)
}
func LogHeaderAddRuntime(builder *flatbuffers.Builder, runtime flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(runtime), 0)
}
func LogHeaderAddProcess(builder *flatbuffers.Builder, process flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(process), 0)
}
func LogHeaderAddSegment(builder *flatbuffers.Builder, segment uint32) {
	builder.PrependUint32Slot(2, segment, 0)
}
func LogHeaderEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
type RecordBatch struct {
	_tab flatbuffers.Table
}

func GetRootAsRecordBatch(buf []byte, offset flatbuffers.UOffsetT) *RecordBatch {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &RecordBatch{}
	x.Init(buf, n+offset)
	return x
}

func GetSizePrefixedRootAsRecordBatch(buf []byte, offset flatbuffers.UOffsetT) *RecordBatch {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &RecordBatch{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func (rcv *RecordBatch) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *RecordBatch) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *RecordBatch) CompressedSize() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *RecordBatch) MutateCompressedSize(n uint32) bool {
	return rcv._tab.MutateUint32Slot(4, n)
}

func (rcv *RecordBatch) UncompressedSize() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *RecordBatch) MutateUncompressedSize(n uint32) bool {
	return rcv._tab.MutateUint32Slot(6, n)
}

func (rcv *RecordBatch) Checksum() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *RecordBatch) MutateChecksum(n uint32) bool {
	return rcv._tab.MutateUint32Slot(8, n)
}

func (rcv *RecordBatch) Compression() Compression {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return Compression(rcv._tab.GetInt8(o + rcv._tab.Pos))
	}
	return 0
}

func (rcv *RecordBatch) MutateCompression(n Compression) bool {
	return rcv._tab.MutateInt8Slot(10, int8(n))
}

func (rcv *RecordBatch) Records(obj *Record, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *RecordBatch) RecordsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func RecordBatchStart(builder *flatbuffers.Builder) {
	builder.StartObject(5)
}
func RecordBatchAddCompressedSize(builder *flatbuffers.Builder, compressedSize uint32) {
	builder.PrependUint32Slot(0, compressedSize, 0)
}
func RecordBatchAddUncompressedSize(builder *flatbuffers.Builder, uncompressedSize uint32) {
	builder.PrependUint32Slot(1, uncompressedSize, 0)
}
func RecordBatchAddChecksum(builder *flatbuffers.Builder, checksum uint32) {
	builder.PrependUint32Slot(2, checksum, 0)
}
func RecordBatchAddCompression(builder *flatbuffers.Builder, compression Compression) {
	builder.PrependInt8Slot(3, int8(compression), 0)
}
func RecordBatchAddRecords(builder *flatbuffers.Builder, records flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(records), 0)
}
func RecordBatchStartRecordsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func RecordBatchEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
type Record struct {
	_tab flatbuffers.Table
}

func GetRootAsRecord(buf []byte, offset flatbuffers.UOffsetT) *Record {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Record{}
	x.Init(buf, n+offset)
	return x
}

func GetSizePrefixedRootAsRecord(buf []byte, offset flatbuffers.UOffsetT) *Record {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &Record{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func (rcv *Record) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Record) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Record) Timestamp() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Record) MutateTimestamp(n int64) bool {
	return rcv._tab.MutateInt64Slot(4, n)
}

func (rcv *Record) Module() uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetUint16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Record) MutateModule(n uint16) bool {
	return rcv._tab.MutateUint16Slot(6, n)
}

func (rcv *Record) Function() uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetUint16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Record) MutateFunction(n uint16) bool {
	return rcv._tab.MutateUint16Slot(8, n)
}

func (rcv *Record) Params(j int) uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetUint64(a + flatbuffers.UOffsetT(j*8))
	}
	return 0
}

func (rcv *Record) ParamsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *Record) MutateParams(j int, n uint64) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateUint64(a+flatbuffers.UOffsetT(j*8), n)
	}
	return false
}

func (rcv *Record) Results(j int) uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetUint64(a + flatbuffers.UOffsetT(j*8))
	}
	return 0
}

func (rcv *Record) ResultsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *Record) MutateResults(j int, n uint64) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateUint64(a+flatbuffers.UOffsetT(j*8), n)
	}
	return false
}

func (rcv *Record) Offset() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Record) MutateOffset(n uint32) bool {
	return rcv._tab.MutateUint32Slot(14, n)
}

func (rcv *Record) Length() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Record) MutateLength(n uint32) bool {
	return rcv._tab.MutateUint32Slot(16, n)
}

func (rcv *Record) MemoryAccess(obj *MemoryAccess, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 16
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *Record) MemoryAccessLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func RecordStart(builder *flatbuffers.Builder) {
	builder.StartObject(8)
}
func RecordAddTimestamp(builder *flatbuffers.Builder, timestamp int64) {
	builder.PrependInt64Slot(0, timestamp, 0)
}
func RecordAddModule(builder *flatbuffers.Builder, module uint16) {
	builder.PrependUint16Slot(1, module, 0)
}
func RecordAddFunction(builder *flatbuffers.Builder, function uint16) {
	builder.PrependUint16Slot(2, function, 0)
}
func RecordAddParams(builder *flatbuffers.Builder, params flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(params), 0)
}
func RecordStartParamsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(8, numElems, 8)
}
func RecordAddResults(builder *flatbuffers.Builder, results flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(results), 0)
}
func RecordStartResultsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(8, numElems, 8)
}
func RecordAddOffset(builder *flatbuffers.Builder, offset uint32) {
	builder.PrependUint32Slot(5, offset, 0)
}
func RecordAddLength(builder *flatbuffers.Builder, length uint32) {
	builder.PrependUint32Slot(6, length, 0)
}
func RecordAddMemoryAccess(builder *flatbuffers.Builder, memoryAccess flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(7, flatbuffers.UOffsetT(memoryAccess), 0)
}
func RecordStartMemoryAccessVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(16, numElems, 4)
}
func RecordEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
type MemoryAccess struct {
	_tab flatbuffers.Struct
}

func (rcv *MemoryAccess) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *MemoryAccess) Table() flatbuffers.Table {
	return rcv._tab.Table
}

func (rcv *MemoryAccess) MemoryOffset() uint32 {
	return rcv._tab.GetUint32(rcv._tab.Pos + flatbuffers.UOffsetT(0))
}
func (rcv *MemoryAccess) MutateMemoryOffset(n uint32) bool {
	return rcv._tab.MutateUint32(rcv._tab.Pos+flatbuffers.UOffsetT(0), n)
}

func (rcv *MemoryAccess) RecordOffset() uint32 {
	return rcv._tab.GetUint32(rcv._tab.Pos + flatbuffers.UOffsetT(4))
}
func (rcv *MemoryAccess) MutateRecordOffset(n uint32) bool {
	return rcv._tab.MutateUint32(rcv._tab.Pos+flatbuffers.UOffsetT(4), n)
}

func (rcv *MemoryAccess) Length() uint32 {
	return rcv._tab.GetUint32(rcv._tab.Pos + flatbuffers.UOffsetT(8))
}
func (rcv *MemoryAccess) MutateLength(n uint32) bool {
	return rcv._tab.MutateUint32(rcv._tab.Pos+flatbuffers.UOffsetT(8), n)
}

func (rcv *MemoryAccess) Access() MemoryAccessType {
	return MemoryAccessType(rcv._tab.GetUint32(rcv._tab.Pos + flatbuffers.UOffsetT(12)))
}
func (rcv *MemoryAccess) MutateAccess(n MemoryAccessType) bool {
	return rcv._tab.MutateUint32(rcv._tab.Pos+flatbuffers.UOffsetT(12), uint32(n))
}

func CreateMemoryAccess(builder *flatbuffers.Builder, memoryOffset uint32, recordOffset uint32, length uint32, access MemoryAccessType) flatbuffers.UOffsetT {
	builder.Prep(4, 16)
	builder.PrependUint32(uint32(access))
	builder.PrependUint32(length)
	builder.PrependUint32(recordOffset)
	builder.PrependUint32(memoryOffset)
	return builder.Offset()
}