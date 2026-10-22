package model

import "fmt"

// ErrorDetail memuat detail standar error sesuai kontrak API
type ErrorDetail struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

// ErrorResponse adalah wrapper standar format respons HTTP error
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// AppError merepresentasikan domain error dengan HTTP status dan kode error
type AppError struct {
	StatusCode int
	Code       string
	Message    string
	Internal   error
}

func (e *AppError) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Internal)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Internal
}

func NewBadRequest(message string, err ...error) *AppError {
	var internal error
	if len(err) > 0 {
		internal = err[0]
	}
	return &AppError{
		StatusCode: 400,
		Code:       "BAD_REQUEST",
		Message:    message,
		Internal:   internal,
	}
}

func NewNotFound(message string) *AppError {
	return &AppError{
		StatusCode: 404,
		Code:       "NOT_FOUND",
		Message:    message,
	}
}

func NewUnauthorized(message string) *AppError {
	return &AppError{
		StatusCode: 401,
		Code:       "UNAUTHORIZED",
		Message:    message,
	}
}

func NewForbidden(message string) *AppError {
	return &AppError{
		StatusCode: 403,
		Code:       "FORBIDDEN",
		Message:    message,
	}
}

func NewConflict(message string) *AppError {
	return &AppError{
		StatusCode: 409,
		Code:       "CONFLICT",
		Message:    message,
	}
}

func NewUnprocessableEntity(message string) *AppError {
	return &AppError{
		StatusCode: 422,
		Code:       "UNPROCESSABLE_ENTITY",
		Message:    message,
	}
}

func NewInternalError(message string, err error) *AppError {
	return &AppError{
		StatusCode: 500,
		Code:       "INTERNAL_SERVER_ERROR",
		Message:    message,
		Internal:   err,
	}
}
