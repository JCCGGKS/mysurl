package middleware

import (
	"context"
	"net/http"

	"mysurl1/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type operationLogWriter interface {
	Insert(ctx context.Context, userID uint64, action, result, reason string, targetCode *string) error
}

type OperationLogMiddleware struct {
	dao operationLogWriter
}

func NewOperationLogMiddleware(dao operationLogWriter) *OperationLogMiddleware {
	return &OperationLogMiddleware{dao: dao}
}

func (m *OperationLogMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := utils.WithOperationLogHolder(r.Context())
		writer := &statusWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next(writer, r.WithContext(ctx))

		payload, ok := utils.GetOperationLogPayload(ctx)
		if !ok {
			return
		}
		if m == nil || m.dao == nil {
			return
		}

		if err := m.dao.Insert(ctx, payload.UserID, payload.Action, payload.Result, payload.Reason, payload.TargetCode); err != nil {
			logx.Errorf("write user operation log failed: %v", err)
		}
	}
}

type statusWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *statusWriter) Write(data []byte) (int, error) {
	return w.ResponseWriter.Write(data)
}
