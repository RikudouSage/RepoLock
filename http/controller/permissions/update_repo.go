package permissions

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	appErrors "go.chrastecky.dev/repolock/errors"
	apphttp "go.chrastecky.dev/repolock/http"
	"go.chrastecky.dev/repolock/http/dto"
	"go.chrastecky.dev/repolock/manager"
)

func (receiver *Controller) UpdateUserRepositoryPermission(writer http.ResponseWriter, request *http.Request) {
	repoID, err := uuid.Parse(chi.URLParam(request, "repoID"))
	if err != nil {
		receiver.writer.WriteErrorResponse(appErrors.NewUserFacingError("invalid repo id"), writer)
		return
	}
	userID, err := uuid.Parse(chi.URLParam(request, "userID"))
	if err != nil {
		receiver.writer.WriteErrorResponse(appErrors.NewUserFacingError("invalid user id"), writer)
		return
	}
	defer request.Body.Close()
	body, err := apphttp.ParseBody[dto.UpdateRepositoryPermission](request.Body)
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

	hasAccess, err := receiver.accessManager.CanUserAccessRepository(
		request.Context(),
		currentUser.ID,
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

	perm, err := receiver.repoPermissionManager.GetUserPermission(request.Context(), repo.ID, userID)
	if err != nil {
		receiver.writer.WriteErrorResponse(err, writer)
		return
	}
	if perm == nil {
		receiver.writer.WriteNotFound(request, writer)
		return
	}

	perm.Permission = body.Permission

	if err = receiver.repoPermissionManager.Update(
		request.Context(),
		perm,
	); err != nil {
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
