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

func (receiver *Controller) CreateOrganizationPermission(writer http.ResponseWriter, request *http.Request) {
	orgID, err := uuid.Parse(chi.URLParam(request, "orgID"))
	if err != nil {
		receiver.writer.WriteErrorResponse(appErrors.NewUserFacingError("invalid id"), writer)
		return
	}

	defer request.Body.Close()
	body, err := apphttp.ParseBody[dto.CreateOrganizationPermission](request.Body)

	if err != nil {
		receiver.writer.WriteErrorResponse(appErrors.NewUserFacingError("invalid request data"), writer)
		return
	}

	user := apphttp.MustGetUser(request.Context())
	org, err := receiver.orgManager.FindByID(request.Context(), orgID)
	if err != nil {
		receiver.writer.WriteErrorResponse(err, writer)
		return
	}

	if org == nil {
		receiver.writer.WriteNotFound(request, writer)
		return
	}

	hasWriteAccess, err := receiver.accessManager.CanUserAccessOrganization(
		request.Context(),
		user.ID,
		org.ID,
		manager.AccessTypeWrite,
	)
	if err != nil {
		receiver.writer.WriteErrorResponse(err, writer)
		return
	}

	if !hasWriteAccess && (body.Approved || !receiver.anyoneCanRequestJoin) {
		receiver.writer.WriteErrorResponse(appErrors.NewAccessDeniedError(), writer)
		return
	}

	membership := &entity.OrganizationMembership{
		OrganizationID: org.ID,
		UserID:         body.UserID,
		Permission:     body.Permission,
		Approved:       body.Approved,
	}
	err = receiver.orgPermissionManager.Create(request.Context(), membership)
	if err != nil {
		if errors.Is(err, manager.ErrOrganizationMembershipAlreadyExists) {
			receiver.writer.WriteErrorResponse(appErrors.NewUserFacingErrorWithStatusCode(
				"the user is already a member of this organization",
				http.StatusConflict,
			), writer)
			return
		}

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
