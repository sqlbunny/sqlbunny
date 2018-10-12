package geo

import (
	"bytes"
	"errors"
	"io"
)

const (
	ewkbSRIDFlag uint32 = 0x20000000
	ewkbMFlag    uint32 = 0x40000000
	ewkbZFlag    uint32 = 0x80000000
	ewkbFlagMask uint32 = 0xF0000000

	ewkbPointType           uint32 = 0x00000001
	ewkbLineStringType      uint32 = 0x00000002
	ewkbPolygonType         uint32 = 0x00000003
	ewkbMultiPointType      uint32 = 0x00000004
	ewkbMultiLineStringType uint32 = 0x00000005
	ewkbMultiPolygonType    uint32 = 0x00000006
)

type geom interface {
	ewkbType() uint32
	ewkbRead(r *ewkbReader)
	ewkbWrite(r *ewkbWriter)
}

type geomS interface {
	geom
	getSRID() uint32
	setSRID(srid uint32)
}

type srid uint32

func (s *srid) getSRID() uint32 {
	return uint32(*s)
}

func (s *srid) setSRID(v uint32) {
	*s = srid(v)
}

func Unmarshal(data []byte, g geom) error {
	return UnmarshalReader(bytes.NewReader(data), g)
}

func UnmarshalReader(rd io.Reader, g geom) error {
	r, err := newEWKBReader(rd)
	if err != nil {
		return err
	}

	ewkbType := r.ReadUint32()
	ewkbFlags := ewkbType & ewkbFlagMask
	ewkbType &= ^ewkbFlagMask
	haveSRID := (ewkbFlags & ewkbSRIDFlag) != 0

	r.flags = ewkbFlags

	if ewkbType != g.ewkbType() & ^ewkbFlagMask {
		return errors.New("Incorrect EWKB type")
	}
	var srid uint32
	if haveSRID {
		srid = r.ReadUint32()
	}
	if s, ok := g.(geomS); ok {
		s.setSRID(srid)
	}

	g.ewkbRead(&r)
	return nil
}

func Marshal(g geom) []byte {
	var b bytes.Buffer
	MarshalWriter(&b, g)
	return b.Bytes()
}

func MarshalWriter(wr io.Writer, g geom) {
	w := newEWKBWriter(wr)

	if s, ok := g.(geomS); ok {
		w.WriteUint32(g.ewkbType() | ewkbSRIDFlag)
		w.WriteUint32(s.getSRID())
	} else {
		w.WriteUint32(g.ewkbType())
	}
	g.ewkbWrite(&w)
}
