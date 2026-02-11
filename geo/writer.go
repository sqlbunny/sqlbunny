package geo

import (
	"encoding/binary"
	"io"
	"math"
)

type ewkbWriter struct {
	w         io.Writer
	byteOrder binary.ByteOrder
}

func newEWKBWriter(w io.Writer) ewkbWriter {
	var buf = [1]byte{0x01}
	_, err := w.Write(buf[:])
	if err != nil {
		panic(err)
	}

	return ewkbWriter{
		w:         w,
		byteOrder: binary.LittleEndian,
	}
}

func (w *ewkbWriter) WriteUint8(v uint8) {
	var buf = [1]byte{v}
	_, err := w.w.Write(buf[:])
	if err != nil {
		panic(err)
	}
}

func (w *ewkbWriter) WriteUint16(v uint16) {
	var data [2]byte
	w.byteOrder.PutUint16(data[:], v)
	_, err := w.w.Write(data[:])
	if err != nil {
		panic(err)
	}
}

func (w *ewkbWriter) WriteUint32(v uint32) {
	var data [4]byte
	w.byteOrder.PutUint32(data[:], v)
	_, err := w.w.Write(data[:])
	if err != nil {
		panic(err)
	}
}

func (w *ewkbWriter) WriteUint64(v uint64) {
	var data [8]byte
	w.byteOrder.PutUint64(data[:], v)
	_, err := w.w.Write(data[:])
	if err != nil {
		panic(err)
	}
}

func (w *ewkbWriter) WriteInt8(v int8) {
	w.WriteUint8(uint8(v))
}

func (w *ewkbWriter) WriteInt16(v int16) {
	w.WriteUint16(uint16(v))
}

func (w *ewkbWriter) WriteInt32(v int32) {
	w.WriteUint32(uint32(v))
}

func (w *ewkbWriter) WriteInt64(v int64) {
	w.WriteUint64(uint64(v))
}

func (w *ewkbWriter) WriteFloat32(v float32) {
	w.WriteUint32(math.Float32bits(v))
}

func (w *ewkbWriter) WriteFloat64(v float64) {
	w.WriteUint64(math.Float64bits(v))
}
