package bunnyid

import (
	"github.com/kernelpayments/sqlbunny/gen"
)

const (
	templatesPackage = "github.com/kernelpayments/sqlbunny/gen/bunnyid"
)

type Plugin struct {
	idTemplates        *gen.TemplateList
	singletonTemplates *gen.TemplateList
}

var _ gen.Plugin = &Plugin{}

func (*Plugin) IsConfigItem() {}

func (p *Plugin) BunnyPlugin() {
	p.idTemplates = gen.MustLoadTemplates(templatesPackage, "templates/id")
	p.singletonTemplates = gen.MustLoadTemplates(templatesPackage, "templates/singleton")

	gen.OnGen(p.gen)
}

func (p *Plugin) gen() {
	var idTypes []*IDType

	for _, t := range gen.Config.Schema.Types {
		switch t := t.(type) {
		case *IDType:
			data := gen.BaseTemplateData()
			data["IDType"] = t

			p.idTemplates.Execute(data, t.Name+".go")

			idTypes = append(idTypes, t)
		}
	}

	data := gen.BaseTemplateData()
	data["IDTypes"] = idTypes

	p.singletonTemplates.ExecuteSingleton(data)
}
