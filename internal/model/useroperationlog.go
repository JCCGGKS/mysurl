package model

import "time"

const (
	UserOperationActionLogin           = "login"
	UserOperationActionCreateLink      = "create_link"
	UserOperationActionCreateLinkBatch = "create_link_batch"
	UserOperationResultSuccess         = "success"
	UserOperationResultFailed          = "failed"
	UserOperationResultPartialSuccess  = "partial_success"
)

type UserOperationLog struct {
	ID        uint64    `db:"id"`
	UserID    uint64    `db:"user_id"`
	Action    string    `db:"action"`
	Result    string    `db:"result"`
	Reason    *string   `db:"reason"`
	CreatedAt time.Time `db:"created_at"`
}
