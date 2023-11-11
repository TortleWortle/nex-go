package nex

import (
	"reflect"

	crunch "github.com/superwhiskers/crunch/v3"
)

// StreamOut is an abstraction of github.com/superwhiskers/crunch with nex type support
type StreamOut struct {
	*crunch.Buffer
	Server ServerInterface
}

// WriteBool writes a bool
func (stream *StreamOut) WriteBool(b bool) {
	var bVar uint8
	if b {
		bVar = 1
	}
	stream.Grow(1)
	stream.WriteByteNext(byte(bVar))
}

// WriteUInt8 writes a uint8
func (stream *StreamOut) WriteUInt8(u8 uint8) {
	stream.Grow(1)
	stream.WriteByteNext(byte(u8))
}

// WriteInt8 writes a int8
func (stream *StreamOut) WriteInt8(s8 int8) {
	stream.Grow(1)
	stream.WriteByteNext(byte(s8))
}

// WriteUInt16LE writes a uint16 as LE
func (stream *StreamOut) WriteUInt16LE(u16 uint16) {
	stream.Grow(2)
	stream.WriteU16LENext([]uint16{u16})
}

// WriteUInt16BE writes a uint16 as BE
func (stream *StreamOut) WriteUInt16BE(u16 uint16) {
	stream.Grow(2)
	stream.WriteU16BENext([]uint16{u16})
}

// WriteInt16LE writes a uint16 as LE
func (stream *StreamOut) WriteInt16LE(s16 int16) {
	stream.Grow(2)
	stream.WriteU16LENext([]uint16{uint16(s16)})
}

// WriteInt16BE writes a uint16 as BE
func (stream *StreamOut) WriteInt16BE(s16 int16) {
	stream.Grow(2)
	stream.WriteU16BENext([]uint16{uint16(s16)})
}

// WriteUInt32LE writes a uint32 as LE
func (stream *StreamOut) WriteUInt32LE(u32 uint32) {
	stream.Grow(4)
	stream.WriteU32LENext([]uint32{u32})
}

// WriteUInt32BE writes a uint32 as BE
func (stream *StreamOut) WriteUInt32BE(u32 uint32) {
	stream.Grow(4)
	stream.WriteU32BENext([]uint32{u32})
}

// WriteInt32LE writes a int32 as LE
func (stream *StreamOut) WriteInt32LE(s32 int32) {
	stream.Grow(4)
	stream.WriteU32LENext([]uint32{uint32(s32)})
}

// WriteInt32BE writes a int32 as BE
func (stream *StreamOut) WriteInt32BE(s32 int32) {
	stream.Grow(4)
	stream.WriteU32BENext([]uint32{uint32(s32)})
}

// WriteUInt64LE writes a uint64 as LE
func (stream *StreamOut) WriteUInt64LE(u64 uint64) {
	stream.Grow(8)
	stream.WriteU64LENext([]uint64{u64})
}

// WriteUInt64BE writes a uint64 as BE
func (stream *StreamOut) WriteUInt64BE(u64 uint64) {
	stream.Grow(8)
	stream.WriteU64BENext([]uint64{u64})
}

// WriteInt64LE writes a int64 as LE
func (stream *StreamOut) WriteInt64LE(s64 int64) {
	stream.Grow(8)
	stream.WriteU64LENext([]uint64{uint64(s64)})
}

// WriteInt64BE writes a int64 as BE
func (stream *StreamOut) WriteInt64BE(s64 int64) {
	stream.Grow(8)
	stream.WriteU64BENext([]uint64{uint64(s64)})
}

// WriteFloat32LE writes a float32 as LE
func (stream *StreamOut) WriteFloat32LE(f32 float32) {
	stream.Grow(4)
	stream.WriteF32LENext([]float32{f32})
}

// WriteFloat32BE writes a float32 as BE
func (stream *StreamOut) WriteFloat32BE(f32 float32) {
	stream.Grow(4)
	stream.WriteF32BENext([]float32{f32})
}

