package pat

import (
	"github.com/go-chi/chi/v5"
	"go.chrastecky.dev/repolock/config/data"
	"go.chrastecky.dev/repolock/http/middleware"
)

func Router(
	controller *Controller,
	requiresValidUserMiddleware middleware.RequiresValidUserMiddleware,
) *data.MountableRouter {
	router := chi.NewRouter()
	router.Use(requiresValidUserMiddleware)

	router.Post("/", controller.CreatePersonalAccessToken)
	router.Delete("/{id}", controller.DeletePersonalAccessToken)

	return data.AsMountableRouter("pat", router)
}
