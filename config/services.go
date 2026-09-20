package config

import (
	"go.chrastecky.dev/repolock/db"
	"go.chrastecky.dev/repolock/db/repo"
	"go.chrastecky.dev/repolock/manager"
	"go.chrastecky.dev/repolock/service"
	"go.uber.org/fx"
)

func provideServices() fx.Option {
	return fx.Provide(
		db.NewQueryFormatter,
		db.NewMapper,
		db.NewTransactionCreator,

		repo.NewUserRepository,
		repo.NewOrganizationRepository,
		repo.NewOrganizationMembershipRepository,

		manager.NewUserManager,
		manager.NewOrganizationManager,

		service.NewUserCreator,
		service.NewPasswordHasher,
		service.NewPasswordVerifier,
	)
}
