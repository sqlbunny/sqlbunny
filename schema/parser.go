package schema

import (
	"fmt"

	p "github.com/andyleap/parser"
)

func MakeGrammar() *p.Grammar {
	Letter := p.Set("\\p{L}")
	Digit := p.Set("\\p{Nd}")
	Underscore := p.Set("_")
	Tick := p.Set("`")
	Any := p.Set("^`")
	OpenPar := p.Set("(")
	ClosePar := p.Set(")")
	Comma := p.Set(",")
	gQuotedString := p.And(p.Set("\""), p.Mult(0, 0, p.Set("^\"")), p.Set("\""))

	WS := p.Ignore(p.Mult(0, 0, p.Set("\t\n\f\r ")))
	RWS := p.Ignore(p.Mult(1, 0, p.Set("\t\n\f\r ")))
	RWSN := p.Ignore(p.Mult(1, 0, p.Set("\t\f ")))
	NL := p.Ignore(p.And(p.Mult(0, 0, p.Set("\t\f\r ")), p.Mult(1, 0, p.And(p.Lit("\n"), WS))))

	gIdentifier := p.And(Letter, p.Mult(0, 0, p.Or(Letter, Digit, Underscore)))
	gIdentifier.Node(func(m p.Match) (p.Match, error) {
		return p.String(m), nil
	})

	gTypeGo := p.Mult(1, 0, p.Or(Letter, Digit, p.Set(".\\[\\]")))
	gTypeGo.Node(func(m p.Match) (p.Match, error) {
		return p.String(m), nil
	})

	gType := p.And(gIdentifier)
	gType.Node(func(m p.Match) (p.Match, error) {
		return p.String(m), nil
	})

	gStructTag := p.And(Tick, p.Mult(0, 0, p.Or(Any)), Tick)
	gStructTag.Node(func(m p.Match) (p.Match, error) {
		return p.String(m), nil
	})

	gField := p.And(
		p.Tag("Name", gIdentifier),
		p.Require(RWSN,
			p.Tag("Type", gType),
		),
		p.Optional(p.And(RWSN, p.Tag("Null", p.Lit("null")))),
		p.Optional(p.And(RWSN, p.Tag("PrimaryKey", p.Lit("primary_key")))),
		p.Optional(p.And(RWSN, p.Tag("Index", p.Lit("index")))),
		p.Optional(p.And(RWSN, p.Tag("Unique", p.Lit("unique")))),
		p.Optional(p.And(RWSN, p.Lit("foreign_key"), OpenPar, p.Tag("ForeignKey", gIdentifier), ClosePar)),
		p.Optional(p.And(RWSN, p.And(Tick, p.Tag("Tag", p.Mult(0, 0, p.Or(Any))), Tick))),
		p.Require(NL),
	)
	gField.Node(func(m p.Match) (p.Match, error) {
		tags, err := TagsFromString(p.String(p.GetTag(m, "Tag")))
		if err != nil {
			return nil, err
		}
		f := &Field{
			Name:       p.GetTag(m, "Name").(string),
			Tags:       tags,
			typeName:   p.GetTag(m, "Type").(string),
			Nullable:   p.GetTag(m, "Null") != nil,
			index:      p.GetTag(m, "Index") != nil,
			unique:     p.GetTag(m, "Unique") != nil,
			primaryKey: p.GetTag(m, "PrimaryKey") != nil,
			foreignKey: p.String(p.GetTag(m, "ForeignKey")),
		}
		return p.TagMatch("Field", f), nil
	})

	gBaseType := p.And(
		p.Lit("type"),
		p.Require(
			RWS,
			p.Tag("Name", gIdentifier), WS,
			p.Lit("{"), WS,
			p.Mult(0, 0, p.And(p.Lit("import"), p.Require(RWS, p.Tag("Import", p.And(p.Optional(p.And(gIdentifier, RWS)), gQuotedString)), NL))),
			p.Lit("not_null"), RWS, p.Tag("Go", gTypeGo), NL,
			p.Optional(p.And(p.Lit("null"), p.Require(RWS, p.Tag("GoNull", gTypeGo), NL))),
			p.Lit("postgres"), RWS, p.Tag("Postgres", gTypeGo), NL,
			p.Lit("}"), WS,
		),
	)
	gBaseType.Node(func(m p.Match) (p.Match, error) {
		var goImports []string
		for _, v := range p.GetTags(m, "Import") {
			goImports = append(goImports, p.String(v))
		}

		var s BaseType
		if p.GetTag(m, "GoNull") != nil {
			s = &BaseTypeNullable{
				Name:      p.GetTag(m, "Name").(string),
				Postgres:  p.GetTag(m, "Postgres").(string),
				Go:        p.GetTag(m, "Go").(string),
				GoNull:    p.GetTag(m, "GoNull").(string),
				GoImports: goImports,
			}
		} else {
			s = &BaseTypeNotNullable{
				Name:      p.GetTag(m, "Name").(string),
				Postgres:  p.GetTag(m, "Postgres").(string),
				Go:        p.GetTag(m, "Go").(string),
				GoImports: goImports,
			}
		}
		return p.TagMatch("BaseType", s), nil
	})

	gIDType := p.And(
		p.Lit("id_type"),
		p.Require(
			RWS,
			p.Tag("Name", gIdentifier),
			RWS,
			p.Tag("Prefix", gIdentifier),
			NL,
		),
	)
	gIDType.Node(func(m p.Match) (p.Match, error) {
		s := &IDType{
			Name:   p.GetTag(m, "Name").(string),
			Prefix: p.GetTag(m, "Prefix").(string),
		}
		return p.TagMatch("IDType", s), nil
	})

	gEnumChoice := p.And(
		p.Tag("Name", gIdentifier),
		p.Require(NL),
	)
	gEnumChoice.Node(func(m p.Match) (p.Match, error) {
		f := p.GetTag(m, "Name").(string)
		return p.TagMatch("EnumChoice", f), nil
	})

	gEnum := p.And(
		p.Lit("enum"),
		p.Require(
			RWS,
			p.Tag("Name", gIdentifier), RWS,
			p.Tag("Type", gType), WS,
			p.Lit("{"), WS,
			p.Mult(0, 0, gEnumChoice),
			p.Lit("}"), WS,
		),
	)
	gEnum.Node(func(m p.Match) (p.Match, error) {
		s := &Enum{
			Name:     p.GetTag(m, "Name").(string),
			typeName: p.GetTag(m, "Type").(string),
		}
		for _, v := range p.GetTags(m, "EnumChoice") {
			s.Choices = append(s.Choices, v.(string))
		}
		return p.TagMatch("Enum", s), nil
	})

	gStruct := p.And(
		p.Lit("struct"),
		p.Require(
			RWS,
			p.Tag("Name", gIdentifier), WS,
			p.Lit("{"), WS,
			p.Mult(0, 0, gField),
			p.Lit("}"), WS,
		),
	)
	gStruct.Node(func(m p.Match) (p.Match, error) {
		s := &Struct{
			Name: p.GetTag(m, "Name").(string),
		}
		for _, v := range p.GetTags(m, "Field") {
			s.Fields = append(s.Fields, v.(*Field))
		}
		return p.TagMatch("Struct", s), nil
	})

	gIdentifierList := p.And(
		p.Tag("Ident", gIdentifier),
		p.Mult(0, 0, p.And(WS, Comma, WS, p.Tag("Ident", gIdentifier))),
	)
	gIdentifierList.Node(func(m p.Match) (p.Match, error) {
		var res []string
		for _, v := range p.GetTags(m, "Ident") {
			res = append(res, p.String(v))
		}
		return p.TagMatch("IdentList", res), nil
	})

	gPrimaryKey := p.And(
		p.Lit("primary_key"),
		p.Require(WS, OpenPar, gIdentifierList, WS, ClosePar, WS),
	)
	gPrimaryKey.Node(func(m p.Match) (p.Match, error) {
		pk := &PrimaryKey{
			Columns: p.GetTag(m, "IdentList").([]string),
		}
		return p.TagMatch("PrimaryKey", pk), nil
	})

	gIndex := p.And(
		p.Lit("index"),
		p.Require(WS, OpenPar, gIdentifierList, WS, ClosePar, WS),
	)
	gIndex.Node(func(m p.Match) (p.Match, error) {
		pk := &Index{
			Columns: p.GetTag(m, "IdentList").([]string),
		}
		return p.TagMatch("Index", pk), nil
	})

	gUnique := p.And(
		p.Lit("unique"),
		p.Require(WS, OpenPar, gIdentifierList, WS, ClosePar, WS),
	)
	gUnique.Node(func(m p.Match) (p.Match, error) {
		k := &Unique{
			Columns: p.GetTag(m, "IdentList").([]string),
		}
		return p.TagMatch("Unique", k), nil
	})

	gModel := p.And(
		p.Lit("model"),
		p.Require(
			RWS,
			p.Tag("Name", gIdentifier), WS,
			p.Lit("{"), WS,
			p.Mult(0, 0, p.Or(gPrimaryKey, gIndex, gUnique, gField)),
			p.Lit("}"), WS,
		),
	)
	gModel.Node(func(m p.Match) (p.Match, error) {
		s := &Model{
			Name: p.GetTag(m, "Name").(string),
		}
		for _, v := range p.GetTags(m, "Field") {
			s.Fields = append(s.Fields, v.(*Field))
		}
		for _, v := range p.GetTags(m, "PrimaryKey") {
			if s.PrimaryKey != nil {
				return nil, fmt.Errorf("Model %s has multiple primary_keys", s.Name)
			}
			s.PrimaryKey = v.(*PrimaryKey)
		}
		for _, v := range p.GetTags(m, "Index") {
			s.Indexes = append(s.Indexes, v.(*Index))
		}
		for _, v := range p.GetTags(m, "Unique") {
			s.Uniques = append(s.Uniques, v.(*Unique))
		}
		return p.TagMatch("Model", s), nil
	})

	gSchema := p.And(WS, p.Mult(0, 0, p.Or(gBaseType, gIDType, gEnum, gStruct, gModel)), WS)
	gSchema.Node(func(m p.Match) (p.Match, error) {
		s := &Schema{}
		for _, v := range p.GetTags(m, "IDType") {
			s.IDTypes = append(s.IDTypes, v.(*IDType))
		}
		for _, v := range p.GetTags(m, "BaseType") {
			s.BaseTypes = append(s.BaseTypes, v.(BaseType))
		}
		for _, v := range p.GetTags(m, "Enum") {
			s.Enums = append(s.Enums, v.(*Enum))
		}
		for _, v := range p.GetTags(m, "Struct") {
			s.Structs = append(s.Structs, v.(*Struct))
		}
		for _, v := range p.GetTags(m, "Model") {
			s.Models = append(s.Models, v.(*Model))
		}
		return s, nil
	})
	return gSchema
}
