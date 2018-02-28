package gen

import (
	"bytes"
	"path/filepath"
	"regexp"
	"text/template"

	"github.com/KernelPay/sqlboiler/common"
	"github.com/pkg/errors"
)

var (
	// templateByteBuffer is re-used by all template construction to avoid
	// allocating more memory than is needed. This will later be a problem for
	// concurrency, address it then.
	templateByteBuffer = &bytes.Buffer{}

	rgxRemoveNumberedPrefix = regexp.MustCompile(`[0-9]+_`)
)

func (s *State) executeTemplates(data *templateData, templates *templateList, filename string) error {
	out := templateByteBuffer
	out.Reset()

	common.WriteFileDisclaimer(out)
	common.WritePackageName(out, s.Config.PkgName)

	for _, tplName := range templates.Templates() {
		if err := executeTemplate(out, templates.Template, tplName, data); err != nil {
			return err
		}
	}

	if err := common.WriteFile(s.Config.OutFolder, filename, out); err != nil {
		return err
	}

	return nil
}

func (s *State) executeSingletonTemplates(data *templateData, templates *templateList) error {
	out := templateByteBuffer
	for _, tplName := range templates.Templates() {
		out.Reset()

		fName := tplName
		ext := filepath.Ext(fName)
		fName = rgxRemoveNumberedPrefix.ReplaceAllString(fName[:len(fName)-len(ext)], "")

		common.WriteFileDisclaimer(out)
		common.WritePackageName(out, s.Config.PkgName)

		if err := executeTemplate(out, templates.Template, tplName, data); err != nil {
			return err
		}

		if err := common.WriteFile(s.Config.OutFolder, fName+".go", out); err != nil {
			return err
		}
	}

	return nil
}

func (s *State) generateTestMainOutput(state *State, data *templateData) error {
	if state.TestMainTemplate == nil {
		return errors.New("No TestMain template located for generation")
	}

	out := templateByteBuffer
	out.Reset()

	common.WriteFileDisclaimer(out)
	common.WritePackageName(out, state.Config.PkgName)

	if err := executeTemplate(out, state.TestMainTemplate, state.TestMainTemplate.Name(), data); err != nil {
		return err
	}

	if err := common.WriteFile(state.Config.OutFolder, "main_test.go", out); err != nil {
		return err
	}

	return nil
}

// executeTemplate takes a template and returns the output of the template
// execution.
func executeTemplate(buf *bytes.Buffer, t *template.Template, name string, data *templateData) error {
	if err := t.ExecuteTemplate(buf, name, data); err != nil {
		return errors.Wrapf(err, "failed to execute template: %s", name)
	}
	return nil
}
