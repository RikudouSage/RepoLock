package org

import (
	"net/http"

	"go.chrastecky.dev/repolock/entity"
	"go.chrastecky.dev/repolock/errors"
	apphttp "go.chrastecky.dev/repolock/http"
	"go.chrastecky.dev/repolock/http/dto"
	"go.uber.org/zap"
)

func (receiver *Controller) Create(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	body, err := apphttp.ParseBody[dto.CreateOrUpdateOrganizationRequest](request.Body)
	if err != nil {
		receiver.logger.Info("unable to parse body", zap.Error(err))
		receiver.writer.WriteErrorResponse(
			errors.NewUserFacingError("invalid body"),
			writer,
		)
		return
	}

	org := &entity.Organization{
		Name: body.Name,
	}
	currentUser := apphttp.MustGetUser(request.Context())
	if err = receiver.creator.CreateWithUser(request.Context(), org, currentUser); err != nil {
		receiver.writer.WriteErrorResponse(err, writer)
		return
	}

	receiver.writer.WriteResponse(http.StatusCreated, org, writer)
}
