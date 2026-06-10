package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"mysurl1/internal/utils"
)

type stubOperationLogWriter struct {
	called     bool
	userID     uint64
	action     string
	result     string
	reason     string
	targetCode *string
	err        error
}

func (s *stubOperationLogWriter) Insert(_ context.Context, userID uint64, action, result, reason string, targetCode *string) error {
	s.called = true
	s.userID = userID
	s.action = action
	s.result = result
	s.reason = reason
	s.targetCode = targetCode
	return s.err
}

func TestOperationLogMiddlewareWriteOnSuccess(t *testing.T) {
	writer := &stubOperationLogWriter{}
	m := NewOperationLogMiddleware(writer)

	handler := m.Handle(func(w http.ResponseWriter, r *http.Request) {
		targetCode := "code9"
		utils.SetOperationLogPayload(r.Context(), utils.OperationLogPayload{
			UserID:     101,
			Action:     "create_link",
			Result:     "success",
			TargetCode: &targetCode,
		})
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/links", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)

	if !writer.called {
		t.Fatalf("expected insert to be called")
	}
	if writer.userID != 101 || writer.action != "create_link" || writer.result != "success" {
		t.Fatalf("unexpected payload persisted: %+v", writer)
	}
	if writer.targetCode == nil || *writer.targetCode != "code9" {
		t.Fatalf("unexpected target code: %+v", writer.targetCode)
	}
}

func TestOperationLogMiddlewareWriteOnFailureStatus(t *testing.T) {
	writer := &stubOperationLogWriter{}
	m := NewOperationLogMiddleware(writer)

	handler := m.Handle(func(w http.ResponseWriter, r *http.Request) {
		utils.SetOperationLogPayload(r.Context(), utils.OperationLogPayload{
			UserID: 101,
			Action: "login",
			Result: "failed",
			Reason: "username or password is invalid",
		})
		w.WriteHeader(http.StatusUnauthorized)
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)

	if !writer.called {
		t.Fatalf("expected insert on failure status when payload exists")
	}
	if writer.result != "failed" || writer.reason != "username or password is invalid" {
		t.Fatalf("unexpected failure payload persisted: %+v", writer)
	}
}

func TestOperationLogMiddlewareSkipWithoutPayload(t *testing.T) {
	writer := &stubOperationLogWriter{}
	m := NewOperationLogMiddleware(writer)

	handler := m.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/user-operation-logs", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)

	if writer.called {
		t.Fatalf("expected insert to be skipped without payload")
	}
}

func TestOperationLogMiddlewareIgnoreInsertError(t *testing.T) {
	writer := &stubOperationLogWriter{err: errors.New("insert failed")}
	m := NewOperationLogMiddleware(writer)

	handler := m.Handle(func(w http.ResponseWriter, r *http.Request) {
		utils.SetOperationLogPayload(r.Context(), utils.OperationLogPayload{
			UserID: 101,
			Action: "login",
			Result: "success",
		})
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	if !writer.called {
		t.Fatalf("expected insert attempt even when writer returns error")
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status code: %d", rec.Code)
	}
}
