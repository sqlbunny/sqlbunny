package geo

import (
	"database/sql/driver"
	"encoding/json"
)

type MultiPolygon struct {
	Polygons []Polygon `json:"polygons"`
}

// MarshalJSON implements json.Marshaler.
func (p *MultiPolygon) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Polygons)
}

// UnmarshalJSON MultiPoint json.Unmarshaler.
func (p *MultiPolygon) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &p.Polygons)
}

type MultiPolygonS struct {
	MultiPolygon
	sridVal
}

func (p *MultiPolygon) ewkbRead(r *ewkbReader) {
	n := r.ReadUint32()
	p.Polygons = make([]Polygon, n)
	for i := range p.Polygons {
		r.ReadUint8()
		r.ReadUint32()
		p.Polygons[i].ewkbRead(r)
	}
}

func (p *MultiPolygon) ewkbWrite(w *ewkbWriter) {
	w.WriteUint32(uint32(len(p.Polygons)))
	for i := range p.Polygons {
		w.WriteUint8(0x01)
		w.WriteUint32(p.Polygons[i].ewkbType())
		p.Polygons[i].ewkbWrite(w)
	}
}

func (p *MultiPolygon) Scan(value interface{}) error {
	return scan(value, p)
}

func (p *MultiPolygonS) Scan(value interface{}) error {
	return scan(value, p)
}

func (p MultiPolygon) Value() (driver.Value, error) {
	return value(&p)
}

func (p MultiPolygon) ewkbType() uint32 {
	return ewkbMultiPolygonType
}
