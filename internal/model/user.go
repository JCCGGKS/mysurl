package model

import "time"

type User struct {
	ID           uint64     `db:"id"`
	Username     string     `db:"username"`
	PasswordHash string     `db:"password_hash"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
}
