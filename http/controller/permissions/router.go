package permissions

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

	router.Route("/org", func(router chi.Router) {
		router.Get("/{orgID}", controller.GetOrganizationPermissions)
		router.Post("/{orgID}", controller.CreateOrganizationPermission)
		router.Get("/{orgID}/{userID}", controller.GetUserOrganizationPermission)
		router.Delete("/{orgID}/{userID}", controller.DeleteUserOrganizationPermission)
		router.Patch("/{orgID}/{userID}", controller.UpdateUserOrganizationPermission)
	})

	router.Route("/repo", func(router chi.Router) {
		router.Get("/{repoID}", controller.GetRepositoryPermissions)
		//router.Post("/{repoID}", controller.CreateRepositoryPermission)
		//router.Get("/{repoID}/{userID}", controller.GetUserRepositoryPermission)
		//router.Delete("/{repoID}/{userID}", controller.DeleteUserRepositoryPermission)
		//router.Patch("/{repoID}/{userID}", controller.UpdateUserRepositoryPermission)
	})

	return data.AsMountableRouter("permissions", router)
}
