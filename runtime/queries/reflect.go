package queries

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/kernelpayments/sqlbunny/runtime/bunny"

	"github.com/kernelpayments/sqlbunny/runtime/strmangle"
	"github.com/kernelpayments/sqlbunny/types/null/convert"
	"github.com/pkg/errors"
)

var (
	bindAccepts = []reflect.Kind{reflect.Ptr, reflect.Slice, reflect.Ptr, reflect.Struct}

	mut         sync.RWMutex
	bindingMaps = make(map[string][]MappedField)
	structMaps  = make(map[string]map[string]MappedField)
)

type MappedField struct {
	Path        uint64
	ParentValid *MappedField
}

// Identifies what kind of object we're binding to
type bindKind int

const (
	kindStruct bindKind = iota
	kindSliceStruct
	kindPtrSliceStruct
)

const (
	loadMethodPrefix       = "Load"
	relationshipStructName = "R"
	loaderStructName       = "L"
)

// Bind executes the query and inserts the
// result into the passed in object pointer
//
// Bind rules:
//   - Struct tags control bind, in the form of: `bunny:"name,bind,null:valid_column_name"`
//   - If the `bunny` struct tag is not present, the field will be completely ignored.
//   - The "name" part specifies the SQL column name that will be bound to this field.
//   - If the ",bind" option is specified on a field of struct type, Bind will recurse into it
//     to look for fields for binding. "name" is appended as a prefix to the SQL column names
//     of the inner fields.
//   - If the ",null:valid_column_name" option is specified in addition to ",bind", the SQL boolean column
//     "valid_column_name" is used to tell whether the nested struct is valid (not null) or not (null).
func Bind(rows *sql.Rows, obj interface{}) error {
	structType, sliceType, singular, err := bindChecks(obj)
	if err != nil {
		return err
	}

	return bind(rows, obj, structType, sliceType, singular)
}

// Bind executes the query and inserts the
// result into the passed in object pointer
//
// See documentation for bunny.Bind()
func (q *Query) Bind(ctx context.Context, obj interface{}) error {
	structType, sliceType, bkind, err := bindChecks(obj)
	if err != nil {
		return err
	}

	rows, err := q.Query(ctx)
	if err != nil {
		return errors.Wrap(err, "bind failed to execute query")
	}
	defer rows.Close()
	if res := bind(rows, obj, structType, sliceType, bkind); res != nil {
		return res
	}

	if len(q.load) != 0 {
		return eagerLoad(ctx, q.load, obj, bkind)
	}

	return nil
}

// bindChecks resolves information about the bind target, and errors if it's not an object
// we can bind to.
func bindChecks(obj interface{}) (structType reflect.Type, sliceType reflect.Type, bkind bindKind, err error) {
	typ := reflect.TypeOf(obj)
	kind := typ.Kind()

	setErr := func() {
		err = errors.Errorf("obj type should be *Type, *[]Type, or *[]*Type but was %q", reflect.TypeOf(obj).String())
	}

	for i := 0; ; i++ {
		switch i {
		case 0:
			if kind != reflect.Ptr {
				setErr()
				return
			}
		case 1:
			switch kind {
			case reflect.Struct:
				structType = typ
				bkind = kindStruct
				return
			case reflect.Slice:
				sliceType = typ
			default:
				setErr()
				return
			}
		case 2:
			switch kind {
			case reflect.Struct:
				structType = typ
				bkind = kindSliceStruct
				return
			case reflect.Ptr:
			default:
				setErr()
				return
			}
		case 3:
			if kind != reflect.Struct {
				setErr()
				return
			}
			structType = typ
			bkind = kindPtrSliceStruct
			return
		}

		typ = typ.Elem()
		kind = typ.Kind()
	}
}

