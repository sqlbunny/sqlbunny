// Package boilingcore has types and methods useful for generating code that
// acts as a fully dynamic ORM might.
package gen

import (
	"go/build"
	"os"
	"path/filepath"
	"regexp"
	"text/template"

	"github.com/KernelPay/sqlboiler/boil/queries"
	"github.com/KernelPay/sqlboiler/boil/strmangle"
	"github.com/KernelPay/sqlboiler/schema"
	"github.com/pkg/errors"
)

const (
	templatesModelDirectory     = "templates/model"
	templatesStructDirectory    = "templates/struct"
	templatesEnumDirectory      = "templates/enum"
	templatesIDDirectory        = "templates/id"
	templatesSingletonDirectory = "templates/singleton"

	templatesTestDirectory          = "templates_test"
	templatesSingletonTestDirectory = "templates_test/singleton"

	templatesTestMainDirectory = "templates_test/main_test"
)

// State holds the global data needed by most pieces to run
type State struct {
	Config *Config

	Schema  *schema.Schema
	Dialect queries.Dialect

	ModelTemplates     *templateList
	ModelTestTemplates *templateList
	StructTemplates    *templateList
	EnumTemplates      *templateList
	IDTemplates        *templateList

	SingletonTemplates     *templateList
	SingletonTestTemplates *templateList

	TestMainTemplate *template.Template
}

// New creates a new state based off of the config
func New(schema *schema.Schema, config *Config) (*State, error) {
	s := &State{
		Schema: schema,
		Config: config,
		Dialect: queries.Dialect{
			LQ:                '"',
			RQ:                '"',
			IndexPlaceholders: true,
			UseTopClause:      false,
		},
	}

	err := s.initOutFolder()
	if err != nil {
		return nil, errors.Wrap(err, "unable to initialize the output folder")
	}

	err = s.initTemplates()
	if err != nil {
		return nil, errors.Wrap(err, "unable to initialize templates")
	}

	err = s.initTags(config.Tags)
	if err != nil {
		return nil, errors.Wrap(err, "unable to initialize struct tags")
	}

	return s, nil
}

// Run executes the sqlboiler templates and outputs them to files based on the
// state given.
func (s *State) Run(includeTests bool) error {
	singletonData := &templateData{
		Models:          s.Schema.Models,
		IDTypes:         s.Schema.IDTypes,
		UseLastInsertID: true,
		PkgName:         s.Config.PkgName,
		NoHooks:         s.Config.NoHooks,
		Dialect:         s.Dialect,
		LQ:              strmangle.QuoteCharacter(s.Dialect.LQ),
		RQ:              strmangle.QuoteCharacter(s.Dialect.RQ),

		StringFuncs: templateStringMappers,
	}

	if err := s.executeSingletonTemplates(singletonData, s.SingletonTemplates); err != nil {
		return errors.Wrap(err, "singleton template output")
	}

	if !s.Config.NoTests && includeTests {
		if err := s.generateTestMainOutput(s, singletonData); err != nil {
			return errors.Wrap(err, "unable to generate TestMain output")
		}

		if err := s.executeSingletonTemplates(singletonData, s.SingletonTestTemplates); err != nil {
			return errors.Wrap(err, "unable to generate singleton test template output")
		}
	}

	for _, e := range s.Schema.IDTypes {
		data := &templateData{
			Models:          s.Schema.Models,
			IDType:          e,
			UseLastInsertID: true,
			PkgName:         s.Config.PkgName,
			NoHooks:         s.Config.NoHooks,
			Dialect:         s.Dialect,
			LQ:              strmangle.QuoteCharacter(s.Dialect.LQ),
			RQ:              strmangle.QuoteCharacter(s.Dialect.RQ),

			StringFuncs: templateStringMappers,
		}

		// Generate the regular templates
		if err := s.executeTemplates(data, s.IDTemplates, e.Name+".go"); err != nil {
			return errors.Wrap(err, "unable to generate output")
		}
	}

	for _, e := range s.Schema.Enums {
		data := &templateData{
			Models:          s.Schema.Models,
			Enum:            e,
			UseLastInsertID: true,
			PkgName:         s.Config.PkgName,
			NoHooks:         s.Config.NoHooks,
			Dialect:         s.Dialect,
			LQ:              strmangle.QuoteCharacter(s.Dialect.LQ),
			RQ:              strmangle.QuoteCharacter(s.Dialect.RQ),

			StringFuncs: templateStringMappers,
		}

		// Generate the regular templates
		if err := s.executeTemplates(data, s.EnumTemplates, e.Name+".go"); err != nil {
			return errors.Wrap(err, "unable to generate output")
		}
	}

	for _, st := range s.Schema.Structs {
		data := &templateData{
			Models:          s.Schema.Models,
			Struct:          st,
			UseLastInsertID: true,
			PkgName:         s.Config.PkgName,
			NoHooks:         s.Config.NoHooks,
			Dialect:         s.Dialect,
			LQ:              strmangle.QuoteCharacter(s.Dialect.LQ),
			RQ:              strmangle.QuoteCharacter(s.Dialect.RQ),

			StringFuncs: templateStringMappers,
		}

		// Generate the regular templates
		if err := s.executeTemplates(data, s.StructTemplates, st.Name+".go"); err != nil {
			return errors.Wrap(err, "unable to generate output")
		}
	}

	for _, model := range s.Schema.Models {
		if model.IsJoinModel {
			continue
		}

		data := &templateData{
			Models:          s.Schema.Models,
			Model:           model,
			UseLastInsertID: true,
			PkgName:         s.Config.PkgName,
			NoHooks:         s.Config.NoHooks,
			Dialect:         s.Dialect,
			LQ:              strmangle.QuoteCharacter(s.Dialect.LQ),
			RQ:              strmangle.QuoteCharacter(s.Dialect.RQ),

			StringFuncs: templateStringMappers,
		}

		// Generate the regular templates
		if err := s.executeTemplates(data, s.ModelTemplates, model.Name+".go"); err != nil {
			return errors.Wrap(err, "unable to generate output")
		}

		// Generate the test templates
		if !s.Config.NoTests && includeTests {
			if err := s.executeTemplates(data, s.ModelTestTemplates, model.Name+"_test.go"); err != nil {
				return errors.Wrap(err, "unable to generate test output")
			}
		}
	}

	return nil
}

