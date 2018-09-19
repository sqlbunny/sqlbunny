package geo

import (
	"database/sql/driver"
	"encoding/json"
)

type MultiLineString struct {
	LineStrings []LineString `json:"line_strings"`
}

// MarshalJSON implements json.Marshaler.
func (p *MultiLineString) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.LineStrings)
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *MultiLineString) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &p.LineStrings)
}

type MultiLineStringS struct {
	MultiLineString
	sridVal
}

func (p *MultiLineString) ewkbRead(r *ewkbReader) {
	n := r.ReadUint32()
	p.LineStrings = make([]LineString, n)
	for i := range p.LineStrings {
		r.ReadUint8()
		r.ReadUint32()
		p.LineStrings[i].ewkbRead(r)
	}
}

func (p *MultiLineString) ewkbWrite(w *ewkbWriter) {
	w.WriteUint32(uint32(len(p.LineStrings)))
	for i := range p.LineStrings {
		w.WriteUint8(0x01)
		w.WriteUint32(p.LineStrings[i].ewkbType())
		p.LineStrings[i].ewkbWrite(w)
	}
}

func (p *MultiLineString) Scan(value interface{}) error {
	return scan(value, p)
}

func (p MultiLineString) Value() (driver.Value, error) {
	return value(&p)
}

func (p MultiLineString) ewkbType() uint32 {
	return ewkbMultiLineStringType
}
