package apperror

import (
	"errors"
	"net/http"
)

func HTTPStatus(err error) int {
	var appErr *AppError
	if !errors.As(err, &appErr) {
		return http.StatusInternalServerError
	}

	switch appErr.Code {
	case CodeEmailAlreadyExists, CodeUsernameAlreadyExists:
		return http.StatusConflict
	case CodeValidation, CodeBadRequest:
		return http.StatusBadRequest
	case CodeNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
