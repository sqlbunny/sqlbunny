package gen

import (
	"bytes"
	"log"
	"path/filepath"
	"text/template"
)

var (
	// templateByteBuffer is re-used by all template construction to avoid
	// allocating more memory than is needed. This will later be a problem for
	// concurrency, address it then.
	templateByteBuffer      = &bytes.Buffer{}
	templateByteBufferInner = &bytes.Buffer{}
)

func (t *TemplateList) ExecuteBuf(data map[string]interface{}, buf *bytes.Buffer) {
	for _, tplName := range t.Templates() {
		executeTemplate(buf, t.Template, tplName, data)
	}
}

func (t *TemplateList) Execute(data map[string]interface{}, filename string) {
	resetImports()
	innerOut := templateByteBufferInner
	innerOut.Reset()

	for _, tplName := range t.Templates() {
		executeTemplate(innerOut, t.Template, tplName, data)
	}

	out := templateByteBuffer
	out.Reset()

	WriteFileDisclaimer(out)
	WritePackageName(out, Config.PkgName)
	WriteImports(out, imports)
	out.Write(innerOut.Bytes())

	WriteFile(Config.OutputPath, filename, out.Bytes())
}

func (t *TemplateList) ExecuteSingleton(data map[string]interface{}) {
	out := templateByteBuffer
	for _, tplName := range t.Templates() {
		out.Reset()

		fName := tplName
		ext := filepath.Ext(fName)
		fName = fName[:len(fName)-len(ext)]

		WriteFileDisclaimer(out)
		WritePackageName(out, Config.PkgName)

		executeTemplate(out, t.Template, tplName, data)
		WriteFile(Config.OutputPath, fName+".go", out.Bytes())
	}
}

// executeTemplate takes a template and returns the output of the template
// execution.
func executeTemplate(buf *bytes.Buffer, t *template.Template, name string, data interface{}) {
	if err := t.ExecuteTemplate(buf, name, data); err != nil {
		log.Fatalf("failed to execute template %s: %v", name, err)
	}
}
