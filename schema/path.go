package schema

import "strings"

type Path []string

func (p Path) DotName() string {
	return strings.Join(p, ".")
}

func (p Path) SQLName() string {
	return strings.Join(p, "__")
}

func (p Path) Equals(q Path) bool {
	if len(p) != len(q) {
		return false
	}
	for i := range p {
		if p[i] != q[i] {
			return false
		}
	}
	return true
}
