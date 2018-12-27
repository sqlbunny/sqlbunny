package bunnyid

import (
	"github.com/kernelpayments/sqlbunny/config"
	"github.com/kernelpayments/sqlbunny/gen"
	"github.com/kernelpayments/sqlbunny/schema"
)

const (
	templatesPackage = "github.com/kernelpayments/sqlbunny/gen/bunnyid"

	templatesIDDirectory        = "templates/id"
	templatesSingletonDirectory = "templates/singleton"
)

type Plugin struct {
	IDTemplates        *gen.TemplateList
	SingletonTemplates *gen.TemplateList
}

func (*Plugin) IsConfigItem() {}

func (p *Plugin) InitPlugin() {
	p.IDTemplates = gen.MustLoadTemplates(templatesPackage, templatesIDDirectory)
	p.SingletonTemplates = gen.MustLoadTemplates(templatesPackage, templatesSingletonDirectory)
}

func (p *Plugin) RunPlugin() {
	var idTypes []*schema.IDType

	for _, t := range config.Config.Schema.Types {
		switch t := t.(type) {
		case *schema.IDType:
			data := &struct {
				*gen.TemplateData
				IDType *schema.IDType
			}{
				TemplateData: gen.BaseTemplateData(),
				IDType:       t,
			}

			p.IDTemplates.Execute(data, t.Name+".go")

			idTypes = append(idTypes, t)
		}
	}

	singletonData := &struct {
		*gen.TemplateData
		IDTypes []*schema.IDType
	}{
		TemplateData: gen.BaseTemplateData(),
		IDTypes:      idTypes,
	}

	p.SingletonTemplates.ExecuteSingleton(singletonData)
}
