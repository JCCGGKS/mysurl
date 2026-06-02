package utils

import (
	"errors"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/rest/httpx"
)

const (
	CodeOK = 0
	MsgOK  = "ok"
)

type Response struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	Data      any    `json:"data,omitempty"`
	Timestamp int64  `json:"timestamp"`
}

func Success(data any) Response {
	if IsNil(data) {
		data = nil
	}

	return Response{
		Code:      CodeOK,
		Msg:       MsgOK,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

func Error(code int, msg string) Response {
	return Response{
		Code:      code,
		Msg:       msg,
		Data:      nil,
		Timestamp: time.Now().Unix(),
	}
}

type HTTPError struct {
	StatusCode int
	Message    string
}

func (e *HTTPError) Error() string {
	return e.Message
}

func BadRequest(message string) *HTTPError {
	return &HTTPError{
		StatusCode: http.StatusBadRequest,
		Message:    message,
	}
}

func NotFound(message string) *HTTPError {
	return &HTTPError{
		StatusCode: http.StatusNotFound,
		Message:    message,
	}
}

func InternalError(message string) *HTTPError {
	return &HTTPError{
		StatusCode: http.StatusInternalServerError,
		Message:    message,
	}
}

func WriteJSONError(w http.ResponseWriter, r *http.Request, err error) {
	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		httpx.WriteJsonCtx(r.Context(), w, httpErr.StatusCode, Error(httpErr.StatusCode, httpErr.Message))
		return
	}

	httpx.WriteJsonCtx(r.Context(), w, http.StatusInternalServerError, Error(http.StatusInternalServerError, err.Error()))
}

func WriteJSONSuccess(w http.ResponseWriter, r *http.Request, data any) {
	httpx.OkJsonCtx(r.Context(), w, Success(data))
}

func WriteRedirectError(w http.ResponseWriter, err error) {
	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		http.Error(w, httpErr.Message, httpErr.StatusCode)
		return
	}

	http.Error(w, err.Error(), http.StatusInternalServerError)
}
