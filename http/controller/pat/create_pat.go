package pat

import (
	"net/http"

	apphttp "go.chrastecky.dev/repolock/http"
	"go.chrastecky.dev/repolock/http/dto"
	"go.chrastecky.dev/repolock/manager"
)

func (receiver *Controller) CreatePersonalAccessToken(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	body, err := apphttp.ParseBody[dto.CreatePersonalAccessTokenRequest](request.Body)
	if err != nil {
		receiver.writer.WriteErrorResponse(err, writer)
		return
	}

	user, _ := apphttp.GetUser(request.Context())
	_, token, err := receiver.manager.CreateForUser(request.Context(), user, &manager.PATConfig{
		Name:      body.Name,
		ExpiresAt: body.ExpiresAt,
	})
	if err != nil {
		receiver.writer.WriteErrorResponse(err, writer)
		return
	}

	receiver.writer.WriteResponse(http.StatusCreated, map[string]string{
		"token": token,
	}, writer)
}
