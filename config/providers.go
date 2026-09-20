package config

import "go.uber.org/fx"

func Providers() fx.Option {
	return fx.Options(
		provideLogger(),
		provideServices(),
		provideGlobalConfig(),
		provideMiddlewares(),
		provideHttp(),
		provideRoutes(),
		provideControllers(),
		provideCron(),
		provideDatabase(),
	)
}
