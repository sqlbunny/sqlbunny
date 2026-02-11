package bunnyid

import (
	"github.com/sqlbunny/sqlbunny/gen"
)

const (
	templatesPackage = "github.com/sqlbunny/bunnyid/gen"
)

type Plugin struct {
	idTemplates        *gen.TemplateList
	singletonTemplates *gen.TemplateList
}

var _ gen.Plugin = &Plugin{}

func (*Plugin) ConfigItem(ctx *gen.Context) {}

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

			p.idTemplates.Execute(data, t.Name+".gen.go")

			idTypes = append(idTypes, t)
		}
	}

	data := gen.BaseTemplateData()
	data["IDTypes"] = idTypes

	p.singletonTemplates.ExecuteSingleton(data)
}
