package errors

import "net/http"

type UserFacingError struct {
	Message string
	Code    string
	Status  int
}

func NewUserFacingError(message string) *UserFacingError {
	return &UserFacingError{Message: message, Status: http.StatusBadRequest}
}

func NewUserFacingErrorWithCode(message, code string) *UserFacingError {
	return &UserFacingError{Message: message, Code: code, Status: http.StatusBadRequest}
}

func NewUserFacingErrorWithStatusCode(message string, status int) *UserFacingError {
	return &UserFacingError{Message: message, Status: status}
}

func NewUserFacingErrorWithCodeAndStatusCode(message, code string, status int) *UserFacingError {
	return &UserFacingError{Message: message, Code: code, Status: status}
}

func (receiver *UserFacingError) StatusCode() int {
	return receiver.Status
}

func (receiver *UserFacingError) SerializeToStruct() any {
	result := map[string]string{
		"error": receiver.Message,
	}
	if receiver.Code != "" {
		result["code"] = receiver.Code
	}

	return result
}

func (receiver *UserFacingError) Error() string {
	return receiver.Message
}

func (receiver *UserFacingError) GetCode() string {
	return receiver.Code
}
