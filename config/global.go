package config

import (
	"go.chrastecky.dev/repolock/config/data"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/fx"
)

func newGlobalConfig() *data.GlobalConfig {
	var cfg data.GlobalConfig
	envconfig.MustProcess("app", &cfg)
	return &cfg
}

func provideGlobalConfig() fx.Option {
	return fx.Provide(newGlobalConfig)
}
