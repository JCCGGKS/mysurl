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
	called bool
	userID uint64
	action string
	result string
	reason string
	err    error
}

func (s *stubOperationLogWriter) Insert(_ context.Context, userID uint64, action, result, reason string) error {
	s.called = true
	s.userID = userID
	s.action = action
	s.result = result
	s.reason = reason
	return s.err
}

func TestOperationLogMiddlewareWriteOnSuccess(t *testing.T) {
	writer := &stubOperationLogWriter{}
	m := NewOperationLogMiddleware(writer)

	handler := m.Handle(func(w http.ResponseWriter, r *http.Request) {
		utils.WriteJSONSuccessWithExtData(w, r, map[string]any{
			"short_code": "code9",
		}, "cache_hit")
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/links", nil)
	req = req.WithContext(utils.WithAuthClaims(req.Context(), &utils.AuthClaims{UserID: 101}))
	rec := httptest.NewRecorder()
	handler(rec, req)

	if !writer.called {
		t.Fatalf("expected insert to be called")
	}
	if writer.userID != 101 || writer.action != "create_link" || writer.result != "success" {
		t.Fatalf("unexpected payload persisted: %+v", writer)
	}
	if writer.reason != "cache_hit" {
		t.Fatalf("unexpected success reason: %q", writer.reason)
	}
}

func TestOperationLogMiddlewareWriteOnFailure(t *testing.T) {
	writer := &stubOperationLogWriter{}
	m := NewOperationLogMiddleware(writer)

	handler := m.Handle(func(w http.ResponseWriter, r *http.Request) {
		utils.WriteJSONError(w, r, utils.Unauthorized("username or password is invalid"))
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)

	if !writer.called {
		t.Fatalf("expected insert on failure response")
	}
	if writer.action != "login" || writer.result != "failed" {
		t.Fatalf("unexpected failure payload persisted: %+v", writer)
	}
	if writer.reason != "username or password is invalid" {
		t.Fatalf("unexpected failure reason: %q", writer.reason)
	}
}

func TestOperationLogMiddlewareWriteLoginUserIDFromResponse(t *testing.T) {
	writer := &stubOperationLogWriter{}
	m := NewOperationLogMiddleware(writer)

	handler := m.Handle(func(w http.ResponseWriter, r *http.Request) {
		utils.WriteJSONSuccess(w, r, map[string]any{
			"token":      "jwt-token",
			"expires_at": 1234567890,
			"user": map[string]any{
				"id":       uint64(88),
				"username": "tester",
			},
		})
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)

	if !writer.called {
		t.Fatalf("expected insert to be called")
	}
	if writer.userID != 88 {
		t.Fatalf("unexpected user id persisted: %d", writer.userID)
	}
	if writer.action != "login" || writer.result != "success" {
		t.Fatalf("unexpected payload persisted: %+v", writer)
	}
}

func TestOperationLogMiddlewareSkipUnknownRoute(t *testing.T) {
	writer := &stubOperationLogWriter{}
	m := NewOperationLogMiddleware(writer)

	handler := m.Handle(func(w http.ResponseWriter, r *http.Request) {
		utils.WriteJSONSuccess(w, r, map[string]any{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/user-operation-logs", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)

	if writer.called {
		t.Fatalf("expected insert to be skipped without process")
	}
}

func TestOperationLogMiddlewareIgnoreInsertError(t *testing.T) {
	writer := &stubOperationLogWriter{err: errors.New("insert failed")}
	m := NewOperationLogMiddleware(writer)

	handler := m.Handle(func(w http.ResponseWriter, r *http.Request) {
		utils.WriteJSONSuccess(w, r, map[string]any{
			"short_code": "code9",
		})
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/links", nil)
	req = req.WithContext(utils.WithAuthClaims(req.Context(), &utils.AuthClaims{UserID: 101}))
	rec := httptest.NewRecorder()
	handler(rec, req)

	if !writer.called {
		t.Fatalf("expected insert attempt even when writer returns error")
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status code: %d", rec.Code)
	}
}
