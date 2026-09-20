package middleware

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	apphttp "go.chrastecky.dev/repolock/http"
)

func SessionAuthMiddleware(
	sessionManager *scs.SessionManager,
) apphttp.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			userID := sessionManager.GetString(request.Context(), apphttp.SessionKeyUserID)
			if userID != "" {
				request = request.WithContext(
					apphttp.WithUserID(request.Context(), uuid.MustParse(userID)),
				)
			}

			next.ServeHTTP(writer, request)
		})
	}
}
