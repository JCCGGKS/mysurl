package utils

type OperationLogPayload struct {
	UserID uint64
	Action string
	Result string
	Reason string
}
