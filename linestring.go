package geo

import (
	"database/sql/driver"
	"encoding/json"
)

type LineString struct {
	Points []Point `json:"points"`
}

// MarshalJSON implements json.Marshaler.
func (p *LineString) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Points)
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *LineString) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &p.Points)
}

type LineStringS struct {
	LineString
	sridVal
}

func (p *LineString) ewkbRead(r *ewkbReader) {
	n := r.ReadUint32()
	p.Points = make([]Point, n)
	for i := range p.Points {
		p.Points[i].ewkbRead(r)
	}
}

func (p *LineString) ewkbWrite(w *ewkbWriter) {
	w.WriteUint32(uint32(len(p.Points)))
	for i := range p.Points {
		p.Points[i].ewkbWrite(w)
	}
}

func (p *LineString) Scan(value interface{}) error {
	return scan(value, p)
}

func (p LineString) Value() (driver.Value, error) {
	return value(&p)
}

func (p LineString) ewkbType() uint32 {
	return ewkbLineStringType
}
