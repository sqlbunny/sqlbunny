package queries

import (
	"context"
	"database/sql/driver"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/KernelPay/sqlbunny/bunny"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func execToContext(exec bunny.Executor) context.Context {
	return bunny.WithExecutor(context.Background(), exec)
}
func bin64(i uint64) string {
	str := strconv.FormatUint(i, 2)
	pad := 64 - len(str)
	if pad > 0 {
		str = strings.Repeat("0", pad) + str
	}

	var newStr string
	for i := 0; i < len(str); i += 8 {
		if i != 0 {
			newStr += " "
		}
		newStr += str[i : i+8]
	}

	return newStr
}

type mockRowMaker struct {
	int
	rows []driver.Value
}

func TestBindStruct(t *testing.T) {
	t.Parallel()

	testResults := struct {
		ID   int    `bunny:"id"`
		Name string `bunny:"test"`
	}{}

	query := &Query{
		from:    []string{"fun"},
		dialect: &Dialect{LQ: '"', RQ: '"', IndexPlaceholders: true},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	ret := sqlmock.NewRows([]string{"id", "test"})
	ret.AddRow(driver.Value(int64(35)), driver.Value("pat"))
	mock.ExpectQuery(`SELECT \* FROM "fun";`).WillReturnRows(ret)

	SetContext(query, execToContext(db))
	err = query.Bind(&testResults)
	if err != nil {
		t.Error(err)
	}

	if id := testResults.ID; id != 35 {
		t.Error("wrong ID:", id)
	}
	if name := testResults.Name; name != "pat" {
		t.Error("wrong name:", name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestBindSlice(t *testing.T) {
	t.Parallel()

	testResults := []struct {
		ID   int    `bunny:"id"`
		Name string `bunny:"test"`
	}{}

	query := &Query{
		from:    []string{"fun"},
		dialect: &Dialect{LQ: '"', RQ: '"', IndexPlaceholders: true},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	ret := sqlmock.NewRows([]string{"id", "test"})
	ret.AddRow(driver.Value(int64(35)), driver.Value("pat"))
	ret.AddRow(driver.Value(int64(12)), driver.Value("cat"))
	mock.ExpectQuery(`SELECT \* FROM "fun";`).WillReturnRows(ret)

	SetContext(query, execToContext(db))
	err = query.Bind(&testResults)
	if err != nil {
		t.Error(err)
	}

	if len(testResults) != 2 {
		t.Fatal("wrong number of results:", len(testResults))
	}
	if id := testResults[0].ID; id != 35 {
		t.Error("wrong ID:", id)
	}
	if name := testResults[0].Name; name != "pat" {
		t.Error("wrong name:", name)
	}

	if id := testResults[1].ID; id != 12 {
		t.Error("wrong ID:", id)
	}
	if name := testResults[1].Name; name != "cat" {
		t.Error("wrong name:", name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestBindPtrSlice(t *testing.T) {
	t.Parallel()

	testResults := []*struct {
		ID   int    `bunny:"id"`
		Name string `bunny:"test"`
	}{}

	query := &Query{
		from:    []string{"fun"},
		dialect: &Dialect{LQ: '"', RQ: '"', IndexPlaceholders: true},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	ret := sqlmock.NewRows([]string{"id", "test"})
	ret.AddRow(driver.Value(int64(35)), driver.Value("pat"))
	ret.AddRow(driver.Value(int64(12)), driver.Value("cat"))
	mock.ExpectQuery(`SELECT \* FROM "fun";`).WillReturnRows(ret)

	SetContext(query, execToContext(db))
	err = query.Bind(&testResults)
	if err != nil {
		t.Error(err)
	}

	if len(testResults) != 2 {
		t.Fatal("wrong number of results:", len(testResults))
	}
	if id := testResults[0].ID; id != 35 {
		t.Error("wrong ID:", id)
	}
	if name := testResults[0].Name; name != "pat" {
		t.Error("wrong name:", name)
	}

	if id := testResults[1].ID; id != 12 {
		t.Error("wrong ID:", id)
	}
	if name := testResults[1].Name; name != "cat" {
		t.Error("wrong name:", name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func testMakeMapping(byt ...byte) uint64 {
	var x uint64
	for i, b := range byt {
		x |= uint64(b) << (uint(i) * 8)
	}
	x |= uint64(255) << uint(len(byt)*8)
	return x
}

func TestMakeStructMapping(t *testing.T) {
	t.Parallel()

	var testStruct = struct {
		LastName    string `bunny:"different"`
		AwesomeName string `bunny:"awesome_name"`
		Face        string
		Nose        string `bunny:"nose"`

		Nested struct {
			LastName    string `bunny:"different"`
			AwesomeName string `bunny:"awesome_name"`
			Face        string
			Nose        string `bunny:"nose"`

			Nested2 struct {
				Nose string `bunny:"nose"`
			} `bunny:"nested2,bind"`
			Nested3 struct {
				Nose string `bunny:"nose"`
			} `bunny:"nested3,structbind"`
		} `bunny:"nested,bind"`
	}{}

	got := MakeStructMapping(reflect.TypeOf(testStruct))

	expectMap := map[string]uint64{
		"different":            testMakeMapping(0),
		"awesome_name":         testMakeMapping(1),
		"nose":                 testMakeMapping(3),
		"nested.different":     testMakeMapping(4, 0),
		"nested.awesome_name":  testMakeMapping(4, 1),
		"nested.nose":          testMakeMapping(4, 3),
		"nested.nested2.nose":  testMakeMapping(4, 4, 0),
		"nested.nested3__nose": testMakeMapping(4, 5, 0),
	}

	for expName, expVal := range expectMap {
		gotVal, ok := got[expName]
		if !ok {
			t.Errorf("%s) had no value", expName)
			continue
		}

		if gotVal != expVal {
			t.Errorf("%s) wrong value,\nwant: %x (%s)\ngot:  %x (%s)", expName, expVal, bin64(expVal), gotVal, bin64(gotVal))
		}
	}
}

func TestPtrFromMapping(t *testing.T) {
	t.Parallel()

	type NestedPtrs struct {
		Int         int
		IntP        *int
		NestedPtrsP *NestedPtrs
	}

	val := &NestedPtrs{
		Int:  5,
		IntP: new(int),
		NestedPtrsP: &NestedPtrs{
			Int:  6,
			IntP: new(int),
		},
	}

	v := ptrFromMapping(reflect.Indirect(reflect.ValueOf(val)), testMakeMapping(0), true)
	if got := *v.Interface().(*int); got != 5 {
		t.Error("flat int was wrong:", got)
	}
	v = ptrFromMapping(reflect.Indirect(reflect.ValueOf(val)), testMakeMapping(1), true)
	if got := *v.Interface().(*int); got != 0 {
		t.Error("flat pointer was wrong:", got)
	}
	v = ptrFromMapping(reflect.Indirect(reflect.ValueOf(val)), testMakeMapping(2, 0), true)
	if got := *v.Interface().(*int); got != 6 {
		t.Error("nested int was wrong:", got)
	}
	v = ptrFromMapping(reflect.Indirect(reflect.ValueOf(val)), testMakeMapping(2, 1), true)
	if got := *v.Interface().(*int); got != 0 {
		t.Error("nested pointer was wrong:", got)
	}
}

func TestValuesFromMapping(t *testing.T) {
	t.Parallel()

	type NestedPtrs struct {
		Int         int
		IntP        *int
		NestedPtrsP *NestedPtrs
	}

	val := &NestedPtrs{
		Int:  5,
		IntP: new(int),
		NestedPtrsP: &NestedPtrs{
			Int:  6,
			IntP: new(int),
		},
	}
	mapping := []uint64{testMakeMapping(0), testMakeMapping(1), testMakeMapping(2, 0), testMakeMapping(2, 1), 0}
	v := ValuesFromMapping(reflect.Indirect(reflect.ValueOf(val)), mapping)

	if got := v[0].(int); got != 5 {
		t.Error("flat int was wrong:", got)
	}
	if got := v[1].(int); got != 0 {
		t.Error("flat pointer was wrong:", got)
	}
	if got := v[2].(int); got != 6 {
		t.Error("nested int was wrong:", got)
	}
	if got := v[3].(int); got != 0 {
		t.Error("nested pointer was wrong:", got)
	}
	if got := *v[4].(*interface{}); got != nil {
		t.Error("nil pointer was not be ignored:", got)
	}
}

func TestPtrsFromMapping(t *testing.T) {
	t.Parallel()

	type NestedPtrs struct {
		Int         int
		IntP        *int
		NestedPtrsP *NestedPtrs
	}

	val := &NestedPtrs{
		Int:  5,
		IntP: new(int),
		NestedPtrsP: &NestedPtrs{
			Int:  6,
			IntP: new(int),
		},
	}

	mapping := []uint64{testMakeMapping(0), testMakeMapping(1), testMakeMapping(2, 0), testMakeMapping(2, 1)}
	v := PtrsFromMapping(reflect.Indirect(reflect.ValueOf(val)), mapping)

	if got := *v[0].(*int); got != 5 {
		t.Error("flat int was wrong:", got)
	}
	if got := *v[1].(*int); got != 0 {
		t.Error("flat pointer was wrong:", got)
	}
	if got := *v[2].(*int); got != 6 {
		t.Error("nested int was wrong:", got)
	}
	if got := *v[3].(*int); got != 0 {
		t.Error("nested pointer was wrong:", got)
	}
}

func TestGetBunnyTag(t *testing.T) {
	t.Parallel()

	type TestStruct struct {
		FirstName   string `bunny:"test_one,bind"`
		LastName    string `bunny:"test_two"`
		MiddleName  string `bunny:"middle_name,bind"`
		AwesomeName string `bunny:"awesome_name,structbind"`
		Age         string `bunny:"age,bind,structbind"`
		Nose        string
		Fail        string `bunny:"fail,invalidflag"`
	}

	var structFields []reflect.StructField
	typ := reflect.TypeOf(TestStruct{})
	removeOk := func(thing reflect.StructField, ok bool) reflect.StructField {
		if !ok {
			panic("Exploded")
		}
		return thing
	}
	structFields = append(structFields, removeOk(typ.FieldByName("FirstName")))
	structFields = append(structFields, removeOk(typ.FieldByName("LastName")))
	structFields = append(structFields, removeOk(typ.FieldByName("MiddleName")))
	structFields = append(structFields, removeOk(typ.FieldByName("AwesomeName")))
	structFields = append(structFields, removeOk(typ.FieldByName("Age")))
	structFields = append(structFields, removeOk(typ.FieldByName("Nose")))
	structFields = append(structFields, removeOk(typ.FieldByName("Fail")))

	expect := []*bunnyTag{
		{present: true, name: "test_one", bind: true},
		{present: true, name: "test_two"},
		{present: true, name: "middle_name", bind: true},
		{present: true, name: "awesome_name", structbind: true},
		nil,
		{present: false},
		nil,
	}
	for i, s := range structFields {
		tag, err := getBunnyTag(s)
		if err != nil {
			if expect[i] != nil {
				t.Errorf("Invalid tag, expected %v, got error %v", expect[i], err)
			}
		} else {
			if expect[i] == nil {
				t.Errorf("Invalid tag, expected error, got %v", tag)
			} else if tag != *expect[i] {
				t.Errorf("Invalid tag, expect %v, got %v", expect[i], tag)
			}
		}
	}
}

func TestBindChecks(t *testing.T) {
	t.Parallel()

	type useless struct {
	}

	var tests = []struct {
		BKind bindKind
		Fail  bool
		Obj   interface{}
	}{
		{BKind: kindStruct, Fail: false, Obj: &useless{}},
		{BKind: kindSliceStruct, Fail: false, Obj: &[]useless{}},
		{BKind: kindPtrSliceStruct, Fail: false, Obj: &[]*useless{}},
		{Fail: true, Obj: 5},
		{Fail: true, Obj: useless{}},
		{Fail: true, Obj: []useless{}},
	}

	for i, test := range tests {
		str, sli, bk, err := bindChecks(test.Obj)

		if err != nil {
			if !test.Fail {
				t.Errorf("%d) should not fail, got: %v", i, err)
			}
			continue
		} else if test.Fail {
			t.Errorf("%d) should fail, got: %v", i, bk)
			continue
		}

		if s := str.Kind(); s != reflect.Struct {
			t.Error("struct kind was wrong:", s)
		}
		if test.BKind != kindStruct {
			if s := sli.Kind(); s != reflect.Slice {
				t.Error("slice kind was wrong:", s)
			}
		}
	}
}

func TestBindSingular(t *testing.T) {
	t.Parallel()

	testResults := struct {
		ID   int    `bunny:"id"`
		Name string `bunny:"test"`
	}{}

	query := &Query{
		from:    []string{"fun"},
		dialect: &Dialect{LQ: '"', RQ: '"', IndexPlaceholders: true},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	ret := sqlmock.NewRows([]string{"id", "test"})
	ret.AddRow(driver.Value(int64(35)), driver.Value("pat"))
	mock.ExpectQuery(`SELECT \* FROM "fun";`).WillReturnRows(ret)

	SetContext(query, execToContext(db))
	err = query.Bind(&testResults)
	if err != nil {
		t.Error(err)
	}

	if id := testResults.ID; id != 35 {
		t.Error("wrong ID:", id)
	}
	if name := testResults.Name; name != "pat" {
		t.Error("wrong name:", name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestBind_InnerJoin(t *testing.T) {
	t.Parallel()

	testResults := []*struct {
		Happy struct {
			ID int `bunny:"identifier"`
		} `bunny:",bind"`
		Fun struct {
			ID int `bunny:"id"`
		} `bunny:",bind"`
	}{}

	query := &Query{
		from:    []string{"fun"},
		joins:   []join{{kind: JoinInner, clause: "happy as h on fun.id = h.fun_id"}},
		dialect: &Dialect{LQ: '"', RQ: '"', IndexPlaceholders: true},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	ret := sqlmock.NewRows([]string{"id"})
	ret.AddRow(driver.Value(int64(10)))
	ret.AddRow(driver.Value(int64(11)))
	mock.ExpectQuery(`SELECT "fun"\.\* FROM "fun" INNER JOIN happy as h on fun.id = h.fun_id;`).WillReturnRows(ret)

	SetContext(query, execToContext(db))
	err = query.Bind(&testResults)
	if err != nil {
		t.Error(err)
	}

	if len(testResults) != 2 {
		t.Fatal("wrong number of results:", len(testResults))
	}
	if id := testResults[0].Happy.ID; id != 0 {
		t.Error("wrong ID:", id)
	}
	if id := testResults[0].Fun.ID; id != 10 {
		t.Error("wrong ID:", id)
	}

	if id := testResults[1].Happy.ID; id != 0 {
		t.Error("wrong ID:", id)
	}
	if id := testResults[1].Fun.ID; id != 11 {
		t.Error("wrong ID:", id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestBind_InnerJoinSelect(t *testing.T) {
	t.Parallel()

	testResults := []*struct {
		Happy struct {
			ID int `bunny:"id"`
		} `bunny:"h,bind"`
		Fun struct {
			ID int `bunny:"id"`
		} `bunny:"fun,bind"`
	}{}

	query := &Query{
		dialect:    &Dialect{LQ: '"', RQ: '"', IndexPlaceholders: true},
		selectCols: []string{"fun.id", "h.id"},
		from:       []string{"fun"},
		joins:      []join{{kind: JoinInner, clause: "happy as h on fun.happy_id = h.id"}},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	ret := sqlmock.NewRows([]string{"fun.id", "h.id"})
	ret.AddRow(driver.Value(int64(10)), driver.Value(int64(11)))
	ret.AddRow(driver.Value(int64(12)), driver.Value(int64(13)))
	mock.ExpectQuery(`SELECT "fun"."id" as "fun.id", "h"."id" as "h.id" FROM "fun" INNER JOIN happy as h on fun.happy_id = h.id;`).WillReturnRows(ret)

	SetContext(query, execToContext(db))
	err = query.Bind(&testResults)
	if err != nil {
		t.Error(err)
	}

	if len(testResults) != 2 {
		t.Fatal("wrong number of results:", len(testResults))
	}
	if id := testResults[0].Happy.ID; id != 11 {
		t.Error("wrong ID:", id)
	}
	if id := testResults[0].Fun.ID; id != 10 {
		t.Error("wrong ID:", id)
	}

	if id := testResults[1].Happy.ID; id != 13 {
		t.Error("wrong ID:", id)
	}
	if id := testResults[1].Fun.ID; id != 12 {
		t.Error("wrong ID:", id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
