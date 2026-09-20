package data

import (
	"strings"

	"github.com/go-chi/chi/v5"
)

type MountableRouter struct {
	MountPoint       string
	Router           chi.Router
	RequireValidSite bool
}

func AsMountableRouter(prefix string, router chi.Router) *MountableRouter {
	return asMountableRouter(prefix, router, true)
}

func AsMountableRouterWithoutValidSite(prefix string, router chi.Router) *MountableRouter {
	return asMountableRouter(prefix, router, false)
}

func asMountableRouter(prefix string, router chi.Router, requireValidSite bool) *MountableRouter {
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}

	return &MountableRouter{
		MountPoint:       prefix,
		Router:           router,
		RequireValidSite: requireValidSite,
	}
}