func bind(rows *sql.Rows, obj interface{}, structType, sliceType reflect.Type, bkind bindKind) error {
	cols, err := rows.Columns()
	if err != nil {
		return errors.Wrap(err, "bind failed to get field names")
	}

	var ptrSlice reflect.Value
	switch bkind {
	case kindSliceStruct, kindPtrSliceStruct:
		ptrSlice = reflect.Indirect(reflect.ValueOf(obj))
	}

	var strMapping map[string]MappedField
	var sok bool
	var mapping []MappedField
	var ok bool

	typStr := structType.String()

	mapKey := makeCacheKey(typStr, cols)
	mut.RLock()
	mapping, ok = bindingMaps[mapKey]
	if !ok {
		if strMapping, sok = structMaps[typStr]; !sok {
			strMapping = MakeStructMapping(structType)
		}
	}
	mut.RUnlock()

	if !ok {
		mapping, err = BindMapping(structType, strMapping, cols)
		if err != nil {
			return err
		}

		mut.Lock()
		if !sok {
			structMaps[typStr] = strMapping
		}
		bindingMaps[mapKey] = mapping
		mut.Unlock()
	}

	var oneStruct reflect.Value
	if bkind == kindSliceStruct {
		oneStruct = reflect.Indirect(reflect.New(structType))
	}

	foundOne := false
	for rows.Next() {
		if bkind == kindStruct && foundOne {
			return bunny.ErrMultipleRows
		}

		foundOne = true
		var newStruct reflect.Value
		var pointers []interface{}

		switch bkind {
		case kindStruct:
			pointers = PtrsFromMapping(reflect.Indirect(reflect.ValueOf(obj)), mapping)
		case kindSliceStruct:
			pointers = PtrsFromMapping(oneStruct, mapping)
		case kindPtrSliceStruct:
			newStruct = reflect.New(structType)
			pointers = PtrsFromMapping(reflect.Indirect(newStruct), mapping)
		}
		if err != nil {
			return err
		}

		if err := rows.Scan(pointers...); err != nil {
			return errors.Wrap(err, "failed to bind pointers to obj")
		}

		switch bkind {
		case kindSliceStruct:
			ptrSlice.Set(reflect.Append(ptrSlice, oneStruct))
		case kindPtrSliceStruct:
			ptrSlice.Set(reflect.Append(ptrSlice, newStruct))
		}
	}

	if bkind == kindStruct && !foundOne {
		return sql.ErrNoRows
	}

	return nil
}

// BindMapping creates a mapping that helps look up the pointer for the
// field given.
func BindMapping(typ reflect.Type, mapping map[string]MappedField, cols []string) ([]MappedField, error) {
	ptrs := make([]MappedField, len(cols))

ColLoop:
	for i, name := range cols {
		ptrMap, ok := mapping[name]
		if ok {
			ptrs[i] = ptrMap
			continue
		}

		suffix := "." + name
		for maybeMatch, mapping := range mapping {
			if strings.HasSuffix(maybeMatch, suffix) {
				ptrs[i] = mapping
				continue ColLoop
			}
		}
		// if c doesn't exist in the model, the pointer will be the zero value in the ptrs array and it's value will be thrown away
		continue
	}

	return ptrs, nil
}

// PtrsFromMapping expects to be passed an addressable struct and a mapping
// of where to find things. It pulls the pointers out referred to by the mapping.
func PtrsFromMapping(val reflect.Value, mapping []MappedField) []interface{} {
	ptrs := make([]interface{}, len(mapping))
	for i, m := range mapping {
		ptrs[i] = ptrFromMapping(val, m, true)
	}
	return ptrs
}

// ValuesFromMapping expects to be passed an addressable struct and a mapping
// of where to find things. It pulls the pointers out referred to by the mapping.
func ValuesFromMapping(val reflect.Value, mapping []MappedField) []interface{} {
	ptrs := make([]interface{}, len(mapping))
	for i, m := range mapping {
		ptrs[i] = ptrFromMapping(val, m, false)
	}
	return ptrs
}

type ignoreNullScan struct {
	dest interface{}
}

// Scan implements the Scanner interface.
func (v *ignoreNullScan) Scan(value interface{}) error {
	if value == nil {
		return convert.ConvertAssignNil(v.dest)
	}
	return convert.ConvertAssign(v.dest, value)
}

