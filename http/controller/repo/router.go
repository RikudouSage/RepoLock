package repo

import (
	"github.com/go-chi/chi/v5"
	"go.chrastecky.dev/repolock/config/data"
	"go.chrastecky.dev/repolock/http/middleware"
)

func Router(
	requiresValidUserMiddleware middleware.RequiresValidUserMiddleware,
	controller *Controller,
) *data.MountableRouter {
	router := chi.NewRouter()
	router.Use(requiresValidUserMiddleware)

	router.Get("/", controller.ListRepositories)
	router.Get("/{id}", controller.GetRepository)
	router.Post("/", controller.CreateRepository)
	//router.Patch("/{id}", controller.UpdateRepository)
	//router.Delete("/{id}", controller.DeleteRepository)

	return data.AsMountableRouter("repositories", router)
}
