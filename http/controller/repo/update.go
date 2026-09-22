package repo

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	appErrors "go.chrastecky.dev/repolock/errors"
	apphttp "go.chrastecky.dev/repolock/http"
	"go.chrastecky.dev/repolock/http/dto"
	"go.chrastecky.dev/repolock/manager"
)

func (receiver *Controller) UpdateRepository(writer http.ResponseWriter, request *http.Request) {
	repositoryID, err := uuid.Parse(chi.URLParam(request, "id"))
	if err != nil {
		receiver.writer.WriteErrorResponse(appErrors.NewUserFacingError("invalid id"), writer)
		return
	}

	defer request.Body.Close()
	body, err := apphttp.ParseBody[dto.UpdateRepositoryRequest](request.Body)
	if err != nil {
		receiver.writer.WriteErrorResponse(err, writer)
		return
	}

	user := apphttp.MustGetUser(request.Context())

	repo, err := receiver.manager.GetRepoForUser(request.Context(), user, repositoryID)
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
		user.ID,
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

	repo.Name = body.Name
	if err = receiver.manager.UpdateRepository(request.Context(), repo); err != nil {
		receiver.writer.WriteErrorResponse(err, writer)
		return
	}

	receiver.writer.WriteResponse(http.StatusOK, repo, writer)
}
