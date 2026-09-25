package org

import (
	"go.chrastecky.dev/repolock/http/response"
	"go.chrastecky.dev/repolock/manager"
)

type Controller struct {
	writer     response.Writer
	orgManager manager.Organization
}

func NewController(
	writer response.Writer,
	orgManager manager.Organization,
) *Controller {
	return &Controller{
		writer:     writer,
		orgManager: orgManager,
	}
}
