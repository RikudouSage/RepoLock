package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"go.chrastecky.dev/repolock/config/data"
	appErrors "go.chrastecky.dev/repolock/errors"

	"github.com/samber/lo"
	"go.uber.org/zap"
)

type Writer interface {
	WriteErrorResponse(err error, writer http.ResponseWriter)
	WriteResponse(status int, body any, responseWriter http.ResponseWriter)
	WriteNotFound(request *http.Request, writer http.ResponseWriter)
}

type writer struct {
	logger *zap.Logger
	cfg    *data.GlobalConfig
}

func (receiver *writer) WriteNotFound(_ *http.Request, writer http.ResponseWriter) {
	receiver.WriteErrorResponse(appErrors.NewUserFacingErrorWithStatusCode(
		"Not found",
		http.StatusNotFound,
	), writer)
}

func NewWriter(logger *zap.Logger, cfg *data.GlobalConfig) Writer {
	return &writer{
		logger: logger,
		cfg:    cfg,
	}
}

func (receiver *writer) WriteErrorResponse(err error, writer http.ResponseWriter) {
	receiver.logger.Error("Writing error response", zap.Error(err))

	var errStruct any = map[string]string{"error": "Internal error"}
	var status = http.StatusInternalServerError

	var errorWithStatus appErrors.ErrorWithStatusCode
	if errors.As(err, &errorWithStatus) {
		status = errorWithStatus.StatusCode()
	}

	var serializableError appErrors.SerializableError
	var codeError appErrors.ErrorWithCode
	if errors.As(err, &serializableError) {
		errStruct = serializableError.SerializeToStruct()
	} else if errors.As(err, &codeError) {
		errStructMap, valid := errStruct.(map[string]string)
		code := codeError.GetCode()
		if valid && code != "" {
			errStructMap["code"] = code
		}
	}

	bytes := lo.Must(json.Marshal(errStruct))

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_, writeErr := writer.Write(bytes)
	if writeErr != nil {
		receiver.logger.Error("Error writing response", zap.Error(err))
	}
}

func (receiver *writer) WriteResponse(status int, body any, responseWriter http.ResponseWriter) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(status)

	var err error

	if strBody, ok := body.(string); ok {
		_, err = responseWriter.Write([]byte(strBody))
	} else if bytesBody, ok := body.([]byte); ok {
		_, err = responseWriter.Write(bytesBody)
	} else if body == nil {
		// do nothing
	} else {
		var bytes []byte
		bytes, err = json.Marshal(body)
		if err != nil {
			receiver.WriteErrorResponse(err, responseWriter)
			return
		}

		_, err = responseWriter.Write(bytes)
	}

	if err != nil {
		receiver.logger.Error("Error writing response", zap.Error(err))
	}
}
