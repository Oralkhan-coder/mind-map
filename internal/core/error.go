package core

import "net/http"

type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func BadRequest(message string) *AppError {
	return New(http.StatusBadRequest, message)
}

func NotFound(message string) *AppError {
	return New(http.StatusNotFound, message)
}

func InternalServerError(message string) *AppError {
	return New(http.StatusInternalServerError, message)
}

func Forbidden(message string) *AppError {
	return New(http.StatusForbidden, message)
}

func Unauthorized(message string) *AppError {
	return New(http.StatusUnauthorized, message)
}
