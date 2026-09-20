package config

import (
	"context"
	"errors"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/alexedwards/scs/v2"
	"go.chrastecky.dev/repolock/config/data"
	"go.chrastecky.dev/repolock/http/response"
	"go.chrastecky.dev/repolock/service"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
)

type routerDi struct {
	fx.In

	SessionManager *scs.SessionManager
	Middlewares    []orderedMiddleware     `group:"middleware"`
	Routers        []*data.MountableRouter `group:"routers"`
}

func newChiRouter(in routerDi) *chi.Mux {
	allMiddlewareWrappers := in.Middlewares
	slices.SortFunc(allMiddlewareWrappers, func(a, b orderedMiddleware) int {
		if a.Order > b.Order {
			return 1
		}

		return -1
	})

	router := chi.NewRouter()
	router.Use(in.SessionManager.LoadAndSave)

	for _, middlewareInstance := range allMiddlewareWrappers {
		router.Use(middlewareInstance.Middleware)
	}

	for _, mountableRouter := range in.Routers {
		router.Mount(mountableRouter.MountPoint, mountableRouter.Router)
	}

	return router
}

func startHttpServer(lifecycle fx.Lifecycle, cfg *data.GlobalConfig, router *chi.Mux) {
	server := &http.Server{
		Addr:              ":" + strconv.FormatUint(uint64(cfg.Port), 10),
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			return server.Shutdown(ctx)
		},
	})
}

func newSessionManager(
	store service.FileSessionStore,
) *scs.SessionManager {
	sessionManager := scs.New()
	sessionManager.Lifetime = 365 * 24 * time.Hour
	sessionManager.IdleTimeout = 24 * time.Hour
	sessionManager.Store = store

	return sessionManager
}

func provideHttp() fx.Option {
	return fx.Module(
		"http",
		fx.Provide(response.NewWriter),
		fx.Provide(newChiRouter),
		fx.Provide(newSessionManager),
		fx.Invoke(startHttpServer),
	)
}
