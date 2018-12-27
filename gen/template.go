package gen

import (
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/kernelpayments/sqlbunny/config"
	"github.com/kernelpayments/sqlbunny/runtime/strmangle"
	"github.com/kernelpayments/sqlbunny/schema"
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

// loadTemplates loads all of the template files in the specified directory.
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
func LoadTemplate(pkg string, path string) (*template.Template, error) {
	pkgPath, err := getPackagePath(pkg)
	if err != nil {
		return nil, err
	}

	pattern := filepath.Join(pkgPath, path)
	tpl, err := template.New("").Funcs(TemplateFunctions).ParseFiles(pattern)

	if err != nil {
		return nil, err
	}

	return tpl.Lookup(filepath.Base(path)), err
}

func MustLoadTemplates(pkg string, path string) *TemplateList {
	res, err := LoadTemplates(pkg, path)
	if err != nil {
		log.Fatalf("Error loading templates, pkg=%s path=%s: %v", pkg, path, err)
	}
	return res
}

func MustLoadTemplate(pkg string, path string) *template.Template {
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
	"titleCase":           strmangle.TitleCase,
	"titleCaseIdentifier": strmangle.TitleCaseIdentifier,
	"camelCase":           strmangle.CamelCase,
}

// templateFunctions is a map of all the functions that get passed into the
// templates. If you wish to pass a new function into your own template,
// add a function pointer here.
var TemplateFunctions = template.FuncMap{
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

	// Imports
	"import":  templateImport,
	"typeGo":  templateTypeGo,
	"typesGo": templateTypesGo,

	// Set operations
	"setInclude": strmangle.SetInclude,

	// Database related mangling
	"whereClause": strmangle.WhereClause,

	// dbdrivers ops
	"filterColumnsByDefault": schema.FilterColumnsByDefault,
	"sqlColDefinitions":      schema.SQLColDefinitions,
	"columnNames":            schema.ColumnNames,
	"getModel":               schema.GetModel,

	"quotes": func(s string) string {
		d := config.Config.Dialect
		lq := strmangle.QuoteCharacter(d.LQ)
		rq := strmangle.QuoteCharacter(d.RQ)

		return fmt.Sprintf("%s%s%s", lq, s, rq)
	},
	"schemaModel": func(model string) string {
		d := config.Config.Dialect
		lq := strmangle.QuoteCharacter(d.LQ)
		rq := strmangle.QuoteCharacter(d.RQ)
		return strmangle.SchemaModel(lq, rq, model)
	},
}
