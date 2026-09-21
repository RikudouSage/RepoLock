package repo

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.chrastecky.dev/repolock/errors"
	apphttp "go.chrastecky.dev/repolock/http"
)

func (receiver *Controller) GetRepository(writer http.ResponseWriter, request *http.Request) {
	id, err := uuid.Parse(chi.URLParam(request, "id"))
	if err != nil {
		receiver.writer.WriteErrorResponse(errors.NewUserFacingError("invalid id"), writer)
		return
	}

	repo, err := receiver.manager.GetRepoForUser(request.Context(), apphttp.MustGetUser(request.Context()), id)
	if err != nil {
		receiver.writer.WriteErrorResponse(err, writer)
		return
	}

	result, err := json.Marshal(repo)
	if err != nil {
		receiver.writer.WriteErrorResponse(err, writer)
		return
	}

	receiver.writer.WriteResponse(http.StatusOK, result, writer)
}
