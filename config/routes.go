package config

import "go.uber.org/fx"

func asRouter(provider any) any {
	return fx.Annotate(
		provider,
		fx.ResultTags(`group:"routers"`),
	)
}

func provideControllers() fx.Option {
	return fx.Provide()
}

func provideRoutes() fx.Option {
	return fx.Provide()
}