// WriteFloat64LE writes a float64 as LE
func (stream *StreamOut) WriteFloat64LE(f64 float64) {
	stream.Grow(8)
	stream.WriteF64LENext([]float64{f64})
}

// WriteFloat64BE writes a float64 as BE
func (stream *StreamOut) WriteFloat64BE(f64 float64) {
	stream.Grow(8)
	stream.WriteF64BENext([]float64{f64})
}

// WriteString writes a NEX string type
func (stream *StreamOut) WriteString(str string) {
	str = str + "\x00"
	strLength := len(str)

	stream.Grow(int64(strLength))
	stream.WriteUInt16LE(uint16(strLength))
	stream.WriteBytesNext([]byte(str))
}

// WriteBuffer writes a NEX Buffer type
func (stream *StreamOut) WriteBuffer(data []byte) {
	dataLength := len(data)

	stream.WriteUInt32LE(uint32(dataLength))

	if dataLength > 0 {
		stream.Grow(int64(dataLength))
		stream.WriteBytesNext(data)
	}
}

// WriteQBuffer writes a NEX qBuffer type
func (stream *StreamOut) WriteQBuffer(data []byte) {
	dataLength := len(data)

	stream.WriteUInt16LE(uint16(dataLength))

	if dataLength > 0 {
		stream.Grow(int64(dataLength))
		stream.WriteBytesNext(data)
	}
}

// WriteResult writes a NEX Result type
func (stream *StreamOut) WriteResult(result *Result) {
	stream.WriteUInt32LE(result.Code)
}

// WriteStructure writes a nex Structure type
func (stream *StreamOut) WriteStructure(structure StructureInterface) {
	if structure.ParentType() != nil {
		stream.WriteStructure(structure.ParentType())
	}

	content := structure.Bytes(NewStreamOut(stream.Server))

	if stream.Server.ProtocolMinorVersion() >= 3 {
		stream.WriteUInt8(structure.StructureVersion())
		stream.WriteUInt32LE(uint32(len(content)))
	}

	stream.Grow(int64(len(content)))
	stream.WriteBytesNext(content)
}

// WriteListUInt8 writes a list of uint8 types
func (stream *StreamOut) WriteListUInt8(list []uint8) {
	stream.WriteUInt32LE(uint32(len(list)))

	for i := 0; i < len(list); i++ {
		stream.WriteUInt8(list[i])
	}
}

// WriteListInt8 writes a list of int8 types
func (stream *StreamOut) WriteListInt8(list []int8) {
	stream.WriteUInt32LE(uint32(len(list)))

	for i := 0; i < len(list); i++ {
		stream.WriteInt8(list[i])
	}
}

// WriteListUInt16LE writes a list of Little-Endian encoded uint16 types
func (stream *StreamOut) WriteListUInt16LE(list []uint16) {
	stream.WriteUInt32LE(uint32(len(list)))

	for i := 0; i < len(list); i++ {
		stream.WriteUInt16LE(list[i])
	}
}

// WriteListUInt16BE writes a list of Big-Endian encoded uint16 types
func (stream *StreamOut) WriteListUInt16BE(list []uint16) {
	stream.WriteUInt32LE(uint32(len(list)))

	for i := 0; i < len(list); i++ {
		stream.WriteUInt16BE(list[i])
	}
}

// WriteListInt16LE writes a list of Little-Endian encoded int16 types
func (stream *StreamOut) WriteListInt16LE(list []int16) {
	stream.WriteUInt32LE(uint32(len(list)))

	for i := 0; i < len(list); i++ {
		stream.WriteInt16LE(list[i])
	}
}

// WriteListInt16BE writes a list of Big-Endian encoded int16 types
func (stream *StreamOut) WriteListInt16BE(list []int16) {
	stream.WriteUInt32LE(uint32(len(list)))

	for i := 0; i < len(list); i++ {
		stream.WriteInt16BE(list[i])
	}
}

