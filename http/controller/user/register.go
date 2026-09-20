package user

import (
	"errors"
	"net/http"

	appErrors "go.chrastecky.dev/repolock/errors"
	apphttp "go.chrastecky.dev/repolock/http"
	"go.chrastecky.dev/repolock/http/dto"
	"go.chrastecky.dev/repolock/manager"
	"go.uber.org/zap"
)

func (receiver *Controller) Register(writer http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	body, err := apphttp.ParseBody[dto.UserRegister](req.Body)

	if err != nil {
		receiver.logger.Info("invalid body", zap.Error(err))
		receiver.responseWriter.WriteErrorResponse(
			appErrors.NewUserFacingError("invalid request body"),
			writer,
		)
		return
	}

	_, err = receiver.userCreator.Create(req.Context(), body)
	if err != nil {
		if errors.Is(err, manager.ErrUserAlreadyExists) {
			receiver.responseWriter.WriteErrorResponse(
				appErrors.NewUserFacingErrorWithStatusCode("the user already exists", http.StatusConflict),
				writer,
			)
			return
		}
		receiver.responseWriter.WriteErrorResponse(err, writer)
		return
	}

	receiver.responseWriter.WriteResponse(http.StatusNoContent, nil, writer)
}
