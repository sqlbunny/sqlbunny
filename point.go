package geo

import (
	"database/sql/driver"
)

type Point struct {
	X float64 `json:"longitude"`
	Y float64 `json:"latitude"`
}

type PointS struct {
	Point
	sridVal
}

func (p *Point) ewkbRead(r *ewkbReader) {
	p.X = r.ReadFloat64()
	p.Y = r.ReadFloat64()
}

func (p *Point) ewkbWrite(w *ewkbWriter) {
	w.WriteFloat64(p.X)
	w.WriteFloat64(p.Y)
}

func (p *Point) Scan(value interface{}) error {
	return scan(value, p)
}

func (p Point) Value() (driver.Value, error) {
	return value(&p)
}

func (p Point) ewkbType() uint32 {
	return ewkbPointType
}
