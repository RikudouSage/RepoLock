package user

import (
	"errors"
	"net/http"

	appErrors "go.chrastecky.dev/repolock/errors"
	apphttp "go.chrastecky.dev/repolock/http"
	"go.chrastecky.dev/repolock/http/dto"
	"go.chrastecky.dev/repolock/manager"
)

func (receiver *Controller) CreateIdentity(writer http.ResponseWriter, request *http.Request) {
	currentUser := apphttp.MustGetUser(request.Context())

	defer request.Body.Close()
	body, err := apphttp.ParseBody[dto.CreateIdentityRequest](request.Body)
	if err != nil {
		receiver.writer.WriteErrorResponse(
			appErrors.NewUserFacingError("invalid data"),
			writer,
		)
		return
	}

	item, err := receiver.identityManager.CreateForUser(request.Context(), body.Identity, currentUser.ID)
	if err != nil {
		if errors.Is(err, manager.ErrIdentityAlreadyExists) {
			receiver.writer.WriteErrorResponse(
				appErrors.NewUserFacingErrorWithStatusCode(
					"this identity already exists for this user",
					http.StatusConflict,
				),
				writer,
			)
			return
		}

		receiver.writer.WriteErrorResponse(err, writer)
		return
	}

	receiver.writer.WriteResponse(http.StatusOK, item, writer)
}
