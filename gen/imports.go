package gen

import (
	"fmt"

	"github.com/KernelPay/sqlbunny/schema"
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

func templateTypesGo(t []schema.TypeGo) []string {
	r := make([]string, len(t))
	for i, j := range t {
		r[i] = templateTypeGo(j)
	}
	return r
}

func templateTypeGo(t schema.TypeGo) string {
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

func removeDuplicates(dedup []string) []string {
	if len(dedup) <= 1 {
		return dedup
	}

	for i := 0; i < len(dedup)-1; i++ {
		for j := i + 1; j < len(dedup); j++ {
			if dedup[i] != dedup[j] {
				continue
			}

			if j != len(dedup)-1 {
				dedup[j] = dedup[len(dedup)-1]
				j--
			}
			dedup = dedup[:len(dedup)-1]
		}
	}

	return dedup
}
