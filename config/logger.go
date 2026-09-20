package config

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func zapLifecycle(lifecycle fx.Lifecycle, log *zap.Logger) {
	lifecycle.Append(fx.Hook{
		OnStop: func(context.Context) error {
			_ = log.Sync()
			return nil
		},
	})
}

func newLogger() (*zap.Logger, error) {
	return zap.NewProduction()
}

func provideLogger() fx.Option {
	return fx.Module(
		"logger",
		fx.Provide(newLogger),
		fx.Invoke(zapLifecycle),
	)
}
