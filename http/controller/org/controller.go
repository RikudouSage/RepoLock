package org

import (
	"go.chrastecky.dev/repolock/http/response"
	"go.chrastecky.dev/repolock/manager"
	"go.chrastecky.dev/repolock/service/organization"
	"go.uber.org/zap"
)

type Controller struct {
	writer     response.Writer
	orgManager manager.Organization
	creator    organization.Creator
	logger     *zap.Logger
}

func NewController(
	writer response.Writer,
	orgManager manager.Organization,
	creator organization.Creator,
	logger *zap.Logger,
) *Controller {
	return &Controller{
		writer:     writer,
		orgManager: orgManager,
		logger:     logger,
		creator:    creator,
	}
}
