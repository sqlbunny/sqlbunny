package core

import (
	"log"
	"os"

	"github.com/sqlbunny/sqlbunny/gen"
	"github.com/sqlbunny/sqlbunny/schema"
)

const (
	templatesPackage = "github.com/sqlbunny/sqlbunny/gen/core"

	templatesModelDirectory     = "templates/model"
	templatesStructDirectory    = "templates/struct"
	templatesEnumDirectory      = "templates/enum"
	templatesSingletonDirectory = "templates/singleton"
)

type Plugin struct {
	ModelTemplates     *gen.TemplateList
	StructTemplates    *gen.TemplateList
	EnumTemplates      *gen.TemplateList
	SingletonTemplates *gen.TemplateList
}

var _ gen.Plugin = &Plugin{}

func (*Plugin) IsConfigItem() {}

func (p *Plugin) BunnyPlugin() {
	schema, err := buildSchema(gen.Config.Items)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	gen.Config.Schema = schema

	p.ModelTemplates = gen.MustLoadTemplates(templatesPackage, templatesModelDirectory)
	p.StructTemplates = gen.MustLoadTemplates(templatesPackage, templatesStructDirectory)
	p.EnumTemplates = gen.MustLoadTemplates(templatesPackage, templatesEnumDirectory)
	p.SingletonTemplates = gen.MustLoadTemplates(templatesPackage, templatesSingletonDirectory)

	gen.OnGen(p.gen)
}

func (p *Plugin) gen() {
	if err := os.MkdirAll(gen.Config.ModelsOutputPath(), os.ModePerm); err != nil {
		log.Fatalf("Error creating output directory %s: %v", gen.Config.ModelsOutputPath(), err)
	}

	var models []*schema.Model
	for _, m := range gen.Config.Schema.Models {
		models = append(models, m)
	}

	data := gen.BaseTemplateData()
	data["Models"] = models
	p.SingletonTemplates.ExecuteSingleton(data)

	for _, t := range gen.Config.Schema.Types {
		switch t := t.(type) {
		case *schema.Enum:
			data := gen.BaseTemplateData()
			data["Enum"] = t
			p.EnumTemplates.Execute(data, t.Name+".go")
		case *schema.Struct:
			data := gen.BaseTemplateData()
			data["Struct"] = t
			p.StructTemplates.Execute(data, t.Name+".go")
		}
	}

	for _, model := range gen.Config.Schema.Models {
		if model.IsJoinModel {
			continue
		}

		data := gen.BaseTemplateData()
		data["Model"] = model
		data["Models"] = models
		p.ModelTemplates.Execute(data, model.Name+".go")
	}
}
