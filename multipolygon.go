package geo

import (
	"database/sql/driver"
)

type MultiPolygon []Polygon

type MultiPolygonS struct {
	MultiPolygon
	srid `json:"srid"`
}

func (g *MultiPolygon) ewkbRead(r *ewkbReader) {
	n := r.ReadUint32()
	*g = make([]Polygon, n)
	for i := range *g {
		UnmarshalReader(r.r, &(*g)[i])
	}
}

func (g *MultiPolygon) ewkbWrite(w *ewkbWriter) {
	w.WriteUint32(uint32(len(*g)))
	for _, p := range *g {
		MarshalWriter(w.w, &p)
	}
}

func (g *MultiPolygon) Scan(value interface{}) error {
	return scan(value, g)
}

func (g *MultiPolygonS) Scan(value interface{}) error {
	return scan(value, g)
}

func (g MultiPolygon) Value() (driver.Value, error) {
	return value(&g)
}

func (g MultiPolygon) ewkbType() uint32 {
	return ewkbMultiPolygonType
}

type MultiPolygonZ []PolygonZ

type MultiPolygonZS struct {
	MultiPolygonZ
	srid `json:"srid"`
}

func (g *MultiPolygonZ) ewkbRead(r *ewkbReader) {
	n := r.ReadUint32()
	*g = make([]PolygonZ, n)
	for i := range *g {
		UnmarshalReader(r.r, &(*g)[i])
	}
}

func (g *MultiPolygonZ) ewkbWrite(w *ewkbWriter) {
	w.WriteUint32(uint32(len(*g)))
	for _, p := range *g {
		MarshalWriter(w.w, &p)
	}
}

func (g *MultiPolygonZ) Scan(value interface{}) error {
	return scan(value, g)
}

func (g *MultiPolygonZS) Scan(value interface{}) error {
	return scan(value, g)
}

func (g MultiPolygonZ) Value() (driver.Value, error) {
	return value(&g)
}

func (g MultiPolygonZ) ewkbType() uint32 {
	return ewkbMultiPolygonType | ewkbZFlag
}

type MultiPolygonM []PolygonM

type MultiPolygonMS struct {
	MultiPolygonM
	srid `json:"srid"`
}

func (g *MultiPolygonM) ewkbRead(r *ewkbReader) {
	n := r.ReadUint32()
	*g = make([]PolygonM, n)
	for i := range *g {
		UnmarshalReader(r.r, &(*g)[i])
	}
}

func (g *MultiPolygonM) ewkbWrite(w *ewkbWriter) {
	w.WriteUint32(uint32(len(*g)))
	for _, p := range *g {
		MarshalWriter(w.w, &p)
	}
}

func (g *MultiPolygonM) Scan(value interface{}) error {
	return scan(value, g)
}

func (g *MultiPolygonMS) Scan(value interface{}) error {
	return scan(value, g)
}

func (g MultiPolygonM) Value() (driver.Value, error) {
	return value(&g)
}

func (g MultiPolygonM) ewkbType() uint32 {
	return ewkbMultiPolygonType | ewkbMFlag
}

type MultiPolygonZM []PolygonZM

type MultiPolygonZMS struct {
	MultiPolygonZM
	srid `json:"srid"`
}

func (g *MultiPolygonZM) ewkbRead(r *ewkbReader) {
	n := r.ReadUint32()
	*g = make([]PolygonZM, n)
	for i := range *g {
		UnmarshalReader(r.r, &(*g)[i])
	}
}

func (g *MultiPolygonZM) ewkbWrite(w *ewkbWriter) {
	w.WriteUint32(uint32(len(*g)))
	for _, p := range *g {
		MarshalWriter(w.w, &p)
	}
}

func (g *MultiPolygonZM) Scan(value interface{}) error {
	return scan(value, g)
}

func (g *MultiPolygonZMS) Scan(value interface{}) error {
	return scan(value, g)
}

func (g MultiPolygonZM) Value() (driver.Value, error) {
	return value(&g)
}

func (g MultiPolygonZM) ewkbType() uint32 {
	return ewkbMultiPolygonType | ewkbZFlag | ewkbMFlag
}
