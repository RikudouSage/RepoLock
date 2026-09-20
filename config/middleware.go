package config

import (
	"fmt"

	apphttp "go.chrastecky.dev/repolock/http"
	appMiddleware "go.chrastecky.dev/repolock/http/middleware"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/fx"
)

var order uint

type orderedMiddleware struct {
	Order      uint
	Middleware apphttp.Middleware
}

type middlewareOut struct {
	fx.Out

	Item orderedMiddleware `group:"middleware"`
}

func provideMiddlewareProviderWithName(provider any, name string) fx.Option {
	defer func() { order++ }()
	order := order

	tag := fmt.Sprintf(`name:"http_middleware_%d"`, order)

	options := []fx.Option{
		fx.Provide(fx.Annotate(provider, fx.ResultTags(tag))),
		fx.Provide(fx.Annotate(func(middlewareInstance apphttp.Middleware) middlewareOut {
			return middlewareOut{Item: orderedMiddleware{
				Order:      order,
				Middleware: middlewareInstance,
			}}
		}, fx.ParamTags(tag))),
	}

	if name != "" {
		aliasTag := fmt.Sprintf(`name:"http_middleware_%s"`, name)

		options = append(options, fx.Provide(fx.Annotate(
			func(middlewareInstance apphttp.Middleware) apphttp.Middleware {
				return middlewareInstance
			},
			fx.ParamTags(tag),
			fx.ResultTags(aliasTag),
		)))
	}

	return fx.Options(options...)
}

func provideMiddlewareProvider(provider any) fx.Option {
	return provideMiddlewareProviderWithName(provider, "")
}

func provideMiddleware(middleware apphttp.Middleware) fx.Option {
	return provideMiddlewareProvider(func() apphttp.Middleware {
		return middleware
	})
}

func provideNamedOptionalMiddleware(provider any, name string) fx.Option {
	return fx.Provide(fx.Annotate(provider, fx.ResultTags(fmt.Sprintf(`name:"%s"`, name))))
}

func provideMiddlewares() fx.Option {
	return fx.Module(
		"middlewares",
		// for all routes
		provideMiddleware(middleware.RequestID),
		provideMiddleware(middleware.RealIP),
		provideMiddleware(middleware.Recoverer),
		provideMiddlewareProvider(appMiddleware.Logger),
		provideMiddleware(middleware.Compress(5)),
		provideMiddleware(middleware.GetHead),
		provideMiddlewareProvider(appMiddleware.SessionAuthMiddleware),
	)
}
