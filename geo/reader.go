package geo

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
)

type ewkbReader struct {
	r         io.Reader
	byteOrder binary.ByteOrder
	flags     uint32
}

func readByte(r io.Reader) byte {
	var buf [1]byte
	_, err := io.ReadFull(r, buf[:])
	if err != nil {
		panic(err)
	}
	return buf[0]
}

func newEWKBReader(r io.Reader) (ewkbReader, error) {
	var byteOrder binary.ByteOrder
	switch readByte(r) {
	case 0:
		byteOrder = binary.BigEndian
	case 1:
		byteOrder = binary.LittleEndian
	default:
		return ewkbReader{}, errors.New("Invalid byte order")
	}

	return ewkbReader{
		r:         r,
		byteOrder: byteOrder,
	}, nil
}

func (r *ewkbReader) ReadUint8() uint8 {
	var buf [1]byte
	_, err := io.ReadFull(r.r, buf[:])
	if err != nil {
		panic(err)
	}
	return buf[0]
}

func (r *ewkbReader) ReadUint16() uint16 {
	var buf [2]byte
	_, err := io.ReadFull(r.r, buf[:])
	if err != nil {
		panic(err)
	}
	return r.byteOrder.Uint16(buf[:])
}

func (r *ewkbReader) ReadUint32() uint32 {
	var buf [4]byte
	_, err := io.ReadFull(r.r, buf[:])
	if err != nil {
		panic(err)
	}
	return r.byteOrder.Uint32(buf[:])
}

func (r *ewkbReader) ReadUint64() uint64 {
	var buf [8]byte
	_, err := io.ReadFull(r.r, buf[:])
	if err != nil {
		panic(err)
	}
	return r.byteOrder.Uint64(buf[:])
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
