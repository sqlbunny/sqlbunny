package geo

import (
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"math"
)

type ewkbWriter struct {
	data      []byte
	byteOrder binary.ByteOrder
}

func newEWKBWriter() ewkbWriter {
	r := ewkbWriter{
		data:      []byte{0x01},
		byteOrder: binary.LittleEndian,
	}
	return r
}

func (r *ewkbWriter) Value() driver.Value {
	return hex.EncodeToString(r.data)
}

func (r *ewkbWriter) WriteUint8(v uint8) {
	r.data = append(r.data, v)
}

func (r *ewkbWriter) WriteUint16(v uint16) {
	var data [2]byte
	r.byteOrder.PutUint16(data[:], v)
	r.data = append(r.data, data[:]...)
}

func (r *ewkbWriter) WriteUint32(v uint32) {
	var data [4]byte
	r.byteOrder.PutUint32(data[:], v)
	r.data = append(r.data, data[:]...)
}

func (r *ewkbWriter) WriteUint64(v uint64) {
	var data [8]byte
	r.byteOrder.PutUint64(data[:], v)
	r.data = append(r.data, data[:]...)
}

func (r *ewkbWriter) WriteInt8(v int8) {
	r.WriteUint8(uint8(v))
}

func (r *ewkbWriter) WriteInt16(v int16) {
	r.WriteUint16(uint16(v))
}

func (r *ewkbWriter) WriteInt32(v int32) {
	r.WriteUint32(uint32(v))
}

func (r *ewkbWriter) WriteInt64(v int64) {
	r.WriteUint64(uint64(v))
}

func (r *ewkbWriter) WriteFloat32(v float32) {
	r.WriteUint32(math.Float32bits(v))
}

func (r *ewkbWriter) WriteFloat64(v float64) {
	r.WriteUint64(math.Float64bits(v))
}
