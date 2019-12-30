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
	"titleCase":           strmangle.TitleCase,
	"titleCaseIdentifier": strmangle.TitleCaseIdentifier,
	"camelCase":           strmangle.CamelCase,
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
	"titleCase":           strmangle.TitleCase,
	"titleCaseIdentifier": strmangle.TitleCaseIdentifier,
	"camelCase":           strmangle.CamelCase,

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
	"whereClause":     strmangle.WhereClause,
	"whereInClause":   strmangle.WhereInClause,
	"joinOnClause":    strmangle.JoinOnClause,
	"joinWhereClause": strmangle.JoinWhereClause,

	// dbdrivers ops
	"sqlColDefinitions": schema.SQLColDefinitions,
	"columnNames":       schema.ColumnNames,
	"getModel":          schema.GetModel,

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

	"doCompare": func(a, b string, ca, cb *schema.Column) string {
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
