package user

import (
	"net/http"

	apphttp "go.chrastecky.dev/repolock/http"
)

func (receiver *Controller) GetIdentities(writer http.ResponseWriter, request *http.Request) {
	currentUser := apphttp.MustGetUser(request.Context())
	identities, err := receiver.identityManager.GetForUser(request.Context(), currentUser.ID)
	if err != nil {
		receiver.writer.WriteErrorResponse(err, writer)
		return
	}

	receiver.writer.WriteResponse(http.StatusOK, identities, writer)
}
