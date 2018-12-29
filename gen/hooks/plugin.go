package hooks

import (
	"bytes"

	"github.com/kernelpayments/sqlbunny/gen"
)

const (
	templatesPackage = "github.com/kernelpayments/sqlbunny/gen/hooks"
)

type Plugin struct {
}

var _ gen.Plugin = &Plugin{}

func (*Plugin) IsConfigItem() {}

func (p *Plugin) BunnyPlugin() {
	gen.OnHook("after_delete_slice", p.hook(gen.MustLoadTemplate(templatesPackage, "templates/after_delete_slice.tpl")))
	gen.OnHook("after_delete", p.hook(gen.MustLoadTemplate(templatesPackage, "templates/after_delete.tpl")))
	gen.OnHook("after_insert", p.hook(gen.MustLoadTemplate(templatesPackage, "templates/after_insert.tpl")))
	gen.OnHook("after_select_slice", p.hook(gen.MustLoadTemplate(templatesPackage, "templates/after_select_slice.tpl")))
	gen.OnHook("after_select_slice_noreturn", p.hook(gen.MustLoadTemplate(templatesPackage, "templates/after_select_slice_noreturn.tpl")))
	gen.OnHook("after_select", p.hook(gen.MustLoadTemplate(templatesPackage, "templates/after_select.tpl")))
	gen.OnHook("after_update", p.hook(gen.MustLoadTemplate(templatesPackage, "templates/after_update.tpl")))
	gen.OnHook("after_upsert", p.hook(gen.MustLoadTemplate(templatesPackage, "templates/after_upsert.tpl")))
	gen.OnHook("before_delete_slice", p.hook(gen.MustLoadTemplate(templatesPackage, "templates/before_delete_slice.tpl")))
	gen.OnHook("before_delete", p.hook(gen.MustLoadTemplate(templatesPackage, "templates/before_delete.tpl")))
	gen.OnHook("before_insert", p.hook(gen.MustLoadTemplate(templatesPackage, "templates/before_insert.tpl")))
	gen.OnHook("before_update", p.hook(gen.MustLoadTemplate(templatesPackage, "templates/before_update.tpl")))
	gen.OnHook("before_upsert", p.hook(gen.MustLoadTemplate(templatesPackage, "templates/before_upsert.tpl")))
	gen.OnHook("model", p.modelHook(gen.MustLoadTemplate(templatesPackage, "templates/model.tpl")))
}

func copyData(m map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range m {
		res[k] = v
	}

	return res
}

func (p *Plugin) hook(tpl *gen.TemplateList) gen.HookFunc {
	return func(buf *bytes.Buffer, data map[string]interface{}, args ...interface{}) {
		data2 := copyData(data)
		data2["Var"] = args[0]
		data2["Model"] = args[1]
		tpl.ExecuteBuf(data2, buf)
	}
}

func (p *Plugin) modelHook(tpl *gen.TemplateList) gen.HookFunc {
	return func(buf *bytes.Buffer, data map[string]interface{}, args ...interface{}) {
		tpl.ExecuteBuf(data, buf)
	}
}
