package user

import (
	"github.com/go-chi/chi/v5"
	"go.chrastecky.dev/repolock/config/data"
	"go.chrastecky.dev/repolock/http/middleware"
)

func Router(
	controller *Controller,
	config *data.GlobalConfig,
	requiresValidUserMiddleware middleware.RequiresValidUserMiddleware,
) *data.MountableRouter {
	router := chi.NewRouter()

	if config.RegistrationsEnabled {
		router.Post("/register", controller.Register)
	}
	if config.PasswordLoginEnabled {
		router.Post("/login/password", controller.UserPasswordLogin)
	}

	router.With(requiresValidUserMiddleware).Route("/identities", func(router chi.Router) {
		router.Get("/", controller.GetIdentities)
		router.Post("/", controller.CreateIdentity)
		router.Delete("/{identity}", controller.DeleteByIdentity)
		router.Delete("/by-id/{id}", controller.DeleteByID)
	})

	return data.AsMountableRouter("users", router)
}
