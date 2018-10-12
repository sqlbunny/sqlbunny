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
	srid `json:"srid"`
}

func (g *Point) ewkbRead(r *ewkbReader) {
	g.X = r.ReadFloat64()
	g.Y = r.ReadFloat64()
	if r.flags&ewkbZFlag != 0 {
		r.ReadFloat64()
	}
	if r.flags&ewkbMFlag != 0 {
		r.ReadFloat64()
	}
}

func (g *Point) ewkbWrite(w *ewkbWriter) {
	w.WriteFloat64(g.X)
	w.WriteFloat64(g.Y)
}

func (g *Point) Scan(value interface{}) error {
	return scan(value, g)
}

func (g Point) Value() (driver.Value, error) {
	return value(&g)
}

func (g Point) ewkbType() uint32 {
	return ewkbPointType
}

type PointZ struct {
	X float64 `json:"longitude"`
	Y float64 `json:"latitude"`
	Z float64 `json:"z"`
}

type PointZS struct {
	PointZ
	srid `json:"srid"`
}

func (g *PointZ) ewkbRead(r *ewkbReader) {
	g.X = r.ReadFloat64()
	g.Y = r.ReadFloat64()
	if r.flags&ewkbZFlag != 0 {
		g.Z = r.ReadFloat64()
	} else {
		g.Z = 0
	}
	if r.flags&ewkbMFlag != 0 {
		r.ReadFloat64()
	}
}

func (g *PointZ) ewkbWrite(w *ewkbWriter) {
	w.WriteFloat64(g.X)
	w.WriteFloat64(g.Y)
	w.WriteFloat64(g.Z)
}

func (g *PointZ) Scan(value interface{}) error {
	return scan(value, g)
}

func (g PointZ) Value() (driver.Value, error) {
	return value(&g)
}

func (g PointZ) ewkbType() uint32 {
	return ewkbPointType | ewkbZFlag
}

type PointM struct {
	X float64 `json:"longitude"`
	Y float64 `json:"latitude"`
	M float64 `json:"m"`
}

type PointMS struct {
	PointM
	srid `json:"srid"`
}

func (g *PointM) ewkbRead(r *ewkbReader) {
	g.X = r.ReadFloat64()
	g.Y = r.ReadFloat64()
	if r.flags&ewkbZFlag != 0 {
		r.ReadFloat64()
	}
	if r.flags&ewkbMFlag != 0 {
		g.M = r.ReadFloat64()
	} else {
		g.M = 0
	}
}

func (g *PointM) ewkbWrite(w *ewkbWriter) {
	w.WriteFloat64(g.X)
	w.WriteFloat64(g.Y)
	w.WriteFloat64(g.M)
}

func (g *PointM) Scan(value interface{}) error {
	return scan(value, g)
}

func (g PointM) Value() (driver.Value, error) {
	return value(&g)
}

func (g PointM) ewkbType() uint32 {
	return ewkbPointType | ewkbMFlag
}

type PointZM struct {
	X float64 `json:"longitude"`
	Y float64 `json:"latitude"`
	Z float64 `json:"z"`
	M float64 `json:"m"`
}

type PointZMS struct {
	PointZM
	srid `json:"srid"`
}

func (g *PointZM) ewkbRead(r *ewkbReader) {
	g.X = r.ReadFloat64()
	g.Y = r.ReadFloat64()
	if r.flags&ewkbZFlag != 0 {
		g.Z = r.ReadFloat64()
	} else {
		g.Z = 0
	}
	if r.flags&ewkbMFlag != 0 {
		g.M = r.ReadFloat64()
	} else {
		g.M = 0
	}
}

func (g *PointZM) ewkbWrite(w *ewkbWriter) {
	w.WriteFloat64(g.X)
	w.WriteFloat64(g.Y)
	w.WriteFloat64(g.Z)
	w.WriteFloat64(g.M)
}

func (g *PointZM) Scan(value interface{}) error {
	return scan(value, g)
}

func (g PointZM) Value() (driver.Value, error) {
	return value(&g)
}

func (g PointZM) ewkbType() uint32 {
	return ewkbPointType | ewkbZFlag | ewkbMFlag
}
