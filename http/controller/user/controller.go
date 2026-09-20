package user

import (
	"github.com/alexedwards/scs/v2"
	"go.chrastecky.dev/repolock/http/response"
	"go.chrastecky.dev/repolock/manager"
	"go.chrastecky.dev/repolock/service/user"
	"go.uber.org/zap"
)

type Controller struct {
	manager        manager.User
	responseWriter response.Writer
	userCreator    user.Creator
	logger         *zap.Logger
	sessionManager *scs.SessionManager
}

func NewController(
	responseWriter response.Writer,
	manager manager.User,
	userCreator user.Creator,
	logger *zap.Logger,
	sessionManager *scs.SessionManager,
) *Controller {
	return &Controller{
		responseWriter: responseWriter,
		manager:        manager,
		userCreator:    userCreator,
		logger:         logger,
		sessionManager: sessionManager,
	}
}
