package user

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.chrastecky.dev/repolock/errors"
	apphttp "go.chrastecky.dev/repolock/http"
)

func (receiver *Controller) DeleteByIdentity(writer http.ResponseWriter, request *http.Request) {
	identityName := chi.URLParam(request, "identity")
	if identityName == "" {
		receiver.writer.WriteNotFound(request, writer)
		return
	}

	user := apphttp.MustGetUser(request.Context())
	identity, err := receiver.identityManager.GetForUserByIdentity(request.Context(), identityName, user.ID)
	if err != nil {
		receiver.writer.WriteErrorResponse(err, writer)
		return
	}

	if err = receiver.identityManager.Delete(request.Context(), identity); err != nil {
		receiver.writer.WriteErrorResponse(err, writer)
		return
	}

	receiver.writer.WriteResponse(http.StatusNoContent, nil, writer)
}

func (receiver *Controller) DeleteByID(writer http.ResponseWriter, request *http.Request) {
	identityID, err := uuid.Parse(chi.URLParam(request, "id"))
	if err != nil {
		receiver.writer.WriteErrorResponse(
			errors.NewUserFacingError("invalid ID"),
			writer,
		)
		return
	}

	user := apphttp.MustGetUser(request.Context())
	identity, err := receiver.identityManager.GetForUserByID(request.Context(), identityID, user.ID)
	if err != nil {
		receiver.writer.WriteErrorResponse(err, writer)
		return
	}

	if err = receiver.identityManager.Delete(request.Context(), identity); err != nil {
		receiver.writer.WriteErrorResponse(err, writer)
		return
	}

	receiver.writer.WriteResponse(http.StatusNoContent, nil, writer)
}
