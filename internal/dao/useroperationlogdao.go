package dao

import (
	"context"
	"errors"
	"strings"

	"mysurl1/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type UserOperationLogDAO struct {
	conn sqlx.SqlConn
}

func NewUserOperationLogDAO(conn sqlx.SqlConn) *UserOperationLogDAO {
	return &UserOperationLogDAO{conn: conn}
}

func (d *UserOperationLogDAO) Insert(ctx context.Context, userID uint64, action, result string, targetID *uint64, targetCode *string) error {
	query := `
INSERT INTO user_operation_logs (
	user_id,
	action,
	result,
	target_id,
	target_code
) VALUES (?, ?, ?, ?, ?)
`

	_, err := d.conn.ExecCtx(ctx, query, userID, action, result, targetID, targetCode)
	return err
}

func (d *UserOperationLogDAO) CountByUserID(ctx context.Context, userID uint64) (int64, error) {
	var result struct {
		Total int64 `db:"total"`
	}

	query := `
SELECT COUNT(1) AS total
FROM user_operation_logs
WHERE user_id = ?
`

	if err := d.conn.QueryRowPartialCtx(ctx, &result, query, userID); err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return 0, nil
		}

		return 0, err
	}

	return result.Total, nil
}

func (d *UserOperationLogDAO) ListByUserIDWithCursor(ctx context.Context, userID, lastID uint64, limit int) ([]model.UserOperationLog, error) {
	var records []model.UserOperationLog
	var builder strings.Builder
	builder.WriteString(`SELECT
	id,
	user_id,
	action,
	result,
	target_id,
	target_code,
	created_at
FROM user_operation_logs
WHERE user_id = ?`)

	args := []any{userID}
	if lastID > 0 {
		builder.WriteString("\n  AND id > ?")
		args = append(args, lastID)
	}

	builder.WriteString("\nORDER BY id ASC\nLIMIT ?\n")
	args = append(args, limit)

	if err := d.conn.QueryRowsPartialCtx(ctx, &records, builder.String(), args...); err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return records, nil
}
