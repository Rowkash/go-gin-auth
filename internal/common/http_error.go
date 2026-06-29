package common

import (
	"net/http"

	"github.com/go-playground/validator/v10"
)

type AppError struct {
	StatusCode int    `json:"statusCode"`
	ErrorType  string `json:"error"`
	Details    string `json:"details"`
}

func (e *AppError) Error() string {
	return e.Details
}

type ValidationErrorResponse struct {
	StatusCode int               `json:"statusCode"`
	ErrorType  string            `json:"error"`
	Details    []ValidationError `json:"details"`
}

func NewNotFoundError(details string) *AppError {
	return &AppError{
		StatusCode: http.StatusNotFound,
		ErrorType:  "NotFound",
		Details:    details,
	}
}

func NewBadRequestError(details string) *AppError {
	return &AppError{
		StatusCode: http.StatusBadRequest,
		ErrorType:  "BadRequest",
		Details:    details,
	}
}

func NewUnauthorizedError(details string) *AppError {
	return &AppError{
		StatusCode: http.StatusUnauthorized,
		ErrorType:  "Unauthorized",
		Details:    details,
	}
}

func NewForbiddenError(details string) *AppError {
	return &AppError{
		StatusCode: http.StatusForbidden,
		ErrorType:  "Forbidden",
		Details:    details,
	}
}

func NewInternalServerError() *AppError {
	return &AppError{
		StatusCode: http.StatusInternalServerError,
		ErrorType:  "InternalServerError",
		Details:    "An unexpected error occurred on the server.",
	}
}

func NewValidationError(err validator.ValidationErrors) ValidationErrorResponse {
	return ValidationErrorResponse{
		StatusCode: http.StatusBadRequest,
		ErrorType:  "BadRequest",
		Details:    FormatValidationErrors(err),
	}
}
