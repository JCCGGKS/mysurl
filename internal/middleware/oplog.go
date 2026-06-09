package middleware

import (
	"net/http"

	"mysurl1/internal/dao"
	"mysurl1/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type OperationLogMiddleware struct {
	dao *dao.UserOperationLogDAO
}

func NewOperationLogMiddleware(dao *dao.UserOperationLogDAO) *OperationLogMiddleware {
	return &OperationLogMiddleware{dao: dao}
}

func (m *OperationLogMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := utils.WithOperationLogHolder(r.Context())
		writer := &statusWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next(writer, r.WithContext(ctx))

		if writer.statusCode >= http.StatusBadRequest {
			return
		}

		payload, ok := utils.GetOperationLogPayload(ctx)
		if !ok {
			return
		}
		if m == nil || m.dao == nil {
			return
		}

		if err := m.dao.Insert(ctx, payload.UserID, payload.Action, payload.Result, payload.TargetID, payload.TargetCode); err != nil {
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
