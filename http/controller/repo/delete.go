package repo

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	appErrors "go.chrastecky.dev/repolock/errors"
	apphttp "go.chrastecky.dev/repolock/http"
	"go.chrastecky.dev/repolock/manager"
)

func (receiver *Controller) DeleteRepository(writer http.ResponseWriter, request *http.Request) {
	repositoryID, err := uuid.Parse(chi.URLParam(request, "id"))
	if err != nil {
		receiver.writer.WriteErrorResponse(
			appErrors.NewUserFacingError("invalid id"),
			writer,
		)
		return
	}

	_, err = receiver.manager.GetRepoForUser(
		request.Context(),
		apphttp.MustGetUser(request.Context()),
		repositoryID,
	)
	if err != nil {
		if errors.Is(err, manager.ErrRepositoryNotFound) {
			receiver.writer.WriteNotFound(request, writer)
			return
		}

		receiver.writer.WriteErrorResponse(err, writer)
		return
	}

	hasAccess, err := receiver.accessManager.CanUserAccessRepository(
		request.Context(),
		apphttp.MustGetUser(request.Context()).ID,
		repositoryID,
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

	if err = receiver.manager.DeleteRepository(request.Context(), repositoryID); err != nil {
		receiver.writer.WriteErrorResponse(err, writer)
		return
	}

	receiver.writer.WriteResponse(http.StatusNoContent, nil, writer)
}
