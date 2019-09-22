package core

import "github.com/sqlbunny/sqlbunny/gen"

type Config struct {
	ModelsPackagePath string
	ModelsPackageName string
}

func (c *Config) ConfigItem(ctx *gen.Context) {
}

func (c *Config) BunnyConfig(s *gen.ConfigStruct) {
	if c.ModelsPackagePath != "" {
		s.ModelsPackagePath = c.ModelsPackagePath
	}
	if c.ModelsPackageName != "" {
		s.ModelsPackageName = c.ModelsPackageName
	}
}
