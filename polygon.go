package geo

import (
	"database/sql/driver"
)

type Polygon struct {
	Coordinates LineString   `json:"coordinates"`
	Holes       []LineString `json:"holes,omitempty"`
}

type PolygonS struct {
	Polygon
	sridVal
}

func (p *Polygon) ewkbRead(r *ewkbReader) {
	n := r.ReadUint32()
	if n == 1 {
		p.Holes = nil
	} else {
		p.Holes = make([]LineString, n-1)
	}
	p.Coordinates.ewkbRead(r)
	for i := range p.Holes {
		p.Holes[i].ewkbRead(r)
	}
}

func (p *Polygon) ewkbWrite(w *ewkbWriter) {
	w.WriteUint32(uint32(len(p.Holes) + 1))
	p.Coordinates.ewkbWrite(w)
	for i := range p.Holes {
		p.Holes[i].ewkbWrite(w)
	}
}

func (p *Polygon) Scan(value interface{}) error {
	return scan(value, p)
}

func (p Polygon) Value() (driver.Value, error) {
	return value(&p)
}

func (p Polygon) ewkbType() uint32 {
	return ewkbPolygonType
}
