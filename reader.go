package geo

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"math"
)

type ewkbReader struct {
	data      []byte
	index     int
	byteOrder binary.ByteOrder
}

func newEWKBReader(value interface{}) (ewkbReader, error) {
	b, ok := value.([]byte)
	if !ok {
		return ewkbReader{}, errors.New("EWKB scan: value is not byte slice")
	}

	data, err := hex.DecodeString(string(b))
	if err != nil {
		return ewkbReader{}, err
	}

	var byteOrder binary.ByteOrder
	switch data[0] {
	case 0:
		byteOrder = binary.BigEndian
	case 1:
		byteOrder = binary.LittleEndian
	default:
		return ewkbReader{}, errors.New("Invalid byte order")
	}

	return ewkbReader{
		data:      data,
		index:     1,
		byteOrder: byteOrder,
	}, nil
}

func (r *ewkbReader) ReadUint8() uint8 {
	res := r.data[r.index]
	r.index += 1
	return res
}

func (r *ewkbReader) ReadUint16() uint16 {
	res := r.byteOrder.Uint16(r.data[r.index:])
	r.index += 2
	return res
}

func (r *ewkbReader) ReadUint32() uint32 {
	res := r.byteOrder.Uint32(r.data[r.index:])
	r.index += 4
	return res
}

func (r *ewkbReader) ReadUint64() uint64 {
	res := r.byteOrder.Uint64(r.data[r.index:])
	r.index += 8
	return res
}

func (r *ewkbReader) ReadInt8() int8 {
	return int8(r.ReadUint8())
}

func (r *ewkbReader) ReadInt16() int16 {
	return int16(r.ReadUint16())
}

func (r *ewkbReader) ReadInt32() int32 {
	return int32(r.ReadUint32())
}

func (r *ewkbReader) ReadInt64() int64 {
	return int64(r.ReadUint64())
}

func (r *ewkbReader) ReadFloat32() float32 {
	return math.Float32frombits(r.ReadUint32())
}

func (r *ewkbReader) ReadFloat64() float64 {
	return math.Float64frombits(r.ReadUint64())
}
