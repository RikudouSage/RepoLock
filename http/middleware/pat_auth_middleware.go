package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	apphttp "go.chrastecky.dev/repolock/http"
	"go.chrastecky.dev/repolock/manager"
	"go.uber.org/zap"
)

func PersonalAccessTokenAuthMiddleware(
	patManager manager.PersonalAccessToken,
	logger *zap.Logger,
) apphttp.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if authorization := request.Header.Get("Authorization"); len(authorization) > 0 {
				authorization = strings.TrimPrefix(authorization, "Bearer ")
				pat, err := patManager.FindByToken(request.Context(), authorization)
				if err != nil {
					logger.Info("failed to find personal access token", zap.Error(err))
				} else {
					request = request.WithContext(
						apphttp.WithUserID(request.Context(), pat.UserID),
					)
					go func() {
						ctx, cancel := context.WithTimeout(request.Context(), 10*time.Second)
						defer cancel()

						err := patManager.UpdateLastUsedAt(
							ctx,
							pat,
							time.Now(),
						)
						if err != nil {
							logger.Warn("failed to update last used at", zap.Error(err))
						}
					}()
				}
			}

			next.ServeHTTP(writer, request)
		})
	}
}
