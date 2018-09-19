package geo

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

type geom interface {
	ewkbType() uint32
	ewkbRead(r *ewkbReader)
	ewkbWrite(r *ewkbWriter)
}

type geomS interface {
	geom
	getSRID() uint32
	setSRID(srid uint32)
}

func scan(value interface{}, g geom) error {
	r, err := newEWKBReader(value)
	if err != nil {
		return err
	}

	ewkbType := r.ReadUint32()
	haveSRID := (ewkbType & ewkbSRIDFlag) != 0
	ewkbType &= ^ewkbSRIDFlag
	if ewkbType != g.ewkbType() {
		return errors.New("Incorrect EWKB type")
	}
	fmt.Println("Have SRID", haveSRID)
	var srid uint32
	if haveSRID {
		srid = r.ReadUint32()
	}
	if s, ok := g.(geomS); ok {
		fmt.Println("asdfasdf", haveSRID)
		s.setSRID(srid)
	}

	g.ewkbRead(&r)
	return nil
}

func value(g geom) (driver.Value, error) {
	w := newEWKBWriter()

	if s, ok := g.(geomS); ok {
		w.WriteUint32(g.ewkbType() | ewkbSRIDFlag)
		w.WriteUint32(s.getSRID())
	} else {
		w.WriteUint32(g.ewkbType())
	}
	g.ewkbWrite(&w)

	return w.Value(), nil
}
