package errors

import (
	"fmt"
	"net/http"
)

type ErrorCode string

const (
	ErrInternal             ErrorCode = "INTERNAL_ERROR"
	ErrInvalidInput         ErrorCode = "INVALID_INPUT"
	ErrNotFound             ErrorCode = "NOT_FOUND"
	ErrUnauthorized         ErrorCode = "UNAUTHORIZED"
	ErrForbidden            ErrorCode = "FORBIDDEN"
	ErrConflict             ErrorCode = "CONFLICT"
	ErrIdempotency          ErrorCode = "IDEMPOTENCY_CONFLICT"
	ErrInventoryUnavailable ErrorCode = "INVENTORY_UNAVAILABLE"
	ErrPaymentDeclined      ErrorCode = "PAYMENT_DECLINED"
	ErrPaymentTimeout       ErrorCode = "PAYMENT_TIMEOUT"
	ErrTenantSuspended      ErrorCode = "TENANT_SUSPENDED"
	ErrRateLimit            ErrorCode = "RATE_LIMIT_EXCEEDED"
)

type AppError struct {
	Code    ErrorCode      `json:"code"`
	Message string         `json:"message"`
	TraceID string         `json:"trace_id,omitempty"`
	Details map[string]any `json:"details,omitempty"`
	Status  int            `json:"-"`
	Cause   error          `json:"-"`
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (cause: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

func New(code ErrorCode, message string) *AppError {
	return &AppError{Code: code, Message: message, Details: make(map[string]any)}
}

func Wrap(code ErrorCode, message string, cause error) *AppError {
	return &AppError{Code: code, Message: message, Cause: cause, Details: make(map[string]any)}
}

func (e *AppError) WithTraceID(traceID string) *AppError {
	e.TraceID = traceID
	return e
}

func (e *AppError) WithDetail(key string, value any) *AppError {
	e.Details[key] = value
	return e
}

func HTTPStatus(err error) int {
	if err == nil {
		return http.StatusOK
	}
	if ae, ok := err.(*AppError); ok {
		if ae.Status != 0 {
			return ae.Status
		}
		switch ae.Code {
		case ErrInvalidInput:
			return http.StatusBadRequest
		case ErrNotFound:
			return http.StatusNotFound
		case ErrUnauthorized:
			return http.StatusUnauthorized
		case ErrForbidden:
			return http.StatusForbidden
		case ErrConflict, ErrIdempotency:
			return http.StatusConflict
		case ErrRateLimit:
			return http.StatusTooManyRequests
		default:
			return http.StatusInternalServerError
		}
	}
	return http.StatusInternalServerError
}

func IsNotFound(err error) bool {
	if ae, ok := err.(*AppError); ok {
		return ae.Code == ErrNotFound
	}
	return false
}

func IsConflict(err error) bool {
	if ae, ok := err.(*AppError); ok {
		return ae.Code == ErrConflict || ae.Code == ErrIdempotency
	}
	return false
}
