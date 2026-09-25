package org

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.chrastecky.dev/repolock/errors"
	apphttp "go.chrastecky.dev/repolock/http"
	"go.chrastecky.dev/repolock/manager"
)

func (receiver *Controller) GetOrg(writer http.ResponseWriter, request *http.Request) {
	id, err := uuid.Parse(chi.URLParam(request, "id"))
	if err != nil {
		receiver.writer.WriteErrorResponse(
			errors.NewUserFacingError("invalid id"),
			writer,
		)
		return
	}

	org, err := receiver.orgManager.FindByID(request.Context(), id)
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
		manager.AccessTypeRead,
	)
	if err != nil {
		receiver.writer.WriteErrorResponse(err, writer)
		return
	}

	if !hasAccess {
		receiver.writer.WriteNotFound(request, writer)
		return
	}

	receiver.writer.WriteResponse(http.StatusOK, org, writer)
}
