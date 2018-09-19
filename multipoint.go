package geo

import (
	"database/sql/driver"
	"encoding/json"
)

type MultiPoint struct {
	Points []Point `json:"points"`
}

// MarshalJSON implements json.Marshaler.
func (p *MultiPoint) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Points)
}

// UnmarshalJSON MultiPoint json.Unmarshaler.
func (p *MultiPoint) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &p.Points)
}

type MultiPointS struct {
	MultiPoint
	sridVal
}

func (p *MultiPoint) ewkbRead(r *ewkbReader) {
	n := r.ReadUint32()
	p.Points = make([]Point, n)
	for i := range p.Points {
		r.ReadUint8()
		r.ReadUint32()
		p.Points[i].ewkbRead(r)
	}
}

func (p *MultiPoint) ewkbWrite(w *ewkbWriter) {
	w.WriteUint32(uint32(len(p.Points)))
	for i := range p.Points {
		w.WriteUint8(0x01)
		w.WriteUint32(p.Points[i].ewkbType())
		p.Points[i].ewkbWrite(w)
	}
}

func (p *MultiPoint) Scan(value interface{}) error {
	return scan(value, p)
}

func (p MultiPoint) Value() (driver.Value, error) {
	return value(&p)
}

func (p MultiPoint) ewkbType() uint32 {
	return ewkbMultiPointType
}
