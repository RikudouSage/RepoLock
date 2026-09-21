package repo

import (
	"go.chrastecky.dev/repolock/http/response"
	"go.chrastecky.dev/repolock/manager"
)

type Controller struct {
	writer        response.Writer
	manager       manager.Repository
	accessManager manager.Access
}

func NewController(
	writer response.Writer,
	manager manager.Repository,
	accessManager manager.Access,
) *Controller {
	return &Controller{
		writer:        writer,
		manager:       manager,
		accessManager: accessManager,
	}
}
