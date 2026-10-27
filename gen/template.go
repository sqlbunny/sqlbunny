package gen

import (
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/sqlbunny/sqlbunny/runtime/strmangle"
	"github.com/sqlbunny/sqlbunny/schema"
	"golang.org/x/tools/go/packages"
)

type TemplateList struct {
	*template.Template
}

// Templates returns the name of all the templates defined in the template list
func (t TemplateList) Templates() []string {
	tplList := t.Template.Templates()

	if len(tplList) == 0 {
		return nil
	}

	ret := make([]string, 0, len(tplList))
	for _, tpl := range tplList {
		if name := tpl.Name(); strings.HasSuffix(name, ".tpl") {
			ret = append(ret, name)
		}
	}

	sort.Strings(ret)

	return ret
}

func getPackagePath(pkg string) (string, error) {
	pkgs, err := packages.Load(nil, pkg)
	if err != nil {
		return "", fmt.Errorf("Error finding package %s to load templates: %v", pkg, err)
	}

	return filepath.Dir(pkgs[0].GoFiles[0]), nil
}

// LoadTemplates loads all of the template files in the specified directory.
func LoadTemplates(pkg string, path string) (*TemplateList, error) {
	pkgPath, err := getPackagePath(pkg)
	if err != nil {
		return nil, err
	}

	pattern := filepath.Join(pkgPath, path, "*.tpl")
	tpl, err := template.New("").Funcs(TemplateFunctions).ParseGlob(pattern)

	if err != nil {
		return nil, err
	}

	return &TemplateList{Template: tpl}, err
}

// LoadTemplate loads a single template file
func LoadTemplate(pkg string, path string) (*TemplateList, error) {
	pkgPath, err := getPackagePath(pkg)
	if err != nil {
		return nil, err
	}

	pattern := filepath.Join(pkgPath, path)
	tpl, err := template.New("").Funcs(TemplateFunctions).ParseFiles(pattern)

	if err != nil {
		return nil, err
	}

	return &TemplateList{Template: tpl}, err
}

func MustLoadTemplates(pkg string, path string) *TemplateList {
	res, err := LoadTemplates(pkg, path)
	if err != nil {
		log.Fatalf("Error loading templates, pkg=%s path=%s: %v", pkg, path, err)
	}
	return res
}

func MustLoadTemplate(pkg string, path string) *TemplateList {
	res, err := LoadTemplate(pkg, path)
	if err != nil {
		log.Fatalf("Error loading template, pkg=%s path=%s: %v", pkg, path, err)
	}
	return res
}

// templateStringMappers are placed into the data to make it easy to use the
// stringMap function.
var templateStringMappers = map[string]func(string) string{
	// String ops
	"quoteWrap":       func(a string) string { return fmt.Sprintf(`"%s"`, a) },
	"replaceReserved": strmangle.ReplaceReservedWords,

	// Casing
	"titleCase": strmangle.TitleCase,
	"camelCase": strmangle.CamelCase,
}

// TemplateFunctions is a map of all the functions that get passed into the
// templates. If you wish to pass a new function into your own template,
// add a function pointer here.
var TemplateFunctions = template.FuncMap{
	// String ops
	"quoteWrap":  func(s string) string { return fmt.Sprintf(`"%s"`, s) },
	"replaceAll": strings.ReplaceAll,

	// Pluralization
	"singular": strmangle.Singular,
	"plural":   strmangle.Plural,

	// Casing
	"titleCase":     strmangle.TitleCase,
	"titleCasePath": titleCasePath,
	"camelCase":     strmangle.CamelCase,

	// String Slice ops
	"join":               func(sep string, slice []string) string { return strings.Join(slice, sep) },
	"joinSlices":         strmangle.JoinSlices,
	"stringMap":          strmangle.StringMap,
	"prefixStringSlice":  strmangle.PrefixStringSlice,
	"containsAny":        strmangle.ContainsAny,
	"generateIgnoreTags": strmangle.GenerateIgnoreTags,

	// String Map ops
	"makeStringMap": strmangle.MakeStringMap,

	// Imports
	"import":  templateImport,
	"goType":  templateGoType,
	"typesGo": templateTypesGo,

	// Set operations
	"setInclude":    strmangle.SetInclude,
	"setComplement": strmangle.SetComplement,

	// Database related mangling
	"whereClause":   WhereClause,
	"whereInClause": WhereInClause,
	"sqlType":       templateSQLType,
	"pkArrayWhere":  PKArrayWhere,
	"modelColumnSQLTypes": func(m *schema.Model) []string {
		names := modelColumns(m)
		sort.Strings(names)
		res := make([]string, len(names))
		for i, n := range names {
			res[i] = fmt.Sprintf("%q: %q", n, m.Table.Columns[n].Type)
		}
		return res
	},
	"joinOnClause":    JoinOnClause,
	"joinWhereClause": JoinWhereClause,

	"sqlNames": func(ps []schema.Path) []string {
		res := make([]string, len(ps))
		for i := range ps {
			res[i] = ps[i].SQLName()
		}
		return res
	},
	"modelColumns":      modelColumns,
	"modelPKColumns":    modelPKColumns,
	"modelNonPKColumns": modelNonPKColumns,

	"schemaModel": func(model string) string {
		return quoteIdent(model)
	},
	"hook": hook,

	"doCompare":  templateCompare,
	"mapKeyType": templateMapKeyType,
	"mapKey":     templateMapKey,
}

