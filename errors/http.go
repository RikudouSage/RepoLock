package errors

import "net/http"

func NewAccessDeniedError() *UserFacingError {
	return NewUserFacingErrorWithStatusCode(
		"you do not have permissions to perform this action",
		http.StatusForbidden,
	)
}
