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

func (receiver *Controller) UpdateUserOrganizationPermission(writer http.ResponseWriter, request *http.Request) {
	orgID, err := uuid.Parse(chi.URLParam(request, "orgID"))
	if err != nil {
		receiver.writer.WriteErrorResponse(appErrors.NewUserFacingError("invalid org id"), writer)
		return
	}
	userID, err := uuid.Parse(chi.URLParam(request, "userID"))
	if err != nil {
		receiver.writer.WriteErrorResponse(appErrors.NewUserFacingError("invalid user id"), writer)
		return
	}
	defer request.Body.Close()
	body, err := apphttp.ParseBody[dto.UpdateOrganizationPermission](request.Body)
	if err != nil {
		receiver.writer.WriteErrorResponse(appErrors.NewUserFacingError("invalid request data"), writer)
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
		receiver.writer.WriteErrorResponse(appErrors.NewAccessDeniedError(), writer)
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

	membership.Approved = body.Approved
	membership.Permission = body.Permission

	if err = receiver.orgPermissionManager.UpdateMembership(
		request.Context(),
		membership,
	); err != nil {
		receiver.writer.WriteErrorResponse(err, writer)
		return
	}

	receiver.writer.WriteResponse(http.StatusOK, &dto.OrganizationPermission{
		OrganizationID: membership.OrganizationID,
		UserID:         membership.UserID,
		Permission:     membership.Permission,
		Approved:       membership.Approved,
	}, writer)
}
