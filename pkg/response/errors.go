package response

import (
	"errors"
	"fmt"
	"net/http"
)

type Type string

const (
	Authorization        Type = "AUTHORIZATION"
	BadRequest           Type = "BAD_REQUEST"
	Conflict             Type = "CONFLICT"
	Internal             Type = "INTERNAL"
	NotFound             Type = "NOT_FOUND"
	PayloadTooLarge      Type = "PAYLOAD_TOO_LARGE"
	ServiceUnavailable   Type = "SERVICE_UNAVAILABLE"
	UnsupportedMediaType Type = "UNSUPPORTED_MEDIA_TYPE"
)

type Error struct {
	Type    Type   `json:"type"`
	Value   string `json:"value,omitempty"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Status() int {
	switch e.Type {
	case Authorization:
		return http.StatusUnauthorized
	case BadRequest:
		return http.StatusBadRequest
	case Conflict:
		return http.StatusConflict
	case Internal:
		return http.StatusInternalServerError
	case NotFound:
		return http.StatusNotFound
	case PayloadTooLarge:
		return http.StatusRequestEntityTooLarge
	case ServiceUnavailable:
		return http.StatusServiceUnavailable
	case UnsupportedMediaType:
		return http.StatusUnsupportedMediaType
	default:
		return http.StatusInternalServerError
	}
}

func Status(err error) int {
	var e *Error
	if errors.As(err, &e) {
		return e.Status()
	}
	return http.StatusInternalServerError
}

func NewAuthorizationError(reason string) *Error {
	return &Error{
		Type:    Authorization,
		Message: reason,
	}
}

func NewBadRequestError(reason string) *Error {
	return &Error{
		Type:    BadRequest,
		Message: fmt.Sprintf("Bad request. Reason: %v", reason),
	}
}

func NewConflictError(name string, value string) *Error {
	return &Error{
		Type:    Conflict,
		Value:   value,
		Message: fmt.Sprintf("resource: %v with value: %v already exists", name, value),
	}
}

func NewInternalServerError() *Error {
	return &Error{
		Type:    Internal,
		Message: "Internal server error.",
	}
}

func NewNotFoundError(name string, value string) *Error {
	return &Error{
		Type:    NotFound,
		Value:   value,
		Message: fmt.Sprintf("resource: %v with value: %v not found", name, value),
	}
}

func NewPayloadTooLargeError(maxBodySize int64, contentLength int64) *Error {
	return &Error{
		Type:    PayloadTooLarge,
		Message: fmt.Sprintf("Max payload size of %v exceeded. Actual payload size: %v", maxBodySize, contentLength),
	}
}

func NewServiceUnavailableError() *Error {
	return &Error{
		Type:    ServiceUnavailable,
		Message: "Service unavailable or timed out",
	}
}

func NewUnsupportedMediaTypeError(reason string) *Error {
	return &Error{
		Type:    UnsupportedMediaType,
		Message: reason,
	}
}
