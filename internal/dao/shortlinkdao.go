package dao

import (
	"context"
	"errors"
	"time"

	"mysurl1/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ShortLinkDAO struct {
	conn sqlx.SqlConn
}

func NewShortLinkDAO(conn sqlx.SqlConn) *ShortLinkDAO {
	return &ShortLinkDAO{conn: conn}
}

func (d *ShortLinkDAO) FindAvailableByHash(ctx context.Context, urlHash string, now time.Time) ([]model.ShortLink, error) {
	var records []model.ShortLink
	query := `
SELECT
	id,
	short_code,
	original_url,
	url_hash,
	visit_count,
	status,
	expires_at,
	created_at,
	updated_at,
	deleted_at
FROM short_links
WHERE url_hash = ?
  AND deleted_at IS NULL
  AND (expires_at IS NULL OR expires_at > ?)
`

	if err := d.conn.QueryRowsCtx(ctx, &records, query, urlHash, now); err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return records, nil
}

func (d *ShortLinkDAO) FindAvailableByCode(ctx context.Context, code string, now time.Time) (*model.ShortLink, error) {
	var record model.ShortLink
	query := `
SELECT
	id,
	short_code,
	original_url,
	url_hash,
	visit_count,
	status,
	expires_at,
	created_at,
	updated_at,
	deleted_at
FROM short_links
WHERE short_code = ?
  AND deleted_at IS NULL
  AND (expires_at IS NULL OR expires_at > ?)
LIMIT 1
`

	if err := d.conn.QueryRowCtx(ctx, &record, query, code, now); err != nil {
		return nil, err
	}

	return &record, nil
}

func (d *ShortLinkDAO) Insert(ctx context.Context, shortCode, originalURL, urlHash string, expiresAt *time.Time) error {
	query := `
INSERT INTO short_links (
	short_code,
	original_url,
	url_hash,
	visit_count,
	status,
	expires_at
) VALUES (?, ?, ?, 0, ?, ?)
`

	_, err := d.conn.ExecCtx(ctx, query, shortCode, originalURL, urlHash, model.ShortLinkStatusActive, expiresAt)
	return err
}

func (d *ShortLinkDAO) IncrementVisitCount(ctx context.Context, id uint64) error {
	_, err := d.conn.ExecCtx(ctx, "UPDATE short_links SET visit_count = visit_count + 1 WHERE id = ?", id)
	return err
}
