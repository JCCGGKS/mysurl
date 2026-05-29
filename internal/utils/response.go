package utils

import "time"

const (
	CodeOK    = 0
	MsgOK     = "ok"
	CodeError = 1
	MsgError  = "error"
)

type Response struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	Data      any    `json:"data,omitempty"`
	Timestamp int64  `json:"timestamp"`
}

func NewResponse(code int, msg string, data any) Response {
	return Response{
		Code:      code,
		Msg:       msg,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

func Success(data any) Response {
	return NewResponse(CodeOK, MsgOK, data)
}

func Error(code int, msg string) Response {
	return NewResponse(code, msg, nil)
}