const nullPkg = "github.com/sqlbunny/sqlbunny/types/null"

// mapKeyConv tells how to use values of a Go type that is not comparable as
// map keys.
type mapKeyConv struct {
	// keyType is a comparable Go type to use as map key instead.
	keyType schema.GoType
	// conv returns Go code that converts expr to keyType.
	conv func(expr string) string
}

// mapKeyConvs has the non-comparable Go types that can be keys (primary keys,
// foreign keys, relationship fields), and how to use them as map keys.
// Types not in this map are used as map keys as they are.
var mapKeyConvs = map[schema.GoType]mapKeyConv{
	{Name: "[]byte"}: {
		keyType: schema.GoType{Name: "string"},
		conv: func(expr string) string {
			return "string(" + expr + ")"
		},
	},
	{Pkg: nullPkg, Name: "Bytes"}: {
		keyType: schema.GoType{Pkg: nullPkg, Name: "String"},
		conv: func(expr string) string {
			newString := templateGoType(schema.GoType{Pkg: nullPkg, Name: "NewString"})
			return newString + "(string(" + expr + ".Bytes), " + expr + ".Valid)"
		},
	},
}

// templateMapKeyType returns the Go type of map keys for values of type t.
// It is t, or a comparable type if t is not comparable (string for []byte).
func templateMapKeyType(t schema.GoType) string {
	if c, ok := mapKeyConvs[t]; ok {
		return templateGoType(c.keyType)
	}
	return templateGoType(t)
}

// templateMapKey returns Go code that converts expr, of type t, to a map key
// of type templateMapKeyType(t).
func templateMapKey(t schema.GoType, expr string) string {
	if c, ok := mapKeyConvs[t]; ok {
		return c.conv(expr)
	}
	return expr
}

func isByteSlice(t schema.GoType) bool {
	return t.Pkg == "" && t.Name == "[]byte"
}

// templateCompare returns Go code that reports whether a (of field ca) and b
// (of field cb) are equal. The fields can be nullable.
func templateCompare(a, b string, ca, cb *schema.Field) string {
	eq := func(x, y string) string {
		return x + " == " + y
	}
	isBytes := isByteSlice(ca.Type.GoType()) || isByteSlice(cb.Type.GoType())
	if isBytes {
		eq = func(x, y string) string {
			return "bytes.Equal(" + x + ", " + y + ")"
		}
	}

	if !ca.Nullable && !cb.Nullable {
		return eq(a, b)
	}

	if ca.Nullable && cb.Nullable {
		if isBytes {
			f := ca.Type.(schema.NullableType).GoTypeNullField()
			return a + ".Valid == " + b + ".Valid && " + eq(a+"."+f, b+"."+f)
		}
		return a + " == " + b
	}

	if cb.Nullable {
		a, b = b, a
		ca, cb = cb, ca
	}

	f := ca.Type.(schema.NullableType).GoTypeNullField()
	return a + ".Valid && " + eq(a+"."+f, b)
}

func modelColumns(m *schema.Model) []string {
	var res []string
	for name := range m.Table.Columns {
		res = append(res, name)
	}
	return res
}

func modelPKColumns(m *schema.Model) []string {
	var res []string
	for _, path := range m.PrimaryKey.Fields {
		res = append(res, path.SQLName())
	}
	return res
}

func modelNonPKColumns(m *schema.Model) []string {
	a := modelColumns(m)
	b := modelPKColumns(m)
	c := strmangle.SetComplement(a, b)
	return c
}

func titleCasePath(p schema.Path) string {
	var res = ""
	for i, n := range p {
		if i != 0 {
			res += "."
		}
		res += strmangle.TitleCase(n)
	}
	return res
}