// WriteListUInt32LE writes a list of Little-Endian encoded uint32 types
func (stream *StreamOut) WriteListUInt32LE(list []uint32) {
	stream.WriteUInt32LE(uint32(len(list)))

	for i := 0; i < len(list); i++ {
		stream.WriteUInt32LE(list[i])
	}
}

// WriteListUInt32BE writes a list of Big-Endian encoded uint32 types
func (stream *StreamOut) WriteListUInt32BE(list []uint32) {
	stream.WriteUInt32LE(uint32(len(list)))

	for i := 0; i < len(list); i++ {
		stream.WriteUInt32BE(list[i])
	}
}

// WriteListInt32LE writes a list of Little-Endian encoded int32 types
func (stream *StreamOut) WriteListInt32LE(list []int32) {
	stream.WriteUInt32LE(uint32(len(list)))

	for i := 0; i < len(list); i++ {
		stream.WriteInt32LE(list[i])
	}
}

// WriteListInt32BE writes a list of Big-Endian encoded int32 types
func (stream *StreamOut) WriteListInt32BE(list []int32) {
	stream.WriteUInt32LE(uint32(len(list)))

	for i := 0; i < len(list); i++ {
		stream.WriteInt32BE(list[i])
	}
}

// WriteListUInt64LE writes a list of Little-Endian encoded uint64 types
func (stream *StreamOut) WriteListUInt64LE(list []uint64) {
	stream.WriteUInt32LE(uint32(len(list)))

	for i := 0; i < len(list); i++ {
		stream.WriteUInt64LE(list[i])
	}
}

// WriteListUInt64BE writes a list of Big-Endian encoded uint64 types
func (stream *StreamOut) WriteListUInt64BE(list []uint64) {
	stream.WriteUInt32LE(uint32(len(list)))

	for i := 0; i < len(list); i++ {
		stream.WriteUInt64BE(list[i])
	}
}

// WriteListInt64LE writes a list of Little-Endian encoded int64 types
func (stream *StreamOut) WriteListInt64LE(list []int64) {
	stream.WriteUInt32LE(uint32(len(list)))

	for i := 0; i < len(list); i++ {
		stream.WriteInt64LE(list[i])
	}
}

// WriteListInt64BE writes a list of Big-Endian encoded int64 types
func (stream *StreamOut) WriteListInt64BE(list []int64) {
	stream.WriteUInt32LE(uint32(len(list)))

	for i := 0; i < len(list); i++ {
		stream.WriteInt64BE(list[i])
	}
}

// WriteListFloat32LE writes a list of Little-Endian encoded float32 types
func (stream *StreamOut) WriteListFloat32LE(list []float32) {
	stream.WriteUInt32LE(uint32(len(list)))

	for i := 0; i < len(list); i++ {
		stream.WriteFloat32LE(list[i])
	}
}

// WriteListFloat32BE writes a list of Big-Endian encoded float32 types
func (stream *StreamOut) WriteListFloat32BE(list []float32) {
	stream.WriteUInt32LE(uint32(len(list)))

	for i := 0; i < len(list); i++ {
		stream.WriteFloat32BE(list[i])
	}
}

// WriteListFloat64LE writes a list of Little-Endian encoded float64 types
func (stream *StreamOut) WriteListFloat64LE(list []float64) {
	stream.WriteUInt32LE(uint32(len(list)))

	for i := 0; i < len(list); i++ {
		stream.WriteFloat64LE(list[i])
	}
}

// WriteListFloat64BE writes a list of Big-Endian encoded float64 types
func (stream *StreamOut) WriteListFloat64BE(list []float64) {
	stream.WriteUInt32LE(uint32(len(list)))

	for i := 0; i < len(list); i++ {
		stream.WriteFloat64BE(list[i])
	}
}

