package pat

import (
	"go.chrastecky.dev/repolock/http/response"
	"go.chrastecky.dev/repolock/manager"
)

type Controller struct {
	writer  response.Writer
	manager manager.PersonalAccessToken
}

func NewController(
	writer response.Writer,
	manager manager.PersonalAccessToken,
) *Controller {
	return &Controller{
		writer:  writer,
		manager: manager,
	}
}
