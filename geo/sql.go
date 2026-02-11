package geo

import (
	"database/sql/driver"
	"encoding/hex"
	"errors"
)

func scan(value interface{}, g geom) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("EWKB scan: value is not byte slice")
	}

	data, err := hex.DecodeString(string(b))
	if err != nil {
		return err
	}

	return Unmarshal(data, g)
}

func value(g geom) (driver.Value, error) {
	data := Marshal(g)
	return hex.EncodeToString(data), nil
}
