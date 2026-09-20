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

func (receiver *Controller) UserPasswordLogin(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	body, err := apphttp.ParseBody[dto.UserPasswordLogin](request.Body)
	if err != nil {
		receiver.logger.Info("invalid body", zap.Error(err))
		receiver.responseWriter.WriteErrorResponse(
			appErrors.NewUserFacingError("invalid request body"),
			writer,
		)
		return
	}

	user, err := receiver.manager.FindByUsernameAndPassword(request.Context(), body.Username, body.Password)
	if err != nil {
		if errors.Is(err, manager.ErrUserNotFound) || errors.Is(err, manager.ErrPasswordNotEnabledForUser) || errors.Is(err, manager.ErrInvalidPassword) {
			receiver.responseWriter.WriteErrorResponse(
				appErrors.NewUserFacingError("user not found or invalid password"),
				writer,
			)
			return
		}

		receiver.responseWriter.WriteErrorResponse(err, writer)
		return
	}

	receiver.sessionManager.Put(request.Context(), apphttp.SessionKeyUserID, user.ID.String())
	receiver.responseWriter.WriteResponse(http.StatusNoContent, nil, writer)
}
