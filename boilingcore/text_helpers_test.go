package boilingcore

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/volatiletech/sqlboiler/bdb"
	"github.com/volatiletech/sqlboiler/bdb/drivers"
)

func TestTxtsFromOne(t *testing.T) {
	t.Parallel()

	tables, err := bdb.Tables(&drivers.MockDriver{}, "public", nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	jets := bdb.GetTable(tables, "jets")
	texts := txtsFromFKey(tables, jets, jets.FKeys[0])
	expect := TxtToOne{}

	expect.ForeignKey = jets.FKeys[0]

	expect.LocalTable.NameGo = "Jet"
	expect.LocalTable.ColumnNameGo = "PilotID"

	expect.ForeignTable.NameGo = "Pilot"
	expect.ForeignTable.NamePluralGo = "Pilots"
	expect.ForeignTable.ColumnName = "id"
	expect.ForeignTable.ColumnNameGo = "ID"

	expect.Function.Name = "pilot"
	expect.Function.ForeignName = "jet"
	expect.Function.NameGo = "Pilot"
	expect.Function.ForeignNameGo = "Jet"

	expect.Function.LocalAssignment = "PilotID.Int"
	expect.Function.ForeignAssignment = "ID"

	if !reflect.DeepEqual(expect, texts) {
		t.Errorf("Want:\n%s\nGot:\n%s\n", spew.Sdump(expect), spew.Sdump(texts))
	}

	texts = txtsFromFKey(tables, jets, jets.FKeys[1])
	expect = TxtToOne{}
	expect.ForeignKey = jets.FKeys[1]

	expect.LocalTable.NameGo = "Jet"
	expect.LocalTable.ColumnNameGo = "AirportID"

	expect.ForeignTable.NameGo = "Airport"
	expect.ForeignTable.NamePluralGo = "Airports"
	expect.ForeignTable.ColumnName = "id"
	expect.ForeignTable.ColumnNameGo = "ID"

	expect.Function.Name = "airport"
	expect.Function.ForeignName = "jets"
	expect.Function.NameGo = "Airport"
	expect.Function.ForeignNameGo = "Jets"

	expect.Function.LocalAssignment = "AirportID"
	expect.Function.ForeignAssignment = "ID"

	if !reflect.DeepEqual(expect, texts) {
		t.Errorf("Want:\n%s\nGot:\n%s\n", spew.Sdump(expect), spew.Sdump(texts))
	}

	if !reflect.DeepEqual(expect, texts) {
		t.Errorf("Want:\n%s\nGot:\n%s\n", spew.Sdump(expect), spew.Sdump(texts))
	}
}

func TestTxtsFromOneToOne(t *testing.T) {
	t.Parallel()

	tables, err := bdb.Tables(&drivers.MockDriver{}, "public", nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	pilots := bdb.GetTable(tables, "pilots")
	texts := txtsFromOneToOne(tables, pilots, pilots.ToOneRelationships[0])
	expect := TxtToOne{}

	expect.ForeignKey = bdb.ForeignKey{
		Name: "none",

		Table:    "jets",
		Column:   "pilot_id",
		Nullable: true,
		Unique:   true,

		ForeignTable:          "pilots",
		ForeignColumn:         "id",
		ForeignColumnNullable: false,
		ForeignColumnUnique:   false,
	}

	expect.LocalTable.NameGo = "Pilot"
	expect.LocalTable.ColumnNameGo = "ID"

	expect.ForeignTable.NameGo = "Jet"
	expect.ForeignTable.NamePluralGo = "Jets"
	expect.ForeignTable.ColumnName = "pilot_id"
	expect.ForeignTable.ColumnNameGo = "PilotID"

	expect.Function.Name = "jet"
	expect.Function.ForeignName = "pilot"
	expect.Function.NameGo = "Jet"
	expect.Function.ForeignNameGo = "Pilot"

	expect.Function.LocalAssignment = "ID"
	expect.Function.ForeignAssignment = "PilotID.Int"

	if !reflect.DeepEqual(expect, texts) {
		t.Errorf("Want:\n%s\nGot:\n%s\n", spew.Sdump(expect), spew.Sdump(texts))
	}
}

func TestTxtsFromMany(t *testing.T) {
	t.Parallel()

	tables, err := bdb.Tables(&drivers.MockDriver{}, "public", nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	pilots := bdb.GetTable(tables, "pilots")
	texts := txtsFromToMany(tables, pilots, pilots.ToManyRelationships[0])
	expect := TxtToMany{}
	expect.LocalTable.NameGo = "Pilot"
	expect.LocalTable.ColumnNameGo = "ID"

	expect.ForeignTable.NameGo = "License"
	expect.ForeignTable.NamePluralGo = "Licenses"
	expect.ForeignTable.NameHumanReadable = "licenses"
	expect.ForeignTable.ColumnNameGo = "PilotID"
	expect.ForeignTable.Slice = "LicenseSlice"

	expect.Function.Name = "licenses"
	expect.Function.ForeignName = "pilot"
	expect.Function.NameGo = "Licenses"
	expect.Function.ForeignNameGo = "Pilot"
	expect.Function.LocalAssignment = "ID"
	expect.Function.ForeignAssignment = "PilotID"

	if !reflect.DeepEqual(expect, texts) {
		t.Errorf("Want:\n%s\nGot:\n%s\n", spew.Sdump(expect), spew.Sdump(texts))
	}

	texts = txtsFromToMany(tables, pilots, pilots.ToManyRelationships[1])
	expect = TxtToMany{}
	expect.LocalTable.NameGo = "Pilot"
	expect.LocalTable.ColumnNameGo = "ID"

	expect.ForeignTable.NameGo = "Language"
	expect.ForeignTable.NamePluralGo = "Languages"
	expect.ForeignTable.NameHumanReadable = "languages"
	expect.ForeignTable.ColumnNameGo = "ID"
	expect.ForeignTable.Slice = "LanguageSlice"

	expect.Function.Name = "languages"
	expect.Function.ForeignName = "pilots"
	expect.Function.NameGo = "Languages"
	expect.Function.ForeignNameGo = "Pilots"
	expect.Function.LocalAssignment = "ID"
	expect.Function.ForeignAssignment = "ID"

	if !reflect.DeepEqual(expect, texts) {
		t.Errorf("Want:\n%s\nGot:\n%s\n", spew.Sdump(expect), spew.Sdump(texts))
	}
}

func TestTxtNameToOne(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Table               string
		Column              string
		Unique              bool
		ForeignTable        string
		ForeignColumn       string
		ForeignColumnUnique bool

		LocalFn   string
		ForeignFn string
	}{
		{"jets", "airport_id", false, "airports", "id", true, "airport", "jets"},
		{"jets", "airport_id", true, "airports", "id", true, "airport", "jet"},

		{"jets", "holiday_id", false, "airports", "id", true, "holiday", "holiday_jets"},
		{"jets", "holiday_id", true, "airports", "id", true, "holiday", "holiday_jet"},

		{"jets", "holiday_airport_id", false, "airports", "id", true, "holiday_airport", "holiday_airport_jets"},
		{"jets", "holiday_airport_id", true, "airports", "id", true, "holiday_airport", "holiday_airport_jet"},

		{"jets", "jet_id", false, "jets", "id", true, "jet", "jets"},
		{"jets", "jet_id", true, "jets", "id", true, "jet", "jet"},
		{"jets", "plane_id", false, "jets", "id", true, "plane", "plane_jets"},
		{"jets", "plane_id", true, "jets", "id", true, "plane", "plane_jet"},

		{"race_result_scratchings", "results_id", false, "race_results", "id", true, "result", "result_race_result_scratchings"},
	}

	for i, test := range tests {
		fk := bdb.ForeignKey{
			Table: test.Table, Column: test.Column, Unique: test.Unique,
			ForeignTable: test.ForeignTable, ForeignColumn: test.ForeignColumn, ForeignColumnUnique: test.ForeignColumnUnique,
		}

		local, foreign := txtNameToOne(fk)
		if local != test.LocalFn {
			t.Error(i, "local wrong:", local, "want:", test.LocalFn)
		}
		if foreign != test.ForeignFn {
			t.Error(i, "foreign wrong:", foreign, "want:", test.ForeignFn)
		}
	}
}

func TestTxtNameToMany(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Table  string
		Column string

		ForeignTable  string
		ForeignColumn string

		ToJoinTable       bool
		JoinLocalColumn   string
		JoinForeignColumn string

		LocalFn   string
		ForeignFn string
	}{
		{"airports", "id", "jets", "airport_id", false, "", "", "jets", "airport"},
		{"airports", "id", "jets", "holiday_airport_id", false, "", "", "holiday_airport_jets", "holiday_airport"},

		{"jets", "id", "jets", "jet_id", false, "", "", "jets", "jet"},
		{"jets", "id", "jets", "plane_id", false, "", "", "plane_jets", "plane"},

		{"pilots", "id", "languages", "id", true, "pilot_id", "language_id", "languages", "pilots"},
		{"pilots", "id", "languages", "id", true, "captain_id", "lingo_id", "lingo_languages", "captain_pilots"},

		{"pilots", "id", "pilots", "id", true, "pilot_id", "mentor_id", "mentor_pilots", "pilots"},
		{"pilots", "id", "pilots", "id", true, "mentor_id", "pilot_id", "pilots", "mentor_pilots"},
		{"pilots", "id", "pilots", "id", true, "captain_id", "mentor_id", "mentor_pilots", "captain_pilots"},

		{"race_results", "id", "race_result_scratchings", "results_id", false, "", "", "result_race_result_scratchings", "result"},
	}

	for i, test := range tests {
		fk := bdb.ToManyRelationship{
			Table: test.Table, Column: test.Column,
			ForeignTable: test.ForeignTable, ForeignColumn: test.ForeignColumn,
			ToJoinTable:     test.ToJoinTable,
			JoinLocalColumn: test.JoinLocalColumn, JoinForeignColumn: test.JoinForeignColumn,
		}

		local, foreign := txtNameToMany(fk)
		if local != test.LocalFn {
			t.Error(i, "local wrong:", local, "want:", test.LocalFn)
		}
		if foreign != test.ForeignFn {
			t.Error(i, "foreign wrong:", foreign, "want:", test.ForeignFn)
		}
	}
}

func TestTrimSuffixes(t *testing.T) {
	t.Parallel()

	for _, s := range identifierSuffixes {
		a := "hello" + s

		if z := trimSuffixes(a); z != "hello" {
			t.Errorf("got %s", z)
		}
	}
}
