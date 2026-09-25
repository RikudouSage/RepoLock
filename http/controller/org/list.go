package org

import (
	"net/http"

	apphttp "go.chrastecky.dev/repolock/http"
)

func (receiver *Controller) ListOrgs(writer http.ResponseWriter, request *http.Request) {
	currentUser := apphttp.MustGetUser(request.Context())
	orgs, err := receiver.orgManager.GetForUser(request.Context(), currentUser.ID)
	if err != nil {
		receiver.writer.WriteErrorResponse(err, writer)
		return
	}

	receiver.writer.WriteResponse(http.StatusOK, orgs, writer)
}
