package utils

import "net/http"

type OperationLogPayload struct {
	UserID uint64
	Action string
	Result string
	Reason string
}

type OperationLogResponseWriter interface {
	SetOperationLogResponse(resp Response)
}

type operationLogResponseReader interface {
	GetOperationLogResponse() *Response
}

func SetOperationLogResponse(w http.ResponseWriter, resp Response) {
	if w == nil {
		return
	}

	writer, ok := w.(OperationLogResponseWriter)
	if !ok {
		return
	}

	writer.SetOperationLogResponse(resp)
}

func GetOperationLogResponse(w http.ResponseWriter) (*Response, bool) {
	if w == nil {
		return nil, false
	}

	writer, ok := w.(operationLogResponseReader)
	if !ok {
		return nil, false
	}

	resp := writer.GetOperationLogResponse()
	return resp, resp != nil
}