// initTemplates loads all template folders into the state object.
func (s *State) initTemplates() error {
	var err error

	basePath, err := getBasePath(s.Config.BaseDir)
	if err != nil {
		return err
	}

	s.ModelTemplates, err = loadTemplates(filepath.Join(basePath, templatesModelDirectory))
	if err != nil {
		return err
	}
	s.StructTemplates, err = loadTemplates(filepath.Join(basePath, templatesStructDirectory))
	if err != nil {
		return err
	}
	s.EnumTemplates, err = loadTemplates(filepath.Join(basePath, templatesEnumDirectory))
	if err != nil {
		return err
	}
	s.IDTemplates, err = loadTemplates(filepath.Join(basePath, templatesIDDirectory))
	if err != nil {
		return err
	}

	s.SingletonTemplates, err = loadTemplates(filepath.Join(basePath, templatesSingletonDirectory))
	if err != nil {
		return err
	}

	if !s.Config.NoTests {
		s.ModelTestTemplates, err = loadTemplates(filepath.Join(basePath, templatesTestDirectory))
		if err != nil {
			return err
		}

		s.SingletonTestTemplates, err = loadTemplates(filepath.Join(basePath, templatesSingletonTestDirectory))
		if err != nil {
			return err
		}

		s.TestMainTemplate, err = loadTemplate(filepath.Join(basePath, templatesTestMainDirectory), "main.tpl")
		if err != nil {
			return err
		}
	}

	return nil
}

var basePackage = "github.com/KernelPay/sqlboiler/gen"

func getBasePath(baseDirConfig string) (string, error) {
	if len(baseDirConfig) > 0 {
		return baseDirConfig, nil
	}

	p, _ := build.Default.Import(basePackage, "", build.FindOnly)
	if p != nil && len(p.Dir) > 0 {
		return p.Dir, nil
	}

	return os.Getwd()
}

// Tags must be in a format like: json, xml, etc.
var rgxValidTag = regexp.MustCompile(`[a-zA-Z_\.]+`)

// initTags removes duplicate tags and validates the format
// of all user tags are simple strings without quotes: [a-zA-Z_\.]+
func (s *State) initTags(tags []string) error {
	s.Config.Tags = removeDuplicates(s.Config.Tags)
	for _, v := range s.Config.Tags {
		if !rgxValidTag.MatchString(v) {
			return errors.New("Invalid tag format %q supplied, only specify name, eg: xml")
		}
	}

	return nil
}

// initOutFolder creates the folder that will hold the generated output.
func (s *State) initOutFolder() error {
	if s.Config.Wipe {
		if err := os.RemoveAll(s.Config.OutFolder); err != nil {
			return err
		}
	}

	return os.MkdirAll(s.Config.OutFolder, os.ModePerm)
}