// ptrFromMapping expects to be passed an addressable struct that it's looking
// for things on.
func ptrFromMapping(val reflect.Value, mapping MappedField, addressOf bool) interface{} {
	if mapping.Path == 0 {
		var ignored interface{}
		return &ignored
	}

	if !addressOf && mapping.ParentValid != nil {
		valid := ptrFromMapping(val, *mapping.ParentValid, false)
		if valid != true {
			var nothing interface{}
			return &nothing
		}
	}

	for i := 0; i < 8; i++ {
		v := (mapping.Path >> uint(i*8)) & 0xFF

		if v == 0 {
			if addressOf && val.Kind() != reflect.Ptr {
				val = val.Addr()
			}
			if !addressOf && val.Kind() == reflect.Ptr {
				val = reflect.Indirect(val)
			}

			if addressOf && mapping.ParentValid != nil {
				// When scanning into a field that's child of a nullable struct,
				// we use a special scan variant that converts DB nulls to
				// Go zero values instead of erroring (unless the field
				// implements sql.Scanner, in which case it's used as usual)
				return &ignoreNullScan{
					dest: val.Interface(),
				}
			}
			return val.Interface()
		}

		val = val.Field(int(v - 1))
		if val.Kind() == reflect.Ptr {
			val = reflect.Indirect(val)
			if !val.IsValid() {
				var nothing interface{}
				return reflect.ValueOf(&nothing)
			}
		}
	}

	panic("could not find pointer from mapping")
}

// MakeStructMapping creates a map of the struct to be able to quickly look
// up its pointers and values by name.
func MakeStructMapping(typ reflect.Type) map[string]MappedField {
	fieldMaps := make(map[string]MappedField)
	makeStructMappingHelper(typ, "", MappedField{}, 0, fieldMaps)

	return fieldMaps
}

func makeStructMappingHelper(typ reflect.Type, prefix string, current MappedField, depth uint, fieldMaps map[string]MappedField) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	n := typ.NumField()
	for i := 0; i < n; i++ {
		f := typ.Field(i)

		tag, err := getBunnyTag(f)
		if err != nil {
			panic(err)
		}
		if !tag.present {
			continue
		}

		name := tag.name
		if len(prefix) != 0 {
			name = prefix + name
		}

		if tag.bind {
			if len(tag.null) != 0 {
				// TODO autodiscover this
				structFieldIdx := 0
				validFieldIdx := 1

				valid := MappedField{
					Path:        current.Path | uint64(i+1)<<depth | (uint64(validFieldIdx+1) << (depth + 8)),
					ParentValid: current.ParentValid,
				}

				fieldMaps[prefix+tag.null] = valid
				next := MappedField{
					Path:        current.Path | uint64(i+1)<<depth | (uint64(structFieldIdx+1) << (depth + 8)),
					ParentValid: &valid,
				}
				makeStructMappingHelper(f.Type.Field(structFieldIdx).Type, name, next, depth+16, fieldMaps)
			} else {
				next := MappedField{
					Path:        current.Path | uint64(i+1)<<depth,
					ParentValid: current.ParentValid,
				}
				makeStructMappingHelper(f.Type, name, next, depth+8, fieldMaps)
			}
			continue
		}

		fieldMaps[name] = MappedField{
			Path:        current.Path | uint64(i+1)<<depth,
			ParentValid: current.ParentValid,
		}
	}
}

type bunnyTag struct {
	present bool
	name    string
	bind    bool
	null    string
}

func getBunnyTag(field reflect.StructField) (bunnyTag, error) {
	tag := field.Tag.Get("bunny")

	// If there is no bunny tag, don't use this field.
	if len(tag) == 0 {
		return bunnyTag{
			present: false,
		}, nil
	}

	parts := strings.Split(tag, ",")
	res := bunnyTag{
		present: true,
		name:    parts[0],
	}
	for _, flag := range parts[1:] {
		if flag == "bind" {
			res.bind = true
		} else if strings.HasPrefix(flag, "null:") {
			res.null = strings.TrimPrefix(flag, "null:")
		} else {
			return bunnyTag{}, fmt.Errorf("Invalid flag in bunny tag in field '%s': '%s'", field.Name, flag)
		}
	}

	if len(res.null) != 0 && !res.bind {
		return bunnyTag{}, fmt.Errorf("Invalid flags in bunny tag in field '%s': null requires bind to be set", field.Name)
	}

	return res, nil
}

func makeCacheKey(typ string, cols []string) string {
	buf := strmangle.GetBuffer()
	buf.WriteString(typ)
	for _, s := range cols {
		buf.WriteString(s)
		buf.WriteByte(',')
	}
	mapKey := buf.String()
	strmangle.PutBuffer(buf)

	return mapKey
}
