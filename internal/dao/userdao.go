package dao

import (
	"context"

	"mysurl1/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type UserDAO struct {
	conn sqlx.SqlConn
}

func NewUserDAO(conn sqlx.SqlConn) *UserDAO {
	return &UserDAO{conn: conn}
}

func (d *UserDAO) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	query := `
SELECT
	id,
	username,
	password_hash
FROM users
WHERE username = ?
  AND deleted_at IS NULL
LIMIT 1
`

	if err := d.conn.QueryRowPartialCtx(ctx, &user, query, username); err != nil {
		return nil, err
	}

	return &user, nil
}

func (d *UserDAO) FindByID(ctx context.Context, id uint64) (*model.User, error) {
	var user model.User
	query := `
SELECT
	id,
	username,
	password_hash
FROM users
WHERE id = ?
  AND deleted_at IS NULL
LIMIT 1
`

	if err := d.conn.QueryRowPartialCtx(ctx, &user, query, id); err != nil {
		return nil, err
	}

	return &user, nil
}

func (d *UserDAO) Insert(ctx context.Context, username, passwordHash string) (uint64, error) {
	query := `
INSERT INTO users (
	username,
	password_hash
) VALUES (?, ?)
`

	result, err := d.conn.ExecCtx(ctx, query, username, passwordHash)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(id), nil
}

func (d *UserDAO) UpdatePassword(ctx context.Context, id uint64, passwordHash string) error {
	query := `
UPDATE users
SET password_hash = ?
WHERE id = ?
  AND deleted_at IS NULL
`

	_, err := d.conn.ExecCtx(ctx, query, passwordHash, id)
	return err
}
