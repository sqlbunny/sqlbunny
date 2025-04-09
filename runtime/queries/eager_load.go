package queries

import (
	"context"
	"reflect"
	"strings"

	"github.com/sqlbunny/errors"
	"github.com/sqlbunny/sqlbunny/runtime/strmangle"
)

type loadRelationshipState struct {
	ctx    context.Context
	loaded map[string]struct{}
	toLoad []string
}

func (l loadRelationshipState) hasLoaded(depth int) bool {
	_, ok := l.loaded[l.buildKey(depth)]
	return ok
}

func (l loadRelationshipState) setLoaded(depth int) {
	l.loaded[l.buildKey(depth)] = struct{}{}
}

func (l loadRelationshipState) buildKey(depth int) string {
	buf := strmangle.GetBuffer()

	for i, piece := range l.toLoad[:depth+1] {
		if i != 0 {
			buf.WriteByte('.')
		}
		buf.WriteString(piece)
	}

	str := buf.String()
	strmangle.PutBuffer(buf)
	return str
}

// eagerLoad loads all of the model's relationships
//
// toLoad should look like:
// []string{"Relationship", "Relationship.NestedRelationship"} ... etc
// obj should be one of:
// *[]*struct or *struct
// bkind should reflect what kind of thing it is above
func eagerLoad(ctx context.Context, toLoad []string, obj any, bkind bindKind) error {
	state := loadRelationshipState{
		ctx:    ctx,
		loaded: map[string]struct{}{},
	}

	val := reflect.ValueOf(obj)
	if bkind == kindStruct {
		r := reflect.MakeSlice(reflect.SliceOf(val.Type()), 1, 1)
		r.Index(0).Set(val)
		val = r
	} else {
		val = val.Elem()
	}

	for _, toLoad := range toLoad {
		state.toLoad = strings.Split(toLoad, ".")
		for i := range state.toLoad {
			state.toLoad[i] = strmangle.TitleCase(state.toLoad[i])
		}
		if err := state.loadRelationships(0, val); err != nil {
			return err
		}
	}

	return nil
}

// loadRelationships dynamically calls the template generated eager load
// functions of the form:
//
//	func (l ModelL) LoadRelationshipName(ctx context.Context, slice []*Model) error
//
// The arguments to this function are:
//   - l is not used, and it is always passed the zero value.
//   - ctx is used to perform additional queries that might be required for loading the relationships.
//   - slice is the slice of model instances, always of the type []*Model.
//
// We start with a normal select before eager loading anything: select * from a;
// Then we start eager loading things, it can be represented by a DAG
//
//	    a1, a2           select id, a_id from b where id in (a1, a2)
//	   / |    \
//	  b1 b2    b3        select id, b_id from c where id in (b2, b3, b4)
//	 /   | \     \
//	c1  c2 c3    c4
//
// That's to say that we descend the graph of relationships, and at each level
// we gather all the things up we want to load into, load them, and then move
// to the next level of the graph.
func (l loadRelationshipState) loadRelationships(depth int, loadingFrom reflect.Value) error {
	if loadingFrom.Len() == 0 {
		return nil
	}

	if !l.hasLoaded(depth) {
		if err := l.callLoadFunction(depth, loadingFrom); err != nil {
			return err
		}
	}

	// Check if we can stop
	if depth+1 >= len(l.toLoad) {
		return nil
	}

	// *[]*struct -> []*struct
	// *struct -> struct
	loadingFrom = reflect.Indirect(loadingFrom)

	// Collect eagerly loaded things to send into next eager load call
	slice, err := collectLoaded(l.toLoad[depth], loadingFrom)
	if err != nil {
		return err
	}

	// If we could collect nothing we're done
	if slice.Len() == 0 {
		return nil
	}

	return l.loadRelationships(depth+1, slice)
}

// callLoadFunction finds the loader struct, finds the method that we need
// to call and calls it.
func (l loadRelationshipState) callLoadFunction(depth int, loadingFrom reflect.Value) error {
	current := l.toLoad[depth]
	sliceType := loadingFrom.Type()
	modelType := sliceType.Elem().Elem()
	ln, found := modelType.FieldByName(loaderStructName)
	// It's possible a Loaders struct doesn't exist on the struct.
	if !found {
		return errors.Errorf("attempted to load %s but no L struct was found", current)
	}

	// Attempt to find the LoadRelationshipName function
	loadMethod, found := ln.Type.MethodByName(loadMethodPrefix + current)
	if !found {
		return errors.Errorf("could not find %s%s method for eager loading", loadMethodPrefix, current)
	}

	methodArgs := []reflect.Value{
		reflect.Zero(ln.Type),
		reflect.ValueOf(l.ctx),
		loadingFrom,
	}

	ret := loadMethod.Func.Call(methodArgs)
	if intf := ret[0].Interface(); intf != nil {
		return errors.Errorf("failed to eager load %s: %w", current, intf.(error))
	}

	l.setLoaded(depth)
	return nil
}

// collectLoaded traverses the next level of the graph and picks up all
// the values that we need for the next eager load query.
//
// For example when loadingFrom is [parent1, parent2]
//
//	parent1 -> child1
//	       \-> child2
//	parent2 -> child3
//
// This should return [child1, child2, child3]
func collectLoaded(key string, loadingFrom reflect.Value) (reflect.Value, error) {
	currentModelType := loadingFrom.Type().Elem().Elem()
	f, ok := currentModelType.FieldByName(relationshipStructName)
	if !ok {
		return reflect.Value{}, errors.New("relationship struct was not found")
	}
	relFieldIndex := f.Index
	f, ok = f.Type.Elem().FieldByName(key)
	if !ok {
		return reflect.Value{}, errors.New("field was not found in relationship struct")
	}
	keyFieldIndex := f.Index

	// Ensure that we get rid of all the helper "XSlice" types
	toMany := f.Type.Kind() == reflect.Slice
	var nextModelType reflect.Type
	if toMany {
		nextModelType = f.Type.Elem().Elem()
	} else {
		nextModelType = f.Type.Elem()
	}
	nextSliceType := reflect.SliceOf(reflect.PtrTo(nextModelType))

	collection := reflect.MakeSlice(nextSliceType, 0, 0)

	lnFrom := loadingFrom.Len()
	for i := 0; i < lnFrom; i++ {
		o := loadingFrom.Index(i).Elem().FieldByIndex(relFieldIndex).Elem().FieldByIndex(keyFieldIndex)
		if o.IsNil() {
			continue
		}
		if toMany {
			collection = reflect.AppendSlice(collection, o)
		} else {
			collection = reflect.Append(collection, o)
		}
	}

	return collection, nil
}
