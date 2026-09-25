package org

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

	router.Get("/", controller.ListOrgs)
	router.Post("/", controller.Create)
	router.Get("/{id}", controller.GetOrg)
	router.Patch("/{id}", controller.UpdateOrg)

	return data.AsMountableRouter("/org", router)
}
