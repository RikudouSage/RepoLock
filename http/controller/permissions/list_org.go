package permissions

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.chrastecky.dev/repolock/entity"
	appErrors "go.chrastecky.dev/repolock/errors"
	apphttp "go.chrastecky.dev/repolock/http"
	"go.chrastecky.dev/repolock/http/dto"
	"go.chrastecky.dev/repolock/manager"
)

func (receiver *Controller) GetOrganizationPermissions(writer http.ResponseWriter, request *http.Request) {
	orgID, err := uuid.Parse(chi.URLParam(request, "orgID"))
	if err != nil {
		receiver.writer.WriteErrorResponse(
			appErrors.NewUserFacingError("invalid id"),
			writer,
		)
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

	user := apphttp.MustGetUser(request.Context())
	hasAccess, err := receiver.accessManager.CanUserAccessOrganization(
		request.Context(),
		user.ID,
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

	perms, err := receiver.orgPermissionManager.GetMemberships(request.Context(), org)
	if err != nil {
		receiver.writer.WriteErrorResponse(err, writer)
		return
	}

	result := lo.Map(perms, func(item *entity.OrganizationMembership, _ int) *dto.OrganizationPermission {
		return &dto.OrganizationPermission{
			OrganizationID: item.OrganizationID,
			UserID:         item.UserID,
			Permission:     item.Permission,
			Approved:       item.Approved,
		}
	})

	receiver.writer.WriteResponse(http.StatusOK, result, writer)
}
