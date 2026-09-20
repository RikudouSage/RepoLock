package main

import (
	"time"

	"go.chrastecky.dev/repolock/config"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	app := fx.New(
		config.Providers(),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.StartTimeout(30*time.Second),
	)

	app.Run()
}
