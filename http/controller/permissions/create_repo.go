package permissions

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.chrastecky.dev/repolock/entity"
	appErrors "go.chrastecky.dev/repolock/errors"
	apphttp "go.chrastecky.dev/repolock/http"
	"go.chrastecky.dev/repolock/http/dto"
	"go.chrastecky.dev/repolock/manager"
)

func (receiver *Controller) CreateRepositoryPermission(writer http.ResponseWriter, request *http.Request) {
	repoID, err := uuid.Parse(chi.URLParam(request, "repoID"))
	if err != nil {
		receiver.writer.WriteErrorResponse(appErrors.NewUserFacingError("invalid id"), writer)
		return
	}

	defer request.Body.Close()
	body, err := apphttp.ParseBody[dto.CreateRepositoryPermission](request.Body)

	if err != nil {
		receiver.writer.WriteErrorResponse(appErrors.NewUserFacingError("invalid request data"), writer)
		return
	}

	currentUser := apphttp.MustGetUser(request.Context())
	repo, err := receiver.repoManager.GetRepoForUser(request.Context(), currentUser, repoID)
	if err != nil {
		receiver.writer.WriteErrorResponse(err, writer)
		return
	}

	if repo == nil {
		receiver.writer.WriteNotFound(request, writer)
		return
	}

	hasWriteAccess, err := receiver.accessManager.CanUserAccessRepository(
		request.Context(),
		currentUser.ID,
		repo.ID,
		manager.AccessTypeManage,
	)
	if err != nil {
		receiver.writer.WriteErrorResponse(err, writer)
		return
	}

	if !hasWriteAccess {
		receiver.writer.WriteErrorResponse(appErrors.NewAccessDeniedError(), writer)
		return
	}

	perm := &entity.RepositoryPermission{
		RepositoryID: repo.ID,
		UserID:       body.UserID,
		Permission:   body.Permission,
	}
	err = receiver.repoPermissionManager.Create(request.Context(), perm)
	if err != nil {
		if errors.Is(err, manager.ErrPermissionAlreadyExists) {
			receiver.writer.WriteErrorResponse(appErrors.NewUserFacingErrorWithStatusCode(
				"the user already has permission for this repository",
				http.StatusConflict,
			), writer)
			return
		}

		receiver.writer.WriteErrorResponse(err, writer)
		return
	}

	receiver.writer.WriteResponse(http.StatusOK, &dto.RepositoryPermission{
		UserID:           perm.UserID,
		RepositoryID:     perm.RepositoryID,
		Permission:       perm.Permission,
		FromOrganization: false,
	}, writer)
}
