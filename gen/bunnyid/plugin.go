package bunnyid

import (
	"bytes"

	"github.com/kernelpayments/sqlbunny/gen"
	"github.com/kernelpayments/sqlbunny/schema"
)

const (
	templatesPackage = "github.com/kernelpayments/sqlbunny/gen/bunnyid"
)

type Plugin struct {
	idTemplates        *gen.TemplateList
	modelTemplates     *gen.TemplateList
	singletonTemplates *gen.TemplateList
}

var _ gen.Plugin = &Plugin{}

func (*Plugin) IsConfigItem() {}

func (p *Plugin) BunnyPlugin() {
	gen.TemplateFunctions["bunnyid_IsStandardModel"] = p.IsStandardModel

	p.idTemplates = gen.MustLoadTemplates(templatesPackage, "templates/id")
	p.modelTemplates = gen.MustLoadTemplates(templatesPackage, "templates/model")
	p.singletonTemplates = gen.MustLoadTemplates(templatesPackage, "templates/singleton")

	gen.OnGen(p.gen)
	gen.OnHook("model", p.modelHook)
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

func (p *Plugin) modelHook(buf *bytes.Buffer, data map[string]interface{}, args ...interface{}) {
	p.modelTemplates.ExecuteBuf(data, buf)
}

func (p *Plugin) IsStandardModel(m *schema.Model) bool {
	for _, c := range m.Fields {
		_, ok := c.Type.(*IDType)
		if c.Name == "id" && ok {
			return true
		}
	}
	return false
}