// WriteListStructure writes a list of NEX Structure types
func (stream *StreamOut) WriteListStructure(structures interface{}) {
	// TODO:
	// Find a better solution that doesn't use reflect

	slice := reflect.ValueOf(structures)
	count := slice.Len()

	stream.WriteUInt32LE(uint32(count))

	for i := 0; i < count; i++ {
		structure := slice.Index(i).Interface().(StructureInterface)
		stream.WriteStructure(structure)
	}
}

// WriteListString writes a list of NEX String types
func (stream *StreamOut) WriteListString(strings []string) {
	length := len(strings)

	stream.WriteUInt32LE(uint32(length))

	for i := 0; i < length; i++ {
		stream.WriteString(strings[i])
	}
}

// WriteListBuffer writes a list of NEX Buffer types
func (stream *StreamOut) WriteListBuffer(buffers [][]byte) {
	length := len(buffers)

	stream.WriteUInt32LE(uint32(length))

	for i := 0; i < length; i++ {
		stream.WriteBuffer(buffers[i])
	}
}

// WriteListQBuffer writes a list of NEX qBuffer types
func (stream *StreamOut) WriteListQBuffer(buffers [][]byte) {
	length := len(buffers)

	stream.WriteUInt32LE(uint32(length))

	for i := 0; i < length; i++ {
		stream.WriteQBuffer(buffers[i])
	}
}

// WriteListResult writes a list of NEX Result types
func (stream *StreamOut) WriteListResult(results []*Result) {
	length := len(results)

	stream.WriteUInt32LE(uint32(length))

	for i := 0; i < length; i++ {
		stream.WriteResult(results[i])
	}
}

// WriteListStationURL writes a list of NEX StationURL types
func (stream *StreamOut) WriteListStationURL(stationURLs []*StationURL) {
	length := len(stationURLs)

	stream.WriteUInt32LE(uint32(length))

	for i := 0; i < length; i++ {
		stream.WriteString(stationURLs[i].EncodeToString())
	}
}

// WriteListDataHolder writes a NEX DataHolder type
func (stream *StreamOut) WriteListDataHolder(dataholders []*DataHolder) {
	length := len(dataholders)

	stream.WriteUInt32LE(uint32(length))

	for i := 0; i < length; i++ {
		stream.WriteDataHolder(dataholders[i])
	}
}

// WriteDataHolder writes a NEX DataHolder type
func (stream *StreamOut) WriteDataHolder(dataholder *DataHolder) {
	content := dataholder.Bytes(NewStreamOut(stream.Server))
	stream.Grow(int64(len(content)))
	stream.WriteBytesNext(content)
}

// WriteDateTime writes a NEX DateTime type
func (stream *StreamOut) WriteDateTime(datetime *DateTime) {
	stream.WriteUInt64LE(datetime.value)
}

// WriteVariant writes a Variant type
func (stream *StreamOut) WriteVariant(variant *Variant) {
	content := variant.Bytes(NewStreamOut(stream.Server))
	stream.Grow(int64(len(content)))
	stream.WriteBytesNext(content)
}

// WriteMap writes a Map type with the given key and value types
func (stream *StreamOut) WriteMap(mapType interface{}) {
	// TODO:
	// Find a better solution that doesn't use reflect

	mapValue := reflect.ValueOf(mapType)
	count := mapValue.Len()

	stream.WriteUInt32LE(uint32(count))

	mapIter := mapValue.MapRange()

	for mapIter.Next() {
		key := mapIter.Key().Interface()
		value := mapIter.Value().Interface()

		switch key := key.(type) {
		case string:
			stream.WriteString(key)
		}

		switch value := value.(type) {
		case *Variant:
			stream.WriteVariant(value)
		}
	}
}

// NewStreamOut returns a new nex output stream
func NewStreamOut(server ServerInterface) *StreamOut {
	return &StreamOut{
		Buffer: crunch.NewBuffer(),
		Server: server,
	}
}
