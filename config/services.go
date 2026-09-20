package config

import (
	"go.uber.org/fx"
)

func provideServices() fx.Option {
	return fx.Provide()
}
