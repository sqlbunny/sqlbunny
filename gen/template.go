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
	"quoteWrap": func(s string) string { return fmt.Sprintf(`"%s"`, s) },

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
	"whereClause":     WhereClause,
	"whereInClause":   WhereInClause,
	"joinOnClause":    JoinOnClause,
	"joinWhereClause": JoinWhereClause,

	"sqlNames": func(ps []schema.Path) []string {
		res := make([]string, len(ps))
		for i := range ps {
			res[i] = ps[i].SQLName()
		}
		return res
	},
	"modelColumns":   modelColumns,
	"modelPKColumns": modelPKColumns,
	"modelNonPKColumns": func(m *schema.Model) []string {
		var res []string
		for _, path := range m.PrimaryKey.Fields {
			res = append(res, path.SQLName())
		}
		return res
	},

	"quotes": func(s string) string {
		d := Config.Dialect
		lq := strmangle.QuoteCharacter(d.LQ)
		rq := strmangle.QuoteCharacter(d.RQ)

		return fmt.Sprintf("%s%s%s", lq, s, rq)
	},
	"schemaModel": func(model string) string {
		d := Config.Dialect
		lq := strmangle.QuoteCharacter(d.LQ)
		rq := strmangle.QuoteCharacter(d.RQ)
		return strmangle.SchemaModel(lq, rq, model)
	},
	"hook": hook,

	"doCompare": func(a, b string, ca, cb *schema.Field) string {
		if ca.Type.GoType().Name == "[]byte" && cb.Type.GoType().Name == "[]byte" {
			return "0 == bytes.Compare(" + a + ", " + b + ")"
		}

		if ca.Nullable == cb.Nullable {
			return a + " == " + b
		}

		if cb.Nullable {
			a, b = b, a
			ca, cb = cb, ca
		}

		f := ca.Type.(schema.NullableType).GoTypeNullField()
		return a + ".Valid && " + a + "." + f + " == " + b
	},
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

// WhereClause returns the where clause using start as the $ flag index
// For example, if start was 2 output would be: "colthing=$2 AND colstuff=$3"
func WhereClause(lq, rq string, start int, cols []schema.Path) string {
	buf := strmangle.GetBuffer()
	defer strmangle.PutBuffer(buf)

	for i, c := range cols {
		if start != 0 {
			buf.WriteString(fmt.Sprintf(`%s%s%s=$%d`, lq, c.SQLName(), rq, start+i))
		} else {
			buf.WriteString(fmt.Sprintf(`%s%s%s=?`, lq, c.SQLName(), rq))
		}

		if i < len(cols)-1 {
			buf.WriteString(" AND ")
		}
	}

	return buf.String()
}

// WhereClauseRepeated returns the where clause repeated with OR clause using start as the $ flag index
// For example, if start was 2 output would be: "(colthing=$2 AND colstuff=$3) OR (colthing=$4 AND colstuff=$5)"
func WhereClauseRepeated(lq, rq string, start int, cols []schema.Path, count int) string {
	var startIndex int
	buf := strmangle.GetBuffer()
	defer strmangle.PutBuffer(buf)
	buf.WriteByte('(')
	for i := 0; i < count; i++ {
		if i != 0 {
			buf.WriteString(") OR (")
		}

		startIndex = 0
		if start > 0 {
			startIndex = start + i*len(cols)
		}

		buf.WriteString(WhereClause(lq, rq, startIndex, cols))
	}
	buf.WriteByte(')')

	return buf.String()
}

// JoinOnClause returns a join on clause
func JoinOnClause(lq, rq string, table1 string, cols1 []schema.Path, table2 string, cols2 []schema.Path) string {
	buf := strmangle.GetBuffer()
	defer strmangle.PutBuffer(buf)

	for i := range cols1 {
		c1 := cols1[i].SQLName()
		c2 := cols2[i].SQLName()
		buf.WriteString(fmt.Sprintf(
			`%s%s%s.%s%s%s=%s%s%s.%s%s%s`,
			lq, table1, rq, lq, c1, rq,
			lq, table2, rq, lq, c2, rq,
		))

		if i < len(cols1)-1 {
			buf.WriteString(" AND ")
		}
	}

	return buf.String()
}

// JoinWhereClause returns a where clause explicitly specifying the table name
func JoinWhereClause(lq, rq string, start int, table string, cols []schema.Path) string {
	buf := strmangle.GetBuffer()
	defer strmangle.PutBuffer(buf)

	for i, c := range cols {
		if start != 0 {
			buf.WriteString(fmt.Sprintf(`%s%s%s.%s%s%s=$%d`, lq, table, rq, lq, c.SQLName(), rq, start+i))
		} else {
			buf.WriteString(fmt.Sprintf(`%s%s%s.%s%s%s=?`, lq, table, rq, lq, c.SQLName(), rq))
		}

		if i < len(cols)-1 {
			buf.WriteString(" AND ")
		}
	}

	return buf.String()
}

// WhereClause returns the where clause using start as the $ flag index
// For example, if start was 2 output would be: "colthing=$2 AND colstuff=$3"
func WhereInClause(lq, rq string, table string, cols []schema.Path) string {
	buf := strmangle.GetBuffer()
	defer strmangle.PutBuffer(buf)

	if len(cols) != 1 {
		buf.WriteString("(")
	}
	for i, c := range cols {
		buf.WriteString(fmt.Sprintf(`%s%s%s.%s%s%s`, lq, table, rq, lq, c.SQLName(), rq))

		if i < len(cols)-1 {
			buf.WriteString(",")
		}
	}
	if len(cols) != 1 {
		buf.WriteString(")")
	}

	return buf.String()
}
