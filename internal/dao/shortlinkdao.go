package dao

import (
	"context"
	"errors"

	"mysurl1/internal/model"

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
	short_code,
	original_url,
	url_hash,
	visit_count,
	status
FROM short_links
WHERE url_hash = ?
  AND deleted_at IS NULL
`

	if err := d.conn.QueryRowsCtx(ctx, &records, query, urlHash); err != nil {
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
	short_code,
	original_url,
	url_hash,
	visit_count,
	status
FROM short_links
WHERE short_code = ?
  AND deleted_at IS NULL
LIMIT 1
`

	if err := d.conn.QueryRowCtx(ctx, &record, query, code); err != nil {
		return nil, err
	}

	return &record, nil
}

// ExistsByCode checks whether a short code already exists in storage.
func (d *ShortLinkDAO) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var result struct {
		Cnt int64 `db:"cnt"`
	}

	query := `
SELECT COUNT(1) AS cnt
FROM short_links
WHERE short_code = ?
LIMIT 1
`

	if err := d.conn.QueryRowCtx(ctx, &result, query, code); err != nil {
		return false, err
	}

	return result.Cnt > 0, nil
}

// Insert creates a new short link record with the default active status.
func (d *ShortLinkDAO) Insert(ctx context.Context, shortCode, originalURL, urlHash string) error {
	query := `
INSERT INTO short_links (
	short_code,
	original_url,
	url_hash
) VALUES (?, ?, ?)
`

	_, err := d.conn.ExecCtx(ctx, query, shortCode, originalURL, urlHash)
	return err
}

// IncrementVisitCount increases the visit counter for a short link record.
func (d *ShortLinkDAO) IncrementVisitCount(ctx context.Context, id uint64) error {
	_, err := d.conn.ExecCtx(ctx, "UPDATE short_links SET visit_count = visit_count + 1 WHERE id = ?", id)
	return err
}
