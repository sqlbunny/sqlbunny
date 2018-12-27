package core

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/kernelpayments/sqlbunny/config"
	"github.com/kernelpayments/sqlbunny/gen"
	"github.com/kernelpayments/sqlbunny/schema"
)

const (
	templatesPackage = "github.com/kernelpayments/sqlbunny/gen/core"

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

func (*Plugin) IsConfigItem() {}

func (p *Plugin) InitPlugin() {
	config.Config.RootCmd.AddCommand(&cobra.Command{
		Use: "gen",
		Run: p.cmdGen,
	})

	p.ModelTemplates = gen.MustLoadTemplates(templatesPackage, templatesModelDirectory)
	p.StructTemplates = gen.MustLoadTemplates(templatesPackage, templatesStructDirectory)
	p.EnumTemplates = gen.MustLoadTemplates(templatesPackage, templatesEnumDirectory)
	p.SingletonTemplates = gen.MustLoadTemplates(templatesPackage, templatesSingletonDirectory)
}

type plugin interface {
	InitPlugin()
	RunPlugin()
}

func (p *Plugin) cmdGen(cmd *cobra.Command, args []string) {
	for _, i := range config.Config.Items {
		if p, ok := i.(plugin); ok {
			p.RunPlugin()
		}
	}
}

func (p *Plugin) RunPlugin() {
	if err := os.MkdirAll(config.Config.OutputPath, os.ModePerm); err != nil {
		log.Fatalf("Error creating output directory %s: %v", config.Config.OutputPath, err)
	}

	var models []*schema.Model
	for _, m := range config.Config.Schema.Models {
		models = append(models, m)
	}

	data := gen.BaseTemplateData()
	data["Models"] = models
	p.SingletonTemplates.ExecuteSingleton(data)

	for _, t := range config.Config.Schema.Types {
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

	for _, model := range config.Config.Schema.Models {
		if model.IsJoinModel {
			continue
		}

		data := gen.BaseTemplateData()
		data["Model"] = model
		data["Models"] = models
		p.ModelTemplates.Execute(data, model.Name+".go")
	}
}
