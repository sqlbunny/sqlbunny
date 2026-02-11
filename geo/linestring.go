package geo

import (
	"database/sql/driver"
)

type LineString []Point

type LineStringS struct {
	LineString
	srid `json:"srid"`
}

func (g *LineString) ewkbRead(r *ewkbReader) {
	n := r.ReadUint32()
	*g = make([]Point, n)
	for i := range *g {
		(*g)[i].ewkbRead(r)
	}
}

func (g *LineString) ewkbWrite(w *ewkbWriter) {
	w.WriteUint32(uint32(len(*g)))
	for _, p := range *g {
		p.ewkbWrite(w)
	}
}

func (g *LineString) Scan(value interface{}) error {
	return scan(value, g)
}

func (g LineString) Value() (driver.Value, error) {
	return value(&g)
}

func (g LineString) ewkbType() uint32 {
	return ewkbLineStringType
}

type LineStringZ []PointZ

type LineStringZS struct {
	LineStringZ
	srid `json:"srid"`
}

func (g *LineStringZ) ewkbRead(r *ewkbReader) {
	n := r.ReadUint32()
	*g = make([]PointZ, n)
	for i := range *g {
		(*g)[i].ewkbRead(r)
	}
}

func (g *LineStringZ) ewkbWrite(w *ewkbWriter) {
	w.WriteUint32(uint32(len(*g)))
	for _, p := range *g {
		p.ewkbWrite(w)
	}
}

func (g *LineStringZ) Scan(value interface{}) error {
	return scan(value, g)
}

func (g LineStringZ) Value() (driver.Value, error) {
	return value(&g)
}

func (g LineStringZ) ewkbType() uint32 {
	return ewkbLineStringType | ewkbZFlag
}

type LineStringM []PointM

type LineStringMS struct {
	LineStringM
	srid `json:"srid"`
}

func (g *LineStringM) ewkbRead(r *ewkbReader) {
	n := r.ReadUint32()
	*g = make([]PointM, n)
	for i := range *g {
		(*g)[i].ewkbRead(r)
	}
}

func (g *LineStringM) ewkbWrite(w *ewkbWriter) {
	w.WriteUint32(uint32(len(*g)))
	for _, p := range *g {
		p.ewkbWrite(w)
	}
}

func (g *LineStringM) Scan(value interface{}) error {
	return scan(value, g)
}

func (g LineStringM) Value() (driver.Value, error) {
	return value(&g)
}

func (g LineStringM) ewkbType() uint32 {
	return ewkbLineStringType | ewkbMFlag
}

type LineStringZM []PointZM

type LineStringZMS struct {
	LineStringZM
	srid `json:"srid"`
}

func (g *LineStringZM) ewkbRead(r *ewkbReader) {
	n := r.ReadUint32()
	*g = make([]PointZM, n)
	for i := range *g {
		(*g)[i].ewkbRead(r)
	}
}

func (g *LineStringZM) ewkbWrite(w *ewkbWriter) {
	w.WriteUint32(uint32(len(*g)))
	for _, p := range *g {
		p.ewkbWrite(w)
	}
}

func (g *LineStringZM) Scan(value interface{}) error {
	return scan(value, g)
}

func (g LineStringZM) Value() (driver.Value, error) {
	return value(&g)
}

func (g LineStringZM) ewkbType() uint32 {
	return ewkbLineStringType | ewkbZFlag | ewkbMFlag
}
