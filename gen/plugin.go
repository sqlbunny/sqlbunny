package gen

import (
	"bytes"

	"github.com/spf13/cobra"
)

type Plugin interface {
	IsConfigItem()
	BunnyPlugin()
}

type HookFunc func(buf *bytes.Buffer, data map[string]interface{}, args ...interface{})

var (
	genFuncs  []func()
	hookFuncs = make(map[string][]HookFunc)
)

func OnGen(f func()) {
	genFuncs = append(genFuncs, f)
}

func AddCommand(cmds ...*cobra.Command) {
	rootCmd.AddCommand(cmds...)
}

func OnHook(name string, f HookFunc) {
	hookFuncs[name] = append(hookFuncs[name], f)
}

func hook(data map[string]interface{}, name string, args ...interface{}) string {
	var buf bytes.Buffer
	for _, f := range hookFuncs[name] {
		f(&buf, data, args...)
	}
	return buf.String()
}
