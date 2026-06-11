package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"mysurl1/internal/model"
	types "mysurl1/internal/schema"
	"mysurl1/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type operationLogWriter interface {
	Insert(ctx context.Context, userID uint64, action, result, reason string) error
}

type operationLogProcess struct {
	Action    string
	OnSuccess func(r *http.Request, resp utils.Response) string
	OnFailure func(r *http.Request, resp utils.Response) string
}

var operationLogProcesses = map[string]map[string]operationLogProcess{
	http.MethodPost: {
		"/api/v1/auth/login": {
			Action: model.UserOperationActionLogin,
			OnSuccess: func(r *http.Request, resp utils.Response) string {
				return buildLoginSuccessReason(r)
			},
			OnFailure: func(r *http.Request, resp utils.Response) string {
				return resp.Msg
			},
		},
		"/api/v1/links": {
			Action: model.UserOperationActionCreateLink,
			OnSuccess: func(r *http.Request, resp utils.Response) string {
				reason, _ := resp.ExtData.(string)
				return reason
			},
			OnFailure: func(r *http.Request, resp utils.Response) string {
				return resp.Msg
			},
		},
		"/api/v1/links/batch": {
			Action: model.UserOperationActionCreateLinkBatch,
			OnSuccess: func(r *http.Request, resp utils.Response) string {
				successCount, failedCount := getBatchCreateCounts(resp)
				return fmt.Sprintf("success_count=%d failed_count=%d", successCount, failedCount)
			},
			OnFailure: func(r *http.Request, resp utils.Response) string {
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

		resp, ok := parseOperationLogResponse(writer.body.Bytes())
		if !ok {
			return
		}

		record := &utils.OperationLogPayload{
			UserID: getOperationUserID(r, resp),
			Action: process.Action,
		}

		if resp.Code == utils.CodeOK {
			record.Result = getOperationResult(process.Action, resp)
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

func getOperationUserID(r *http.Request, resp utils.Response) uint64 {
	if r == nil {
		return 0
	}

	claims, ok := utils.GetAuthClaims(r.Context())
	if !ok || claims == nil {
		return getOperationUserIDFromResponse(resp)
	}

	return claims.UserID
}

func getOperationUserIDFromResponse(resp utils.Response) uint64 {
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

	var data types.AuthResponse
	raw, err := json.Marshal(resp.Data)
	if err != nil {
		return 0
	}
	if err := json.Unmarshal(raw, &data); err != nil {
		return 0
	}

	return data.User.ID
}

func parseOperationLogResponse(body []byte) (utils.Response, bool) {
	var resp utils.Response
	if len(bytes.TrimSpace(body)) == 0 {
		return utils.Response{}, false
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return utils.Response{}, false
	}

	return resp, true
}

func buildLoginSuccessReason(r *http.Request) string {
	platform := strings.TrimSpace(r.Header.Get("User-Agent"))
	if platform == "" {
		platform = "unknown"
	}

	ip := getRequestIP(r)
	if ip == "" {
		ip = "unknown"
	}

	return fmt.Sprintf("platform=%s ip=%s", platform, ip)
}

func getRequestIP(r *http.Request) string {
	if r == nil {
		return ""
	}

	if forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); forwarded != "" {
		parts := strings.Split(forwarded, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}

	if realIP := strings.TrimSpace(r.Header.Get("X-Real-Ip")); realIP != "" {
		return realIP
	}

	hostPort := strings.TrimSpace(r.RemoteAddr)
	if hostPort == "" {
		return ""
	}

	host, _, err := strings.Cut(hostPort, ":")
	if err {
		return hostPort
	}

	return host
}

func getBatchCreateCounts(resp utils.Response) (int, int) {
	if resp.Data == nil {
		return 0, 0
	}

	if dataMap, ok := resp.Data.(map[string]any); ok {
		return getIntValue(dataMap["success_count"]), getIntValue(dataMap["failed_count"])
	}

	var data types.BatchCreateLinksResponse
	raw, err := json.Marshal(resp.Data)
	if err != nil {
		return 0, 0
	}
	if err := json.Unmarshal(raw, &data); err != nil {
		return 0, 0
	}

	return data.SuccessCount, data.FailedCount
}

func getOperationResult(action string, resp utils.Response) string {
	if action != model.UserOperationActionCreateLinkBatch {
		return model.UserOperationResultSuccess
	}

	successCount, failedCount := getBatchCreateCounts(resp)
	switch {
	case successCount > 0 && failedCount > 0:
		return model.UserOperationResultPartialSuccess
	case successCount > 0:
		return model.UserOperationResultSuccess
	default:
		return model.UserOperationResultFailed
	}
}

func getIntValue(value any) int {
	switch v := value.(type) {
	case int:
		return v
	case int32:
		return int(v)
	case int64:
		return int(v)
	case float64:
		return int(v)
	default:
		return 0
	}
}

type responseCaptureWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (w *responseCaptureWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseCaptureWriter) Write(data []byte) (int, error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}
