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
	ExtData   any    `json:"extdata,omitempty"`
	Timestamp int64  `json:"timestamp"`
}

type ResponseOption func(*Response)

func WithResponseExtData(extData any) ResponseOption {
	return func(resp *Response) {
		resp.ExtData = extData
	}
}

func Success(data any, opts ...ResponseOption) Response {
	if IsNil(data) {
		data = nil
	}

	resp := Response{
		Code:      CodeOK,
		Msg:       MsgOK,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}

	for _, opt := range opts {
		if opt != nil {
			opt(&resp)
		}
	}

	return resp
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

func Unauthorized(message string) *HTTPError {
	return &HTTPError{
		StatusCode: http.StatusUnauthorized,
		Message:    message,
	}
}

func Conflict(message string) *HTTPError {
	return &HTTPError{
		StatusCode: http.StatusConflict,
		Message:    message,
	}
}

func WriteJSONError(w http.ResponseWriter, r *http.Request, err error) {
	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		writeJSONResponse(w, r, httpErr.StatusCode, Error(httpErr.StatusCode, httpErr.Message))
		return
	}

	writeJSONResponse(w, r, http.StatusInternalServerError, Error(http.StatusInternalServerError, err.Error()))
}

func WriteJSONSuccess(w http.ResponseWriter, r *http.Request, data any, opts ...ResponseOption) {
	writeJSONResponse(w, r, http.StatusOK, Success(data, opts...))
}

func WriteRedirectError(w http.ResponseWriter, err error) {
	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		http.Error(w, httpErr.Message, httpErr.StatusCode)
		return
	}

	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func writeJSONResponse(w http.ResponseWriter, r *http.Request, statusCode int, resp Response) {
	httpx.WriteJsonCtx(r.Context(), w, statusCode, resp)
}
