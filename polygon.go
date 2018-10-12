package geo

import (
	"database/sql/driver"
)

type Polygon struct {
	Coordinates LineString   `json:"coordinates"`
	Holes       []LineString `json:"holes,omitempty"`
}

type PolygonS struct {
	Polygon
	srid `json:"srid"`
}

func (g *Polygon) ewkbRead(r *ewkbReader) {
	n := r.ReadUint32()
	g.Holes = make([]LineString, n-1)
	g.Coordinates.ewkbRead(r)
	for i := range g.Holes {
		g.Holes[i].ewkbRead(r)
	}
}

func (g *Polygon) ewkbWrite(w *ewkbWriter) {
	w.WriteUint32(uint32(len(g.Holes) + 1))
	g.Coordinates.ewkbWrite(w)
	for i := range g.Holes {
		g.Holes[i].ewkbWrite(w)
	}
}

func (g *Polygon) Scan(value interface{}) error {
	return scan(value, g)
}

func (g Polygon) Value() (driver.Value, error) {
	return value(&g)
}

func (g Polygon) ewkbType() uint32 {
	return ewkbPolygonType
}

type PolygonZ struct {
	Coordinates LineStringZ   `json:"coordinates"`
	Holes       []LineStringZ `json:"holes,omitempty"`
}

type PolygonZS struct {
	PolygonZ
	srid `json:"srid"`
}

func (g *PolygonZ) ewkbRead(r *ewkbReader) {
	n := r.ReadUint32()
	g.Holes = make([]LineStringZ, n-1)
	g.Coordinates.ewkbRead(r)
	for i := range g.Holes {
		g.Holes[i].ewkbRead(r)
	}
}

func (g *PolygonZ) ewkbWrite(w *ewkbWriter) {
	w.WriteUint32(uint32(len(g.Holes) + 1))
	g.Coordinates.ewkbWrite(w)
	for i := range g.Holes {
		g.Holes[i].ewkbWrite(w)
	}
}

func (g *PolygonZ) Scan(value interface{}) error {
	return scan(value, g)
}

func (g PolygonZ) Value() (driver.Value, error) {
	return value(&g)
}

func (g PolygonZ) ewkbType() uint32 {
	return ewkbPolygonType | ewkbZFlag
}

type PolygonM struct {
	Coordinates LineStringM   `json:"coordinates"`
	Holes       []LineStringM `json:"holes,omitempty"`
}

type PolygonMS struct {
	PolygonM
	srid `json:"srid"`
}

func (g *PolygonM) ewkbRead(r *ewkbReader) {
	n := r.ReadUint32()
	g.Holes = make([]LineStringM, n-1)
	g.Coordinates.ewkbRead(r)
	for i := range g.Holes {
		g.Holes[i].ewkbRead(r)
	}
}

func (g *PolygonM) ewkbWrite(w *ewkbWriter) {
	w.WriteUint32(uint32(len(g.Holes) + 1))
	g.Coordinates.ewkbWrite(w)
	for i := range g.Holes {
		g.Holes[i].ewkbWrite(w)
	}
}

func (g *PolygonM) Scan(value interface{}) error {
	return scan(value, g)
}

func (g PolygonM) Value() (driver.Value, error) {
	return value(&g)
}

func (g PolygonM) ewkbType() uint32 {
	return ewkbPolygonType | ewkbMFlag
}

type PolygonZM struct {
	Coordinates LineStringZM   `json:"coordinates"`
	Holes       []LineStringZM `json:"holes,omitempty"`
}

type PolygonZMS struct {
	PolygonZM
	srid `json:"srid"`
}

func (g *PolygonZM) ewkbRead(r *ewkbReader) {
	n := r.ReadUint32()
	g.Holes = make([]LineStringZM, n-1)
	g.Coordinates.ewkbRead(r)
	for i := range g.Holes {
		g.Holes[i].ewkbRead(r)
	}
}

func (g *PolygonZM) ewkbWrite(w *ewkbWriter) {
	w.WriteUint32(uint32(len(g.Holes) + 1))
	g.Coordinates.ewkbWrite(w)
	for i := range g.Holes {
		g.Holes[i].ewkbWrite(w)
	}
}

func (g *PolygonZM) Scan(value interface{}) error {
	return scan(value, g)
}

func (g PolygonZM) Value() (driver.Value, error) {
	return value(&g)
}

func (g PolygonZM) ewkbType() uint32 {
	return ewkbPolygonType | ewkbZFlag | ewkbMFlag
}
