package user

import (
	"go.chrastecky.dev/repolock/http/response"
	"go.chrastecky.dev/repolock/manager"
	"go.chrastecky.dev/repolock/service"
	"go.uber.org/zap"
)

type Controller struct {
	manager        manager.User
	responseWriter response.Writer
	userCreator    service.UserCreator
	logger         *zap.Logger
}

func NewController(
	responseWriter response.Writer,
	manager manager.User,
	userCreator service.UserCreator,
	logger *zap.Logger,
) *Controller {
	return &Controller{
		responseWriter: responseWriter,
		manager:        manager,
		userCreator:    userCreator,
		logger:         logger,
	}
}
