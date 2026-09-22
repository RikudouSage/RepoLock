package repo

import (
	"errors"
	"net/http"

	appErrors "go.chrastecky.dev/repolock/errors"
	apphttp "go.chrastecky.dev/repolock/http"
	"go.chrastecky.dev/repolock/http/dto"
	"go.chrastecky.dev/repolock/manager"
)

func (receiver *Controller) CreateRepository(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	body, err := apphttp.ParseBody[dto.CreateRepositoryRequest](request.Body)
	if err != nil {
		receiver.writer.WriteErrorResponse(err, writer)
		return
	}

	if err = body.Validate(); err != nil {
		receiver.writer.WriteErrorResponse(appErrors.NewUserFacingError(err.Error()), writer)
		return
	}

	hasAccess, err := receiver.accessManager.CanUserAccessOrganization(
		request.Context(),
		apphttp.MustGetUser(request.Context()).ID,
		body.OrganizationID,
		manager.AccessTypeWrite,
	)
	if err != nil {
		receiver.writer.WriteErrorResponse(err, writer)
		return
	}

	if !hasAccess {
		receiver.writer.WriteErrorResponse(appErrors.NewAccessDeniedError(), writer)
		return
	}

	repo, err := receiver.manager.CreateRepository(request.Context(), body.OrganizationID, body.URL, body.Name)
	if err != nil {
		if errors.Is(err, manager.ErrRepositoryAlreadyExists) {
			receiver.writer.WriteErrorResponse(
				appErrors.NewUserFacingErrorWithStatusCode("repository already exists", http.StatusConflict),
				writer,
			)
			return
		}

		receiver.writer.WriteErrorResponse(err, writer)
		return
	}

	receiver.writer.WriteResponse(http.StatusCreated, repo, writer)
}
