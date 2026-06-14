package dao

import (
	"context"
	"errors"
	"time"

	"mysurl1/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type UserRefreshTokenDAO struct {
	conn sqlx.SqlConn
}

func NewUserRefreshTokenDAO(conn sqlx.SqlConn) *UserRefreshTokenDAO {
	return &UserRefreshTokenDAO{conn: conn}
}

func (d *UserRefreshTokenDAO) Insert(ctx context.Context, userID uint64, tokenHash string, expiresAt time.Time) error {
	query := `
INSERT INTO user_refresh_tokens (
	user_id,
	token_hash,
	expires_at
) VALUES (?, ?, ?)
`

	_, err := d.conn.ExecCtx(ctx, query, userID, tokenHash, expiresAt)
	return err
}

func (d *UserRefreshTokenDAO) FindActiveByTokenHash(ctx context.Context, tokenHash string) (*model.UserRefreshToken, error) {
	var token model.UserRefreshToken
	query := `
SELECT
	id,
	user_id,
	token_hash,
	expires_at,
	revoked_at,
	created_at,
	updated_at
FROM user_refresh_tokens
WHERE token_hash = ?
  AND revoked_at IS NULL
  AND expires_at > NOW()
LIMIT 1
`

	if err := d.conn.QueryRowPartialCtx(ctx, &token, query, tokenHash); err != nil {
		return nil, err
	}

	return &token, nil
}

func (d *UserRefreshTokenDAO) RevokeByTokenHash(ctx context.Context, tokenHash string) error {
	query := `
UPDATE user_refresh_tokens
SET revoked_at = NOW(),
	updated_at = NOW()
WHERE token_hash = ?
  AND revoked_at IS NULL
`

	_, err := d.conn.ExecCtx(ctx, query, tokenHash)
	return err
}

func (d *UserRefreshTokenDAO) Rotate(ctx context.Context, oldTokenHash string, userID uint64, newTokenHash string, newExpiresAt time.Time) error {
	return d.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		revokeQuery := `
UPDATE user_refresh_tokens
SET revoked_at = NOW(),
	updated_at = NOW()
WHERE token_hash = ?
  AND revoked_at IS NULL
`

		result, err := session.ExecCtx(ctx, revokeQuery, oldTokenHash)
		if err != nil {
			return err
		}

		rows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rows == 0 {
			return sqlx.ErrNotFound
		}

		insertQuery := `
INSERT INTO user_refresh_tokens (
	user_id,
	token_hash,
	expires_at
) VALUES (?, ?, ?)
`

		_, err = session.ExecCtx(ctx, insertQuery, userID, newTokenHash, newExpiresAt)
		return err
	})
}

func (d *UserRefreshTokenDAO) RevokeIfExists(ctx context.Context, tokenHash string) error {
	if tokenHash == "" {
		return nil
	}

	err := d.RevokeByTokenHash(ctx, tokenHash)
	if err != nil && errors.Is(err, sqlx.ErrNotFound) {
		return nil
	}

	return err
}
