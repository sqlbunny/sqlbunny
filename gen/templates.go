package gen

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/KernelPay/sqlboiler/boil/queries"
	"github.com/KernelPay/sqlboiler/boil/strmangle"
	"github.com/KernelPay/sqlboiler/schema"
	"github.com/pkg/errors"
)

// templateData for sqlboiler templates
type templateData struct {
	Models  []*schema.Model
	IDTypes []*schema.IDType

	Model  *schema.Model
	Struct *schema.Struct
	Enum   *schema.Enum
	IDType *schema.IDType

	Imports []string

	// Controls what names are output
	PkgName string
	Schema  string

	// Controls which code is output (mysql vs postgres ...)
	UseLastInsertID bool

	// Turn off hook generation
	NoHooks bool

	// Tags control which
	Tags []string

	// Generate struct tags as camelCase or snake_case
	StructTagCasing string

	// StringFuncs are usable in templates with stringMap
	StringFuncs map[string]func(string) string

	// Dialect controls quoting
	Dialect queries.Dialect
	LQ      string
	RQ      string
}

func (t templateData) Quotes(s string) string {
	return fmt.Sprintf("%s%s%s", t.LQ, s, t.RQ)
}

func (t templateData) SchemaModel(model string) string {
	return strmangle.SchemaModel(t.LQ, t.RQ, model)
}

type templateList struct {
	*template.Template
}

type templateNameList []string

func (t templateNameList) Len() int {
	return len(t)
}

func (t templateNameList) Swap(k, j int) {
	t[k], t[j] = t[j], t[k]
}

func (t templateNameList) Less(k, j int) bool {
	// Make sure "struct" goes to the front
	if t[k] == "struct.tpl" {
		return true
	}

	res := strings.Compare(t[k], t[j])
	if res <= 0 {
		return true
	}

	return false
}

// Templates returns the name of all the templates defined in the template list
func (t templateList) Templates() []string {
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

	sort.Sort(templateNameList(ret))

	return ret
}

// loadTemplates loads all of the template files in the specified directory.
func loadTemplates(dir string) (*templateList, error) {
	pattern := filepath.Join(dir, "*.tpl")
	tpl, err := template.New("").Funcs(templateFunctions).ParseGlob(pattern)

	if err != nil {
		return nil, err
	}

	return &templateList{Template: tpl}, err
}

// loadTemplate loads a single template file
func loadTemplate(dir string, filename string) (*template.Template, error) {
	pattern := filepath.Join(dir, filename)
	tpl, err := template.New("").Funcs(templateFunctions).ParseFiles(pattern)

	if err != nil {
		return nil, err
	}

	return tpl.Lookup(filename), err
}

// replaceTemplate finds the template matching with name and replaces its
// contents with the contents of the template located at filename
func replaceTemplate(tpl *template.Template, name, filename string) error {
	if tpl == nil {
		return fmt.Errorf("template for %s is nil", name)
	}

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.Wrapf(err, "failed reading template file: %s", filename)
	}

	if tpl, err = tpl.New(name).Funcs(templateFunctions).Parse(string(b)); err != nil {
		return errors.Wrapf(err, "failed to parse template file: %s", filename)
	}

	return nil
}

// templateStringMappers are placed into the data to make it easy to use the
// stringMap function.
var templateStringMappers = map[string]func(string) string{
	// String ops
	"quoteWrap":       func(a string) string { return fmt.Sprintf(`"%s"`, a) },
	"replaceReserved": strmangle.ReplaceReservedWords,

	// Casing
	"titleCase":           strmangle.TitleCase,
	"titleCaseIdentifier": strmangle.TitleCaseIdentifier,
	"camelCase":           strmangle.CamelCase,
}

// templateFunctions is a map of all the functions that get passed into the
// templates. If you wish to pass a new function into your own template,
// add a function pointer here.
var templateFunctions = template.FuncMap{
	// String ops
	"quoteWrap": func(s string) string { return fmt.Sprintf(`"%s"`, s) },
	"id":        strmangle.Identifier,

	// Pluralization
	"singular": strmangle.Singular,
	"plural":   strmangle.Plural,

	// Casing
	"titleCase":           strmangle.TitleCase,
	"titleCaseIdentifier": strmangle.TitleCaseIdentifier,
	"camelCase":           strmangle.CamelCase,

	// String Slice ops
	"join":               func(sep string, slice []string) string { return strings.Join(slice, sep) },
	"joinSlices":         strmangle.JoinSlices,
	"stringMap":          strmangle.StringMap,
	"prefixStringSlice":  strmangle.PrefixStringSlice,
	"containsAny":        strmangle.ContainsAny,
	"generateTags":       strmangle.GenerateTags,
	"generateIgnoreTags": strmangle.GenerateIgnoreTags,

	// String Map ops
	"makeStringMap": strmangle.MakeStringMap,

	// Set operations
	"setInclude": strmangle.SetInclude,

	// Database related mangling
	"whereClause": strmangle.WhereClause,

	// Relationship text helpers
	"txtsFromFKey":     txtsFromFKey,
	"txtsFromOneToOne": txtsFromOneToOne,
	"txtsFromToMany":   txtsFromToMany,

	// dbdrivers ops
	"filterColumnsByDefault": schema.FilterColumnsByDefault,
	"sqlColDefinitions":      schema.SQLColDefinitions,
	"columnNames":            schema.ColumnNames,
	"getModel":               schema.GetModel,
}
