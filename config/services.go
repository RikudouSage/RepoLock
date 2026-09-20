package config

import (
	"go.chrastecky.dev/repolock/db"
	"go.chrastecky.dev/repolock/db/repo"
	"go.uber.org/fx"
)

func provideServices() fx.Option {
	return fx.Provide(
		db.NewQueryFormatter,
		db.NewMapper,

		repo.NewUserRepository,
	)
}
