package dao

import (
	"context"
	"errors"
	"fmt"

	"mysurl1/internal/model"
	"mysurl1/internal/utils"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ShortLinkDAO struct {
	conn sqlx.SqlConn
}

// NewShortLinkDAO creates a DAO instance for short link persistence operations.
func NewShortLinkDAO(conn sqlx.SqlConn) *ShortLinkDAO {
	return &ShortLinkDAO{conn: conn}
}

// FindAvailableByHash loads non-deleted short links by url hash for reuse checks.
func (d *ShortLinkDAO) FindAvailableByHash(ctx context.Context, urlHash string) ([]model.ShortLink, error) {
	var records []model.ShortLink
	query := `
SELECT
	id,
	user_id,
	short_code,
	original_url,
	url_hash
FROM short_links
WHERE url_hash = ?
  AND deleted_at IS NULL
`

	if err := d.conn.QueryRowsPartialCtx(ctx, &records, query, urlHash); err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return records, nil
}

// FindAvailableByCode loads a non-deleted short link by short code.
func (d *ShortLinkDAO) FindAvailableByCode(ctx context.Context, code string) (*model.ShortLink, error) {
	var record model.ShortLink
	query := `
SELECT
	id,
	original_url
FROM short_links
WHERE short_code = ?
  AND deleted_at IS NULL
LIMIT 1
`

	if err := d.conn.QueryRowPartialCtx(ctx, &record, query, code); err != nil {
		return nil, err
	}

	return &record, nil
}

// FindAvailableByOriginalURL loads a non-deleted short link by user id and normalized original URL.
func (d *ShortLinkDAO) FindAvailableByOriginalURL(ctx context.Context, userID uint64, normalizedURL string) (*model.ShortLink, error) {
	var record model.ShortLink
	query := `
SELECT
	id,
	short_code,
	original_url
FROM short_links
WHERE user_id = ?
  AND original_url = ?
  AND deleted_at IS NULL
LIMIT 1
`

	if err := d.conn.QueryRowPartialCtx(ctx, &record, query, userID, normalizedURL); err != nil {
		return nil, err
	}

	return &record, nil
}

// ListByUserID loads non-deleted short links for a specific user ordered by newest first.
func (d *ShortLinkDAO) ListByUserID(ctx context.Context, userID uint64) ([]model.ShortLink, error) {
	var records []model.ShortLink
	query := `
SELECT
	id,
	user_id,
	short_code,
	original_url,
	visit_count,
	created_at
FROM short_links
WHERE user_id = ?
  AND deleted_at IS NULL
ORDER BY id DESC
`

	if err := d.conn.QueryRowsPartialCtx(ctx, &records, query, userID); err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return records, nil
}

// Insert creates a new short link record with the normalized long URL.
func (d *ShortLinkDAO) Insert(ctx context.Context, userID *uint64, shortCode, normalizedURL, urlHash string) error {
	query := `
INSERT INTO short_links (
	user_id,
	short_code,
	original_url,
	url_hash
) VALUES (?, ?, ?, ?)
`

	_, err := d.conn.ExecCtx(ctx, query, userID, shortCode, normalizedURL, urlHash)
	return err
}

// CreateWithShortCode inserts a short link row with the provided short code.
func (d *ShortLinkDAO) CreateWithShortCode(ctx context.Context, userID *uint64, shortCode, normalizedURL, urlHash string) (string, error) {
	if d == nil || d.conn == nil {
		return "", fmt.Errorf("short link dao is not configured")
	}

	if err := d.Insert(ctx, userID, shortCode, normalizedURL, urlHash); err != nil {
		return "", err
	}

	return shortCode, nil
}

// CreateWithAutoIncrement inserts a short link row first, then derives and updates
// the short code from the generated auto increment id in one transaction.
func (d *ShortLinkDAO) CreateWithAutoIncrement(ctx context.Context, userID *uint64, normalizedURL, urlHash string) (string, error) {
	if d == nil || d.conn == nil {
		return "", fmt.Errorf("short link dao is not configured")
	}

	var shortCode string

	err := d.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		insertQuery := `
INSERT INTO short_links (
	user_id,
	short_code,
	original_url,
	url_hash
) VALUES (?, NULL, ?, ?)
`

		result, err := session.ExecCtx(ctx, insertQuery, userID, normalizedURL, urlHash)
		if err != nil {
			return err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		if id <= 0 {
			return fmt.Errorf("invalid short link id: %d", id)
		}

		shortCode = utils.EncodeBase62(uint64(id))
		updateQuery := "UPDATE short_links SET short_code = ? WHERE id = ?"
		if _, err := session.ExecCtx(ctx, updateQuery, shortCode, id); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	return shortCode, nil
}

// IncrementVisitCount increases the visit counter for a short link record.
func (d *ShortLinkDAO) IncrementVisitCount(ctx context.Context, id uint64) error {
	_, err := d.conn.ExecCtx(ctx, "UPDATE short_links SET visit_count = visit_count + 1 WHERE id = ?", id)
	return err
}

// AddVisitCount increases visit_count by the given delta.
func (d *ShortLinkDAO) AddVisitCount(ctx context.Context, id, delta uint64) error {
	_, err := d.conn.ExecCtx(ctx, "UPDATE short_links SET visit_count = visit_count + ? WHERE id = ?", delta, id)
	return err
}
