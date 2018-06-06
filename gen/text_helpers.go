package gen

import (
	"fmt"
	"strings"

	"github.com/KernelPay/sqlboiler/boil/strmangle"
	"github.com/KernelPay/sqlboiler/schema"
)

// TxtToOne contains text that will be used by templates for a one-to-many or
// a one-to-one relationship.
type TxtToOne struct {
	ForeignKey *schema.ForeignKey

	LocalModel struct {
		NameGo       string
		ColumnNameGo string
	}

	ForeignModel struct {
		NameGo       string
		NamePluralGo string
		ColumnNameGo string
		ColumnName   string
	}

	Function struct {
		Name          string
		ForeignName   string
		NameGo        string
		ForeignNameGo string

		UsesBytes bool

		LocalAssignment   string
		ForeignAssignment string
	}
}

func txtsFromFKey(models []*schema.Model, model *schema.Model, fkey *schema.ForeignKey) TxtToOne {
	r := TxtToOne{}

	r.ForeignKey = fkey

	r.LocalModel.NameGo = strmangle.TitleCase(strmangle.Singular(model.Name))
	r.LocalModel.ColumnNameGo = strmangle.TitleCase(strmangle.Singular(fkey.Column))

	r.ForeignModel.NameGo = strmangle.TitleCase(strmangle.Singular(fkey.ForeignModel))
	r.ForeignModel.NamePluralGo = strmangle.TitleCase(strmangle.Plural(fkey.ForeignModel))
	r.ForeignModel.ColumnName = fkey.ForeignColumn
	r.ForeignModel.ColumnNameGo = strmangle.TitleCase(strmangle.Singular(fkey.ForeignColumn))

	r.Function.Name, r.Function.ForeignName = txtNameToOne(fkey)
	r.Function.NameGo = strmangle.TitleCase(r.Function.Name)
	r.Function.ForeignNameGo = strmangle.TitleCase(r.Function.ForeignName)

	r.Function.LocalAssignment = strmangle.TitleCaseIdentifier(fkey.Column)
	if fkey.Nullable {
		col := model.GetColumn(fkey.Column)
		r.Function.LocalAssignment += "." + col.Type.(schema.NullableType).TypeGoNullField()
	}

	foreignModel := schema.GetModel(models, fkey.ForeignModel)
	ForeignColumn := foreignModel.GetColumn(fkey.ForeignColumn)

	r.Function.ForeignAssignment = strmangle.TitleCaseIdentifier(fkey.ForeignColumn)
	if fkey.ForeignColumnNullable {
		r.Function.ForeignAssignment += "." + ForeignColumn.Type.(schema.NullableType).TypeGoNullField()
	}

	r.Function.UsesBytes = ForeignColumn.Type.TypeGo() == "[]byte"

	return r
}

func txtsFromOneToOne(models []*schema.Model, model *schema.Model, oneToOne *schema.ToOneRelationship) TxtToOne {
	fkey := &schema.ForeignKey{
		Model:    oneToOne.Model,
		Column:   oneToOne.Column,
		Nullable: oneToOne.Nullable,
		Unique:   oneToOne.Unique,

		ForeignModel:          oneToOne.ForeignModel,
		ForeignColumn:         oneToOne.ForeignColumn,
		ForeignColumnNullable: oneToOne.ForeignColumnNullable,
		ForeignColumnUnique:   oneToOne.ForeignColumnUnique,
	}

	rel := txtsFromFKey(models, model, fkey)
	col := model.GetColumn(oneToOne.Column)

	// Reverse foreign key
	rel.ForeignKey.Model, rel.ForeignKey.ForeignModel = rel.ForeignKey.ForeignModel, rel.ForeignKey.Model
	rel.ForeignKey.Column, rel.ForeignKey.ForeignColumn = rel.ForeignKey.ForeignColumn, rel.ForeignKey.Column
	rel.ForeignKey.Nullable, rel.ForeignKey.ForeignColumnNullable = rel.ForeignKey.ForeignColumnNullable, rel.ForeignKey.Nullable
	rel.ForeignKey.Unique, rel.ForeignKey.ForeignColumnUnique = rel.ForeignKey.ForeignColumnUnique, rel.ForeignKey.Unique
	rel.Function.UsesBytes = col.Type.TypeGo() == "[]byte"
	rel.Function.ForeignName, rel.Function.Name = txtNameToOne(&schema.ForeignKey{
		Model:         oneToOne.ForeignModel,
		Column:        oneToOne.ForeignColumn,
		Unique:        true,
		ForeignModel:  oneToOne.Model,
		ForeignColumn: oneToOne.Column,
	})
	rel.Function.NameGo = strmangle.TitleCase(rel.Function.Name)
	rel.Function.ForeignNameGo = strmangle.TitleCase(rel.Function.ForeignName)
	return rel
}

// TxtToMany contains text that will be used by many-to-one relationships.
type TxtToMany struct {
	LocalModel struct {
		NameGo       string
		ColumnNameGo string
	}

	ForeignModel struct {
		NameGo            string
		NamePluralGo      string
		NameHumanReadable string
		ColumnNameGo      string
		Slice             string
	}

	Function struct {
		Name          string
		ForeignName   string
		NameGo        string
		ForeignNameGo string

		UsesBytes bool

		LocalAssignment   string
		ForeignAssignment string
	}
}

