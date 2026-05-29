package utils

import "net/http"

type HTTPError struct {
	StatusCode int
	Message    string
}

func (e *HTTPError) Error() string {
	return e.Message
}

func NewHTTPError(statusCode int, message string) *HTTPError {
	return &HTTPError{
		StatusCode: statusCode,
		Message:    message,
	}
}

func BadRequest(message string) *HTTPError {
	return NewHTTPError(http.StatusBadRequest, message)
}

func NotFound(message string) *HTTPError {
	return NewHTTPError(http.StatusNotFound, message)
}

func InternalError(message string) *HTTPError {
	return NewHTTPError(http.StatusInternalServerError, message)
}
