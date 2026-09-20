package middleware

import (
	"net/http"

	"go.chrastecky.dev/repolock/errors"
	apphttp "go.chrastecky.dev/repolock/http"
	"go.chrastecky.dev/repolock/http/response"
)

type RequiresValidUserMiddleware apphttp.Middleware

func ProvideRequiresValidUserMiddleware(
	responseWriter response.Writer,
) RequiresValidUserMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if _, ok := apphttp.GetUser(request.Context()); !ok {
				responseWriter.WriteErrorResponse(
					errors.NewUserFacingErrorWithStatusCode("unauthorized", http.StatusUnauthorized),
					writer,
				)
				return
			}

			next.ServeHTTP(writer, request)
		})
	}
}