// quoteIdent quotes a SQL identifier. The quotes are escaped, because the
// result is written inside a Go string literal in the generated code.
func quoteIdent(s string) string {
	return `\"` + s + `\"`
}

// WhereClause returns the where clause using start as the $ flag index
// For example, if start was 2 output would be: "colthing=$2 AND colstuff=$3"
func WhereClause(start int, cols []schema.Path) string {
	buf := strmangle.GetBuffer()
	defer strmangle.PutBuffer(buf)

	for i, c := range cols {
		if start != 0 {
			buf.WriteString(fmt.Sprintf(`%s=$%d`, quoteIdent(c.SQLName()), start+i))
		} else {
			buf.WriteString(fmt.Sprintf(`%s=?`, quoteIdent(c.SQLName())))
		}

		if i < len(cols)-1 {
			buf.WriteString(" AND ")
		}
	}

	return buf.String()
}

// JoinOnClause returns a join on clause
func JoinOnClause(table1 string, cols1 []schema.Path, table2 string, cols2 []schema.Path) string {
	buf := strmangle.GetBuffer()
	defer strmangle.PutBuffer(buf)

	for i := range cols1 {
		c1 := cols1[i].SQLName()
		c2 := cols2[i].SQLName()
		buf.WriteString(fmt.Sprintf(
			`%s.%s=%s.%s`,
			quoteIdent(table1), quoteIdent(c1),
			quoteIdent(table2), quoteIdent(c2),
		))

		if i < len(cols1)-1 {
			buf.WriteString(" AND ")
		}
	}

	return buf.String()
}

// JoinWhereClause returns a where clause explicitly specifying the table name
func JoinWhereClause(start int, table string, cols []schema.Path) string {
	buf := strmangle.GetBuffer()
	defer strmangle.PutBuffer(buf)

	for i, c := range cols {
		if start != 0 {
			buf.WriteString(fmt.Sprintf(`%s.%s=$%d`, quoteIdent(table), quoteIdent(c.SQLName()), start+i))
		} else {
			buf.WriteString(fmt.Sprintf(`%s.%s=?`, quoteIdent(table), quoteIdent(c.SQLName())))
		}

		if i < len(cols)-1 {
			buf.WriteString(" AND ")
		}
	}

	return buf.String()
}

// WhereClause returns the where clause using start as the $ flag index
// For example, if start was 2 output would be: "colthing=$2 AND colstuff=$3"
func WhereInClause(table string, cols []schema.Path) string {
	buf := strmangle.GetBuffer()
	defer strmangle.PutBuffer(buf)

	if len(cols) != 1 {
		buf.WriteString("(")
	}
	for i, c := range cols {
		buf.WriteString(fmt.Sprintf(`%s.%s`, quoteIdent(table), quoteIdent(c.SQLName())))

		if i < len(cols)-1 {
			buf.WriteString(",")
		}
	}
	if len(cols) != 1 {
		buf.WriteString(")")
	}

	return buf.String()
}

// templateSQLType returns the Postgres type name of a field, for casts like "?::bytea[]".
func templateSQLType(f *schema.Field) string {
	bt, ok := f.Type.(schema.BaseType)
	if !ok {
		panic(fmt.Sprintf("field %s has type %s with no SQL type", f.Name, f.Type.GetName()))
	}
	return bt.SQLType().Type
}

// PKArrayWhere returns a where clause that matches the given key columns
// against one array parameter per column, starting at placeholder "start":
//
//	"id" = ANY($1::bytea[])
//	("a","b") IN (SELECT * FROM unnest($1::text[], $2::int[]))
func PKArrayWhere(start int, m *schema.Model, cols []schema.Path) string {
	param := func(i int) string {
		if start == 0 {
			return "?"
		}
		return fmt.Sprintf("$%d", start+i)
	}
	if len(cols) == 1 {
		f := m.FindField(cols[0])
		return fmt.Sprintf("%s = ANY(%s::%s[])", quoteIdent(cols[0].SQLName()), param(0), templateSQLType(f))
	}
	names := make([]string, len(cols))
	params := make([]string, len(cols))
	for i, c := range cols {
		names[i] = quoteIdent(c.SQLName())
		params[i] = fmt.Sprintf("%s::%s[]", param(i), templateSQLType(m.FindField(c)))
	}
	return fmt.Sprintf("(%s) IN (SELECT * FROM unnest(%s))", strings.Join(names, ","), strings.Join(params, ", "))
}
