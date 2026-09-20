package user

import (
	"github.com/go-chi/chi/v5"
	"go.chrastecky.dev/repolock/config/data"
)

func Router(
	controller *Controller,
	config *data.GlobalConfig,
) *data.MountableRouter {
	router := chi.NewRouter()

	if config.RegistrationsEnabled {
		router.Post("/register", controller.Register)
	}
	if config.PasswordLoginEnabled {
		router.Post("/login/password", controller.UserPasswordLogin)
	}

	return data.AsMountableRouter("users", router)
}
