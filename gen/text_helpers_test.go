package gen

import (
	"reflect"
	"testing"

	"github.com/kernelpayments/sqlbunny/schema"
	"github.com/davecgh/go-spew/spew"
)

func TestTxtsFromOne(t *testing.T) {
	t.Parallel()

	models, err := schema.Models(&drivers.MockDriver{}, "public", nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	jets := schema.GetModel(models, "jets")
	texts := txtsFromFKey(models, jets, jets.ForeignKeys[0])
	expect := TxtToOne{}

	expect.ForeignKey = jets.ForeignKeys[0]

	expect.LocalModel.NameGo = "Jet"
	expect.LocalModel.ColumnNameGo = "PilotID"

	expect.ForeignModel.NameGo = "Pilot"
	expect.ForeignModel.NamePluralGo = "Pilots"
	expect.ForeignModel.ColumnName = "id"
	expect.ForeignModel.ColumnNameGo = "ID"

	expect.Function.Name = "pilot"
	expect.Function.ForeignName = "jet"
	expect.Function.NameGo = "Pilot"
	expect.Function.ForeignNameGo = "Jet"

	expect.Function.LocalAssignment = "PilotID.Int"
	expect.Function.ForeignAssignment = "ID"

	if !reflect.DeepEqual(expect, texts) {
		t.Errorf("Want:\n%s\nGot:\n%s\n", spew.Sdump(expect), spew.Sdump(texts))
	}

	texts = txtsFromFKey(models, jets, jets.ForeignKeys[1])
	expect = TxtToOne{}
	expect.ForeignKey = jets.ForeignKeys[1]

	expect.LocalModel.NameGo = "Jet"
	expect.LocalModel.ColumnNameGo = "AirportID"

	expect.ForeignModel.NameGo = "Airport"
	expect.ForeignModel.NamePluralGo = "Airports"
	expect.ForeignModel.ColumnName = "id"
	expect.ForeignModel.ColumnNameGo = "ID"

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

	models, err := schema.Models(&drivers.MockDriver{}, "public", nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	pilots := schema.GetModel(models, "pilots")
	texts := txtsFromOneToOne(models, pilots, pilots.ToOneRelationships[0])
	expect := TxtToOne{}

	expect.ForeignKey = schema.ForeignKey{
		Name: "none",

		Model:    "jets",
		Field:    "pilot_id",
		Nullable: true,
		Unique:   true,

		ForeignModel:         "pilots",
		ForeignColumn:         "id",
		ForeignColumnNullable: false,
		ForeignColumnUnique:   false,
	}

	expect.LocalModel.NameGo = "Pilot"
	expect.LocalModel.ColumnNameGo = "ID"

	expect.ForeignModel.NameGo = "Jet"
	expect.ForeignModel.NamePluralGo = "Jets"
	expect.ForeignModel.ColumnName = "pilot_id"
	expect.ForeignModel.ColumnNameGo = "PilotID"

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

	models, err := schema.Models(&drivers.MockDriver{}, "public", nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	pilots := schema.GetModel(models, "pilots")
	texts := txtsFromToMany(models, pilots, pilots.ToManyRelationships[0])
	expect := TxtToMany{}
	expect.LocalModel.NameGo = "Pilot"
	expect.LocalModel.ColumnNameGo = "ID"

	expect.ForeignModel.NameGo = "License"
	expect.ForeignModel.NamePluralGo = "Licenses"
	expect.ForeignModel.NameHumanReadable = "licenses"
	expect.ForeignModel.ColumnNameGo = "PilotID"
	expect.ForeignModel.Slice = "LicenseSlice"

	expect.Function.Name = "licenses"
	expect.Function.ForeignName = "pilot"
	expect.Function.NameGo = "Licenses"
	expect.Function.ForeignNameGo = "Pilot"
	expect.Function.LocalAssignment = "ID"
	expect.Function.ForeignAssignment = "PilotID"

	if !reflect.DeepEqual(expect, texts) {
		t.Errorf("Want:\n%s\nGot:\n%s\n", spew.Sdump(expect), spew.Sdump(texts))
	}

	texts = txtsFromToMany(models, pilots, pilots.ToManyRelationships[1])
	expect = TxtToMany{}
	expect.LocalModel.NameGo = "Pilot"
	expect.LocalModel.ColumnNameGo = "ID"

	expect.ForeignModel.NameGo = "Language"
	expect.ForeignModel.NamePluralGo = "Languages"
	expect.ForeignModel.NameHumanReadable = "languages"
	expect.ForeignModel.ColumnNameGo = "ID"
	expect.ForeignModel.Slice = "LanguageSlice"

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
		Model              string
		Field              string
		Unique             bool
		ForeignModel       string
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
		fk := schema.ForeignKey{
			Model: test.Model, Field: test.Field, Unique: test.Unique,
			ForeignModel: test.ForeignModel, ForeignColumn: test.ForeignColumn, ForeignColumnUnique: test.ForeignColumnUnique,
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
		Model string
		Field string

		ForeignModel string
		ForeignColumn string

		ToJoinModel      bool
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
		fk := schema.ToManyRelationship{
			Model: test.Model, Field: test.Field,
			ForeignModel: test.ForeignModel, ForeignColumn: test.ForeignColumn,
			ToJoinModel:    test.ToJoinModel,
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
