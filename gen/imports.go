package gen

import (
	"fmt"

	"github.com/kernelpayments/sqlbunny/schema"
)

var (
	imports     map[string]string
	importCount int
)

func resetImports() {
	imports = make(map[string]string)
	importCount = 0
}

func templateImport(name string, pkg string) string {
	oldName, ok := imports[pkg]
	if ok && oldName != name {
		panic(fmt.Sprintf("package %s can't be imported with name %s, was already imported with name %s", pkg, name, oldName))
	}
	imports[pkg] = name
	return ""
}

func templateTypesGo(t []schema.GoType) []string {
	r := make([]string, len(t))
	for i, j := range t {
		r[i] = templateGoType(j)
	}
	return r
}

func templateGoType(t schema.GoType) string {
	if t.Pkg == "" {
		return t.Name
	}

	pkgName, ok := imports[t.Pkg]
	if !ok {
		pkgName = fmt.Sprintf("_import%02d", importCount)
		imports[t.Pkg] = pkgName
		importCount++
	}
	return pkgName + "." + t.Name
}
