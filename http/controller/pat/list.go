package pat

import (
	"net/http"

	apphttp "go.chrastecky.dev/repolock/http"
)

func (receiver *Controller) ListPersonalAccessTokens(writer http.ResponseWriter, request *http.Request) {
	items, err := receiver.manager.FindAllForUser(
		request.Context(),
		apphttp.MustGetUser(request.Context()),
	)
	if err != nil {
		receiver.writer.WriteErrorResponse(err, writer)
		return
	}

	receiver.writer.WriteResponse(http.StatusOK, items, writer)
}
