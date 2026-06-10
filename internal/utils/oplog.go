package utils

import "context"

type operationLogContextKey struct{}

type OperationLogPayload struct {
	UserID     uint64
	Action     string
	Result     string
	Reason     string
	TargetCode *string
}

type operationLogHolder struct {
	payload *OperationLogPayload
}

func WithOperationLogHolder(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return context.WithValue(ctx, operationLogContextKey{}, &operationLogHolder{})
}

func SetOperationLogPayload(ctx context.Context, payload OperationLogPayload) {
	if ctx == nil || payload.Action == "" || payload.Result == "" {
		return
	}

	holder, ok := ctx.Value(operationLogContextKey{}).(*operationLogHolder)
	if !ok || holder == nil {
		return
	}

	copyPayload := payload
	holder.payload = &copyPayload
}

func GetOperationLogPayload(ctx context.Context) (*OperationLogPayload, bool) {
	if ctx == nil {
		return nil, false
	}

	holder, ok := ctx.Value(operationLogContextKey{}).(*operationLogHolder)
	if !ok || holder == nil || holder.payload == nil {
		return nil, false
	}

	return holder.payload, true
}
