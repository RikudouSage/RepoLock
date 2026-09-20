package config

import (
	"go.chrastecky.dev/repolock/db"
	"go.chrastecky.dev/repolock/db/repo"
	"go.chrastecky.dev/repolock/manager"
	"go.chrastecky.dev/repolock/service"
	"go.chrastecky.dev/repolock/service/user"
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
		repo.NewPersonalAccessTokenRepository,

		manager.NewUserManager,
		manager.NewOrganizationManager,
		manager.NewPersonalAccessTokenManager,

		user.NewCreator,
		service.NewPasswordHasher,
		service.NewPasswordVerifier,
		service.NewFileSessionStore,
		service.NewRandomStringGenerator,
	)
}
