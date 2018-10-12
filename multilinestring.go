package geo

import (
	"database/sql/driver"
)

type MultiLineString []LineString

type MultiLineStringS struct {
	MultiLineString
	srid `json:"srid"`
}

func (g *MultiLineString) ewkbRead(r *ewkbReader) {
	n := r.ReadUint32()
	*g = make([]LineString, n)
	for i := range *g {
		UnmarshalReader(r.r, &(*g)[i])
	}
}

func (g *MultiLineString) ewkbWrite(w *ewkbWriter) {
	w.WriteUint32(uint32(len(*g)))
	for _, p := range *g {
		MarshalWriter(w.w, &p)
	}
}

func (g *MultiLineString) Scan(value interface{}) error {
	return scan(value, g)
}

func (g MultiLineString) Value() (driver.Value, error) {
	return value(&g)
}

func (g MultiLineString) ewkbType() uint32 {
	return ewkbMultiLineStringType
}

type MultiLineStringZ []LineStringZ

type MultiLineStringZS struct {
	MultiLineStringZ
	srid `json:"srid"`
}

func (g *MultiLineStringZ) ewkbRead(r *ewkbReader) {
	n := r.ReadUint32()
	*g = make([]LineStringZ, n)
	for i := range *g {
		UnmarshalReader(r.r, &(*g)[i])
	}
}

func (g *MultiLineStringZ) ewkbWrite(w *ewkbWriter) {
	w.WriteUint32(uint32(len(*g)))
	for _, p := range *g {
		MarshalWriter(w.w, &p)
	}
}

func (g *MultiLineStringZ) Scan(value interface{}) error {
	return scan(value, g)
}

func (g MultiLineStringZ) Value() (driver.Value, error) {
	return value(&g)
}

func (g MultiLineStringZ) ewkbType() uint32 {
	return ewkbMultiLineStringType | ewkbZFlag
}

type MultiLineStringM []LineStringM

type MultiLineStringMS struct {
	MultiLineStringM
	srid `json:"srid"`
}

func (g *MultiLineStringM) ewkbRead(r *ewkbReader) {
	n := r.ReadUint32()
	*g = make([]LineStringM, n)
	for i := range *g {
		UnmarshalReader(r.r, &(*g)[i])
	}
}

func (g *MultiLineStringM) ewkbWrite(w *ewkbWriter) {
	w.WriteUint32(uint32(len(*g)))
	for _, p := range *g {
		MarshalWriter(w.w, &p)
	}
}

func (g *MultiLineStringM) Scan(value interface{}) error {
	return scan(value, g)
}

func (g MultiLineStringM) Value() (driver.Value, error) {
	return value(&g)
}

func (g MultiLineStringM) ewkbType() uint32 {
	return ewkbMultiLineStringType | ewkbMFlag
}

type MultiLineStringZM []LineStringZM

type MultiLineStringZMS struct {
	MultiLineStringZM
	srid `json:"srid"`
}

func (g *MultiLineStringZM) ewkbRead(r *ewkbReader) {
	n := r.ReadUint32()
	*g = make([]LineStringZM, n)
	for i := range *g {
		UnmarshalReader(r.r, &(*g)[i])
	}
}

func (g *MultiLineStringZM) ewkbWrite(w *ewkbWriter) {
	w.WriteUint32(uint32(len(*g)))
	for _, p := range *g {
		MarshalWriter(w.w, &p)
	}
}

func (g *MultiLineStringZM) Scan(value interface{}) error {
	return scan(value, g)
}

func (g MultiLineStringZM) Value() (driver.Value, error) {
	return value(&g)
}

func (g MultiLineStringZM) ewkbType() uint32 {
	return ewkbMultiLineStringType | ewkbZFlag | ewkbMFlag
}
