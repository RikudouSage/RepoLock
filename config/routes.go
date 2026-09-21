package config

import (
	"go.chrastecky.dev/repolock/http/controller/pat"
	"go.chrastecky.dev/repolock/http/controller/repo"
	"go.chrastecky.dev/repolock/http/controller/user"
	"go.uber.org/fx"
)

func asRouter(provider any) any {
	return fx.Annotate(
		provider,
		fx.ResultTags(`group:"routers"`),
	)
}

func provideControllers() fx.Option {
	return fx.Provide(
		user.NewController,
		pat.NewController,
		repo.NewController,
	)
}

func provideRoutes() fx.Option {
	return fx.Provide(
		asRouter(user.Router),
		asRouter(pat.Router),
		asRouter(repo.Router),
	)
}
