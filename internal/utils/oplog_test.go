package utils

import "testing"

func TestOperationLogPayloadStoresCoreFields(t *testing.T) {
	payload := OperationLogPayload{
		UserID: 1001,
		Action: "create_link",
		Result: "success",
		Reason: "",
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
}
