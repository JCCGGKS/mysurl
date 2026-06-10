package utils

import (
	"context"
	"testing"
)

func TestOperationLogPayloadLifecycle(t *testing.T) {
	ctx := WithOperationLogHolder(context.Background())
	targetCode := "abc123"

	SetOperationLogPayload(ctx, OperationLogPayload{
		UserID:     1001,
		Action:     "create_link",
		Result:     "success",
		TargetCode: &targetCode,
	})

	payload, ok := GetOperationLogPayload(ctx)
	if !ok {
		t.Fatalf("expected payload to be stored")
	}
	if payload.UserID != 1001 {
		t.Fatalf("unexpected user id: %d", payload.UserID)
	}
	if payload.Action != "create_link" {
		t.Fatalf("unexpected action: %s", payload.Action)
	}
	if payload.Result != "success" {
		t.Fatalf("unexpected result: %s", payload.Result)
	}
	if payload.TargetCode == nil || *payload.TargetCode != targetCode {
		t.Fatalf("unexpected target code: %+v", payload.TargetCode)
	}
}

func TestSetOperationLogPayloadIgnoreInvalidPayload(t *testing.T) {
	ctx := WithOperationLogHolder(context.Background())

	SetOperationLogPayload(ctx, OperationLogPayload{
		Action: "login",
		Result: "success",
	})

	if _, ok := GetOperationLogPayload(ctx); !ok {
		t.Fatalf("expected payload without user id to be stored")
	}
}

func TestSetOperationLogPayloadWithoutHolder(t *testing.T) {
	ctx := context.Background()
	SetOperationLogPayload(ctx, OperationLogPayload{
		UserID: 1001,
		Action: "login",
		Result: "success",
	})

	if _, ok := GetOperationLogPayload(ctx); ok {
		t.Fatalf("expected no payload without holder")
	}
}
