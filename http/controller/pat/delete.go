package pat

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.chrastecky.dev/repolock/errors"
)

func (receiver *Controller) DeletePersonalAccessToken(writer http.ResponseWriter, request *http.Request) {
	id, err := uuid.Parse(chi.URLParam(request, "id"))
	if err != nil {
		receiver.writer.WriteErrorResponse(
			errors.NewUserFacingError("invalid pat id"),
			writer,
		)
		return
	}

	pat, err := receiver.manager.FindByID(request.Context(), id)
	if err != nil {
		receiver.writer.WriteErrorResponse(err, writer)
		return
	}

	if pat == nil {
		receiver.writer.WriteNotFound(request, writer)
		return
	}

	if err = receiver.manager.DeleteToken(request.Context(), pat); err != nil {
		receiver.writer.WriteErrorResponse(err, writer)
	}

	receiver.writer.WriteResponse(http.StatusNoContent, nil, writer)
}
