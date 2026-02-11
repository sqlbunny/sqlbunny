package geo

import (
	"database/sql/driver"
)

type MultiPoint []Point

type MultiPointS struct {
	MultiPoint
	srid `json:"srid"`
}

func (g *MultiPoint) ewkbRead(r *ewkbReader) {
	n := r.ReadUint32()
	*g = make([]Point, n)
	for i := range *g {
		UnmarshalReader(r.r, &(*g)[i])
	}
}

func (g *MultiPoint) ewkbWrite(w *ewkbWriter) {
	w.WriteUint32(uint32(len(*g)))
	for _, p := range *g {
		MarshalWriter(w.w, &p)
	}
}

func (g *MultiPoint) Scan(value interface{}) error {
	return scan(value, g)
}

func (g MultiPoint) Value() (driver.Value, error) {
	return value(&g)
}

func (g MultiPoint) ewkbType() uint32 {
	return ewkbMultiPointType
}

type MultiPointZ []PointZ

type MultiPointZS struct {
	MultiPointZ
	srid `json:"srid"`
}

func (g *MultiPointZ) ewkbRead(r *ewkbReader) {
	n := r.ReadUint32()
	*g = make([]PointZ, n)
	for i := range *g {
		UnmarshalReader(r.r, &(*g)[i])
	}
}

func (g *MultiPointZ) ewkbWrite(w *ewkbWriter) {
	w.WriteUint32(uint32(len(*g)))
	for _, p := range *g {
		MarshalWriter(w.w, &p)
	}
}

func (g *MultiPointZ) Scan(value interface{}) error {
	return scan(value, g)
}

func (g MultiPointZ) Value() (driver.Value, error) {
	return value(&g)
}

func (g MultiPointZ) ewkbType() uint32 {
	return ewkbMultiPointType | ewkbZFlag
}

type MultiPointM []PointM

type MultiPointMS struct {
	MultiPointM
	srid `json:"srid"`
}

func (g *MultiPointM) ewkbRead(r *ewkbReader) {
	n := r.ReadUint32()
	*g = make([]PointM, n)
	for i := range *g {
		UnmarshalReader(r.r, &(*g)[i])
	}
}

func (g *MultiPointM) ewkbWrite(w *ewkbWriter) {
	w.WriteUint32(uint32(len(*g)))
	for _, p := range *g {
		MarshalWriter(w.w, &p)
	}
}

func (g *MultiPointM) Scan(value interface{}) error {
	return scan(value, g)
}

func (g MultiPointM) Value() (driver.Value, error) {
	return value(&g)
}

func (g MultiPointM) ewkbType() uint32 {
	return ewkbMultiPointType | ewkbMFlag
}

type MultiPointZM []PointZM

type MultiPointZMS struct {
	MultiPointZM
	srid `json:"srid"`
}

func (g *MultiPointZM) ewkbRead(r *ewkbReader) {
	n := r.ReadUint32()
	*g = make([]PointZM, n)
	for i := range *g {
		UnmarshalReader(r.r, &(*g)[i])
	}
}

func (g *MultiPointZM) ewkbWrite(w *ewkbWriter) {
	w.WriteUint32(uint32(len(*g)))
	for _, p := range *g {
		MarshalWriter(w.w, &p)
	}
}

func (g *MultiPointZM) Scan(value interface{}) error {
	return scan(value, g)
}

func (g MultiPointZM) Value() (driver.Value, error) {
	return value(&g)
}

func (g MultiPointZM) ewkbType() uint32 {
	return ewkbMultiPointType | ewkbZFlag | ewkbMFlag
}
