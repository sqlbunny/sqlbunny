{{- $modelName := .IDType.Name | titleCase -}}

type {{$modelName}}Array []{{$modelName}}

// Scan implements the sql.Scanner interface.
func (a *{{$modelName}}Array) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case string:
		return a.scanBytes([]byte(src))
	case nil:
		*a = nil
		return nil
	}

	return fmt.Errorf("pq: cannot convert %T to ByteaArray", src)
}

func (a *{{$modelName}}Array) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "{{$modelName}}Array")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make({{$modelName}}Array, len(elems))
		for i, v := range elems {
			bytes, err := parseBytea(v)
			if err != nil {
				return fmt.Errorf("could not parse id array index %d: %s", i, err.Error())
			}
			if len(bytes) != bunny.IDRawLen {
				return fmt.Errorf("could not parse id array index %d: got len %d, expected %d", i, len(bytes), bunny.IDRawLen)
			}
			copy(b[i].IDData[:], bytes[:])
		}
		*a = b
	}
	return nil
}

// Value implements the driver.Valuer interface. It uses the "hex" format which
// is only supported on PostgreSQL 9.0 or newer.
func (a {{$modelName}}Array) Value() (driver.Value, error) {
	if a == nil {
		//return nil, nil
		return "{}", nil
	}

	if n := len(a); n > 0 {
		// There will be at least two curly brackets, 2*N bytes of quotes,
		// 3*N bytes of hex formatting, and N-1 bytes of delimiters.
		size := 1 + 6*n + hex.EncodedLen(bunny.IDRawLen)*len(a)
		b := make([]byte, size)

		for i, s := 0, b; i < n; i++ {
			o := copy(s, `,"\\x`)
			o += hex.Encode(s[o:], a[i].IDData[:])
			s[o] = '"'
			s = s[o+1:]
		}

		b[0] = '{'
		b[size-1] = '}'

		return string(b), nil
	}

	return "{}", nil
}
