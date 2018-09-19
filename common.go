package geo

const (
	ewkbSRIDFlag uint32 = 0x20000000
	ewkbMFlag    uint32 = 0x40000000
	ewkbZFlag    uint32 = 0x80000000

	ewkbPointType           uint32 = 0x00000001
	ewkbLineStringType      uint32 = 0x00000002
	ewkbPolygonType         uint32 = 0x00000003
	ewkbMultiPointType      uint32 = 0x00000004
	ewkbMultiLineStringType uint32 = 0x00000005
	ewkbMultiPolygonType    uint32 = 0x00000006
)

type sridVal struct {
	SRID uint32 `json:"srid"`
}

func (s *sridVal) getSRID() uint32 {
	return s.SRID
}

func (s *sridVal) setSRID(srid uint32) {
	s.SRID = srid
}
