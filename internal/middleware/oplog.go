package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"mysurl1/internal/model"
	"mysurl1/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type operationLogWriter interface {
	Insert(ctx context.Context, userID uint64, action, result, reason string) error
}

type operationLogRecord struct {
	UserID uint64
	Action string
	Result string
	Reason string
}

type operationLogResponse struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Data    any    `json:"data"`
	ExtData any    `json:"-"`
}

type operationLogAuthResponseData struct {
	User struct {
		ID uint64 `json:"id"`
	} `json:"user"`
}

type operationLogProcess struct {
	Action    string
	OnSuccess func(r *http.Request, resp operationLogResponse) string
	OnFailure func(r *http.Request, resp operationLogResponse) string
}

var operationLogProcesses = map[string]map[string]operationLogProcess{
	http.MethodPost: {
		"/api/v1/auth/login": {
			Action: model.UserOperationActionLogin,
			OnSuccess: func(r *http.Request, resp operationLogResponse) string {
				return ""
			},
			OnFailure: func(r *http.Request, resp operationLogResponse) string {
				return resp.Msg
			},
		},
		"/api/v1/links": {
			Action: model.UserOperationActionCreateLink,
			OnSuccess: func(r *http.Request, resp operationLogResponse) string {
				reason, _ := resp.ExtData.(string)
				return reason
			},
			OnFailure: func(r *http.Request, resp operationLogResponse) string {
				return resp.Msg
			},
		},
		"/api/v1/links/batch": {
			Action: model.UserOperationActionCreateLinkBatch,
			OnSuccess: func(r *http.Request, resp operationLogResponse) string {
				return ""
			},
			OnFailure: func(r *http.Request, resp operationLogResponse) string {
				return resp.Msg
			},
		},
	},
	http.MethodGet:    {},
	http.MethodDelete: {},
}

type OperationLogMiddleware struct {
	dao operationLogWriter
}

func NewOperationLogMiddleware(dao operationLogWriter) *OperationLogMiddleware {
	return &OperationLogMiddleware{dao: dao}
}

func (m *OperationLogMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		process, ok := lookupOperationLogProcess(r.Method, r.URL.Path)
		if !ok {
			next(w, r)
			return
		}

		writer := &responseCaptureWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next(writer, r)

		if m == nil || m.dao == nil {
			return
		}

		resp, ok := getOperationLogResponse(writer)
		if !ok {
			resp, ok = parseOperationLogResponse(writer.body.Bytes())
			if !ok {
				return
			}
		}

		record := &operationLogRecord{
			UserID: getOperationUserID(r, resp),
			Action: process.Action,
		}

		if resp.Code == utils.CodeOK {
			record.Result = model.UserOperationResultSuccess
			if process.OnSuccess != nil {
				record.Reason = process.OnSuccess(r, resp)
			}
		} else {
			record.Result = model.UserOperationResultFailed
			if process.OnFailure != nil {
				record.Reason = process.OnFailure(r, resp)
			}
		}
		if record.Action == "" || record.Result == "" {
			return
		}

		if err := m.dao.Insert(r.Context(), record.UserID, record.Action, record.Result, record.Reason); err != nil {
			logx.Errorf("write user operation log failed: %v", err)
		}
	}
}

func lookupOperationLogProcess(method, path string) (operationLogProcess, bool) {
	routeMap, ok := operationLogProcesses[method]
	if !ok {
		return operationLogProcess{}, false
	}

	process, ok := routeMap[path]
	return process, ok
}

func getOperationUserID(r *http.Request, resp operationLogResponse) uint64 {
	if r == nil {
		return 0
	}

	claims, ok := utils.GetAuthClaims(r.Context())
	if !ok || claims == nil {
		return getOperationUserIDFromResponse(resp)
	}

	return claims.UserID
}

func getOperationUserIDFromResponse(resp operationLogResponse) uint64 {
	if resp.Data == nil {
		return 0
	}

	if dataMap, ok := resp.Data.(map[string]any); ok {
		userMap, ok := dataMap["user"].(map[string]any)
		if !ok {
			return 0
		}

		idValue, ok := userMap["id"]
		if !ok {
			return 0
		}

		switch id := idValue.(type) {
		case float64:
			return uint64(id)
		case uint64:
			return id
		}
	}

	var data operationLogAuthResponseData
	raw, err := json.Marshal(resp.Data)
	if err != nil {
		return 0
	}
	if err := json.Unmarshal(raw, &data); err != nil {
		return 0
	}

	return data.User.ID
}

func parseOperationLogResponse(body []byte) (operationLogResponse, bool) {
	var resp operationLogResponse
	if len(bytes.TrimSpace(body)) == 0 {
		return resp, false
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return resp, false
	}

	return resp, true
}

func getOperationLogResponse(w http.ResponseWriter) (operationLogResponse, bool) {
	resp, ok := utils.GetOperationLogResponse(w)
	if !ok || resp == nil {
		return operationLogResponse{}, false
	}

	return operationLogResponse{
		Code:    resp.Code,
		Msg:     resp.Msg,
		Data:    resp.Data,
		ExtData: resp.ExtData,
	}, true
}

type responseCaptureWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
	resp       *utils.Response
}

func (w *responseCaptureWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseCaptureWriter) Write(data []byte) (int, error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

func (w *responseCaptureWriter) SetOperationLogResponse(resp utils.Response) {
	respCopy := resp
	w.resp = &respCopy
}

func (w *responseCaptureWriter) GetOperationLogResponse() *utils.Response {
	return w.resp
}
