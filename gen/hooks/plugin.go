package hooks

import (
	"bytes"

	"github.com/sqlbunny/sqlbunny/gen"
)

const (
	templatesPackage = "github.com/sqlbunny/sqlbunny/gen/hooks"
)

type Plugin struct {
}

var _ gen.Plugin = &Plugin{}

func (*Plugin) ConfigItem(ctx *gen.Context) {}

func (p *Plugin) BunnyPlugin() {
	gen.OnHook("after_delete_slice", p.hook(gen.MustLoadTemplate(templatesPackage, "templates/after_delete_slice.tpl")))
	gen.OnHook("after_delete", p.hook(gen.MustLoadTemplate(templatesPackage, "templates/after_delete.tpl")))
	gen.OnHook("after_insert", p.hook(gen.MustLoadTemplate(templatesPackage, "templates/after_insert.tpl")))
	gen.OnHook("after_select_slice", p.hook(gen.MustLoadTemplate(templatesPackage, "templates/after_select_slice.tpl")))
	gen.OnHook("after_select_slice_noreturn", p.hook(gen.MustLoadTemplate(templatesPackage, "templates/after_select_slice_noreturn.tpl")))
	gen.OnHook("after_select", p.hook(gen.MustLoadTemplate(templatesPackage, "templates/after_select.tpl")))
	gen.OnHook("after_update", p.hook(gen.MustLoadTemplate(templatesPackage, "templates/after_update.tpl")))
	gen.OnHook("before_delete_slice", p.hook(gen.MustLoadTemplate(templatesPackage, "templates/before_delete_slice.tpl")))
	gen.OnHook("before_delete", p.hook(gen.MustLoadTemplate(templatesPackage, "templates/before_delete.tpl")))
	gen.OnHook("before_insert", p.hook(gen.MustLoadTemplate(templatesPackage, "templates/before_insert.tpl")))
	gen.OnHook("before_update", p.hook(gen.MustLoadTemplate(templatesPackage, "templates/before_update.tpl")))
	gen.OnHook("model", p.modelHook(gen.MustLoadTemplate(templatesPackage, "templates/model.tpl")))
}

func copyData(m map[string]any) map[string]any {
	res := make(map[string]any)
	for k, v := range m {
		res[k] = v
	}

	return res
}

func (p *Plugin) hook(tpl *gen.TemplateList) gen.HookFunc {
	return func(buf *bytes.Buffer, data map[string]any, args ...any) {
		data2 := copyData(data)
		data2["Var"] = args[0]
		data2["Model"] = args[1]
		tpl.ExecuteBuf(data2, buf)
	}
}

func (p *Plugin) modelHook(tpl *gen.TemplateList) gen.HookFunc {
	return func(buf *bytes.Buffer, data map[string]any, args ...any) {
		tpl.ExecuteBuf(data, buf)
	}
}
