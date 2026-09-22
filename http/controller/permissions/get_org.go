package permissions

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.chrastecky.dev/repolock/errors"
	apphttp "go.chrastecky.dev/repolock/http"
	"go.chrastecky.dev/repolock/http/dto"
	"go.chrastecky.dev/repolock/manager"
)

func (receiver *Controller) GetUserOrganizationPermission(writer http.ResponseWriter, request *http.Request) {
	orgID, err := uuid.Parse(chi.URLParam(request, "orgID"))
	if err != nil {
		receiver.writer.WriteErrorResponse(errors.NewUserFacingError("invalid org id"), writer)
		return
	}
	userID, err := uuid.Parse(chi.URLParam(request, "userID"))
	if err != nil {
		receiver.writer.WriteErrorResponse(errors.NewUserFacingError("invalid user id"), writer)
		return
	}

	org, err := receiver.orgManager.FindByID(request.Context(), orgID)
	if err != nil {
		receiver.writer.WriteErrorResponse(err, writer)
		return
	}
	if org == nil {
		receiver.writer.WriteNotFound(request, writer)
		return
	}

	currentUser := apphttp.MustGetUser(request.Context())
	hasAccess, err := receiver.accessManager.CanUserAccessOrganization(
		request.Context(),
		currentUser.ID,
		org.ID,
		manager.AccessTypeManage,
	)
	if err != nil {
		receiver.writer.WriteErrorResponse(err, writer)
		return
	}
	if !hasAccess {
		receiver.writer.WriteErrorResponse(errors.NewAccessDeniedError(), writer)
		return
	}

	membership, err := receiver.orgPermissionManager.GetUserMembership(request.Context(), org.ID, userID)
	if err != nil {
		receiver.writer.WriteErrorResponse(err, writer)
		return
	}
	if membership == nil {
		receiver.writer.WriteNotFound(request, writer)
		return
	}

	receiver.writer.WriteResponse(http.StatusOK, &dto.OrganizationPermission{
		OrganizationID: membership.OrganizationID,
		UserID:         membership.UserID,
		Permission:     membership.Permission,
		Approved:       membership.Approved,
	}, writer)
}
