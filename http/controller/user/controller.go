package user

import (
	"github.com/alexedwards/scs/v2"
	"go.chrastecky.dev/repolock/http/response"
	"go.chrastecky.dev/repolock/manager"
	"go.chrastecky.dev/repolock/service/user"
	"go.uber.org/zap"
)

type Controller struct {
	userManager     manager.User
	identityManager manager.VCSIdentity
	userCreator     user.Creator

	writer         response.Writer
	sessionManager *scs.SessionManager
	logger         *zap.Logger
}

func NewController(
	responseWriter response.Writer,
	userManager manager.User,
	identityManager manager.VCSIdentity,
	userCreator user.Creator,
	logger *zap.Logger,
	sessionManager *scs.SessionManager,
) *Controller {
	return &Controller{
		writer:          responseWriter,
		userManager:     userManager,
		identityManager: identityManager,
		userCreator:     userCreator,
		logger:          logger,
		sessionManager:  sessionManager,
	}
}
