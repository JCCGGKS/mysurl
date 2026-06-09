package model

import "time"

const (
	UserOperationActionLogin      = "login"
	UserOperationActionCreateLink = "create_link"
	UserOperationResultSuccess    = "success"
)

type UserOperationLog struct {
	ID         uint64    `db:"id"`
	UserID     uint64    `db:"user_id"`
	Action     string    `db:"action"`
	Result     string    `db:"result"`
	TargetID   *uint64   `db:"target_id"`
	TargetCode *string   `db:"target_code"`
	CreatedAt  time.Time `db:"created_at"`
}
