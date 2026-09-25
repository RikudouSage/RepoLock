package permissions

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	appErrors "go.chrastecky.dev/repolock/errors"
	apphttp "go.chrastecky.dev/repolock/http"
	"go.chrastecky.dev/repolock/manager"
)

func (receiver *Controller) GetRepositoryPermissions(writer http.ResponseWriter, request *http.Request) {
	repoID, err := uuid.Parse(chi.URLParam(request, "repoID"))
	if err != nil {
		receiver.writer.WriteErrorResponse(
			appErrors.NewUserFacingError("invalid id"),
			writer,
		)
		return
	}

	user := apphttp.MustGetUser(request.Context())
	repo, err := receiver.repoManager.GetRepoForUser(request.Context(), user, repoID)
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
		repo.ID,
		manager.AccessTypeManage,
	)
	if err != nil {
		receiver.writer.WriteErrorResponse(err, writer)
		return
	}
	if !hasAccess {
		receiver.writer.WriteErrorResponse(appErrors.NewAccessDeniedError(), writer)
		return
	}

	perms, err := receiver.repoPermissionManager.GetPermissions(request.Context(), repo)
	if err != nil {
		receiver.writer.WriteErrorResponse(err, writer)
		return
	}

	receiver.writer.WriteResponse(http.StatusOK, perms, writer)
}