// txtsFromToMany creates a struct that does a lot of the text
// transformation in advance for a given relationship.
func txtsFromToMany(models []*schema.Model, model *schema.Model, rel *schema.ToManyRelationship) TxtToMany {
	r := TxtToMany{}
	r.LocalModel.NameGo = strmangle.TitleCase(strmangle.Singular(model.Name))
	r.LocalModel.ColumnNameGo = strmangle.TitleCaseIdentifier(rel.Column)

	foreignNameSingular := strmangle.Singular(rel.ForeignModel)
	r.ForeignModel.NamePluralGo = strmangle.TitleCase(strmangle.Plural(rel.ForeignModel))
	r.ForeignModel.NameGo = strmangle.TitleCase(foreignNameSingular)
	r.ForeignModel.ColumnNameGo = strmangle.TitleCase(rel.ForeignColumn)
	r.ForeignModel.Slice = fmt.Sprintf("%sSlice", strmangle.TitleCase(foreignNameSingular))
	r.ForeignModel.NameHumanReadable = strings.Replace(rel.ForeignModel, "_", " ", -1)

	r.Function.Name, r.Function.ForeignName = txtNameToMany(rel)
	r.Function.NameGo = strmangle.TitleCase(r.Function.Name)
	r.Function.ForeignNameGo = strmangle.TitleCase(r.Function.ForeignName)

	col := model.GetColumn(rel.Column)
	r.Function.LocalAssignment = strmangle.TitleCaseIdentifier(rel.Column)
	if rel.Nullable {
		r.Function.LocalAssignment += "." + col.Type.(schema.NullableType).TypeGoNullField()
	}

	r.Function.ForeignAssignment = strmangle.TitleCaseIdentifier(rel.ForeignColumn)
	if rel.ForeignColumnNullable {
		foreignModel := schema.GetModel(models, rel.ForeignModel)
		ForeignColumn := foreignModel.GetColumn(rel.ForeignColumn)
		r.Function.ForeignAssignment += "." + ForeignColumn.Type.(schema.NullableType).TypeGoNullField()
	}

	r.Function.UsesBytes = col.Type.TypeGo() == "[]byte"

	return r
}

// txtNameToOne creates the local and foreign function names for
// one-to-many and one-to-one relationships, where local == lhs (one).
//
// = many-to-one
// users - videos : user_id
// users - videos : producer_id
//
// fk == model = user.Videos         | video.User
// fk != model = user.ProducerVideos | video.Producer
//
// = many-to-one
// industries - industries : parent_id
//
// fk == model = industry.Industries | industry.Industry
// fk != model = industry.ParentIndustries | industry.Parent
//
// = one-to-one
// users - videos : user_id
// users - videos : producer_id
//
// fk == model = user.Video         | video.User
// fk != model = user.ProducerVideo | video.Producer
//
// = one-to-one
// industries - industries : parent_id
//
// fk == model = industry.Industry | industry.Industry
// fk != model = industry.ParentIndustry | industry.Industry
func txtNameToOne(fk *schema.ForeignKey) (localFn, foreignFn string) {
	localFn = strmangle.Singular(trimSuffixes(fk.Column))
	fkeyIsModelName := localFn != strmangle.Singular(fk.ForeignModel)

	if fkeyIsModelName {
		foreignFn = localFn + "_"
	}

	plurality := strmangle.Plural
	if fk.Unique {
		plurality = strmangle.Singular
	}
	foreignFn += plurality(fk.Model)

	return localFn, foreignFn
}

// txtNameToMany creates the local and foreign function names for
// many-to-one and many-to-many relationship, where local == lhs (many)
//
// cases:
// = many-to-many
// sponsors - constests
// sponsor_id contest_id
// fk == model = sponsor.Contests | contest.Sponsors
//
// = many-to-many
// sponsors - constests
// wiggle_id jiggle_id
// fk != model = sponsor.JiggleSponsors | contest.WiggleContests
//
// = many-to-many
// industries - industries
// industry_id  mapped_industry_id
//
// fk == model = industry.Industries
// fk != model = industry.MappedIndustryIndustry
func txtNameToMany(toMany *schema.ToManyRelationship) (localFn, foreignFn string) {
	if toMany.ToJoinModel {
		localFkey := strmangle.Singular(trimSuffixes(toMany.JoinLocalField))
		foreignFkey := strmangle.Singular(trimSuffixes(toMany.JoinForeignColumn))

		if localFkey != strmangle.Singular(toMany.Model) {
			foreignFn = localFkey + "_"
		}
		foreignFn += strmangle.Plural(toMany.Model)

		if foreignFkey != strmangle.Singular(toMany.ForeignModel) {
			localFn = foreignFkey + "_"
		}
		localFn += strmangle.Plural(toMany.ForeignModel)

		return localFn, foreignFn
	}

	fkeyName := strmangle.Singular(trimSuffixes(toMany.ForeignColumn))
	if fkeyName != strmangle.Singular(toMany.Model) {
		localFn = fkeyName + "_"
	}
	localFn += strmangle.Plural(toMany.ForeignModel)
	foreignFn = strmangle.Singular(fkeyName)
	return localFn, foreignFn
}

var identifierSuffixes = []string{"_id", "_uuid", "_guid", "_oid"}

// trimSuffixes from the identifier
func trimSuffixes(str string) string {
	ln := len(str)
	for _, s := range identifierSuffixes {
		str = strings.TrimSuffix(str, s)
		if len(str) != ln {
			break
		}
	}

	return str
}
