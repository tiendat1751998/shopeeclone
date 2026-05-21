package errors

import (
	"fmt"
	"net/http"
)

type ErrorCode string

const (
	ErrInternal         ErrorCode = "INTERNAL_ERROR"
	ErrNotFound         ErrorCode = "NOT_FOUND"
	ErrBadRequest       ErrorCode = "BAD_REQUEST"
	ErrUnauthorized     ErrorCode = "UNAUTHORIZED"
	ErrForbidden        ErrorCode = "FORBIDDEN"
	ErrConflict         ErrorCode = "CONFLICT"
	ErrValidation       ErrorCode = "VALIDATION_ERROR"
	ErrRateLimited      ErrorCode = "RATE_LIMITED"
	ErrServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
	ErrStockInsufficient ErrorCode = "STOCK_INSUFFICIENT"
	ErrDuplicate        ErrorCode = "DUPLICATE_ENTRY"
)

type AppError struct {
	Code       ErrorCode    `json:"error_code"`
	Message    string       `json:"message"`
	HTTPStatus int          `json:"-"`
	Details    []ErrorDetail `json:"details,omitempty"`
	Err        error        `json:"-"`
}

type ErrorDetail struct {
	Field string `json:"field,omitempty"`
	Issue string `json:"issue"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) WithDetail(field, issue string) *AppError {
	e.Details = append(e.Details, ErrorDetail{Field: field, Issue: issue})
	return e
}

func (e *AppError) WithDetails(details []ErrorDetail) *AppError {
	e.Details = append(e.Details, details...)
	return e
}

func New(code ErrorCode, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

func NewInternal(message string) *AppError {
	return New(ErrInternal, message, http.StatusInternalServerError)
}

func NewInternalError(err error) *AppError {
	return &AppError{
		Code:       ErrInternal,
		Message:    "An internal error occurred",
		HTTPStatus: http.StatusInternalServerError,
		Err:        err,
	}
}

func NewNotFound(message string) *AppError {
	return New(ErrNotFound, message, http.StatusNotFound)
}

func NewBadRequest(message string) *AppError {
	return New(ErrBadRequest, message, http.StatusBadRequest)
}

func NewUnauthorized(message string) *AppError {
	return New(ErrUnauthorized, message, http.StatusUnauthorized)
}

func NewForbidden(message string) *AppError {
	return New(ErrForbidden, message, http.StatusForbidden)
}

func NewConflict(message string) *AppError {
	return New(ErrConflict, message, http.StatusConflict)
}

func NewValidation(message string) *AppError {
	return New(ErrValidation, message, http.StatusUnprocessableEntity)
}

func NewRateLimited(message string) *AppError {
	return New(ErrRateLimited, message, http.StatusTooManyRequests)
}

func NewServiceUnavailable(message string) *AppError {
	return New(ErrServiceUnavailable, message, http.StatusServiceUnavailable)
}

func NewStockInsufficient(message string) *AppError {
	return New(ErrStockInsufficient, message, http.StatusConflict)
}

func NewDuplicate(message string) *AppError {
	return New(ErrDuplicate, message, http.StatusConflict)
}

func FromError(err error) *AppError {
	if err == nil {
		return nil
	}
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return NewInternalError(err)
}

func HTTPStatus(err error) int {
	if appErr, ok := err.(*AppError); ok {
		return appErr.HTTPStatus
	}
	return http.StatusInternalServerError
}

func ErrorCodeFromError(err error) ErrorCode {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code
	}
	return ErrInternal
}
