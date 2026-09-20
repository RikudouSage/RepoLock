package middleware

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"go.chrastecky.dev/repolock/entity"
	apphttp "go.chrastecky.dev/repolock/http"
	"go.chrastecky.dev/repolock/manager"
	"go.chrastecky.dev/repolock/types"
	"go.uber.org/zap"
)

func ResolveUserMiddleware(
	userManager manager.User,
	logger *zap.Logger,
) apphttp.Middleware {
	// todo handle eviction
	cache := types.NewSyncMap[uuid.UUID, *entity.User]()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if userID, ok := apphttp.GetUserID(request.Context()); ok {
				user, ok := cache.Get(userID)
				if !ok {
					var err error
					user, err = userManager.FindByID(request.Context(), userID)
					cache.Set(userID, user)
					if err != nil && !errors.Is(err, manager.ErrUserNotFound) {
						logger.Error("failed finding user by id", zap.Error(err))
					}
				}

				if user != nil {
					request = request.WithContext(
						apphttp.WithUser(request.Context(), user),
					)
				}
			}

			next.ServeHTTP(writer, request)
		})
	}
}
