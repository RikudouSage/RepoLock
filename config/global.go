package config

import (
	"go.chrastecky.dev/repolock/config/data"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/fx"
)

func newGlobalConfig() (*data.GlobalConfig, error) {
	var cfg data.GlobalConfig
	envconfig.MustProcess("app", &cfg)

	if err := cfg.Normalize(); err != nil {
		return nil, err
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func provideGlobalConfig() fx.Option {
	return fx.Provide(newGlobalConfig)
}
