package repo

import (
	"net/http"

	apphttp "go.chrastecky.dev/repolock/http"
)

func (receiver *Controller) ListRepositories(writer http.ResponseWriter, request *http.Request) {
	repos, err := receiver.manager.GetReposForUser(request.Context(), apphttp.MustGetUser(request.Context()))
	if err != nil {
		receiver.writer.WriteErrorResponse(err, writer)
		return
	}

	receiver.writer.WriteResponse(http.StatusOK, repos, writer)
}
