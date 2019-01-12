package core

import (
	"github.com/kernelpayments/sqlbunny/gen"
)

func Run(items ...gen.ConfigItem) {
	var items2 []gen.ConfigItem
	items2 = append(items2, &Plugin{})
	items2 = append(items2, items...)

	gen.Run(items2)
}
