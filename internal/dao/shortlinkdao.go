package dao

import (
	"context"
	"errors"
	"fmt"
	"strings"

	codestrategy "mysurl1/internal/logic/code_strategy"
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

// FindAvailableByOriginalURLs loads non-deleted short links by user id and a batch of normalized URLs.
func (d *ShortLinkDAO) FindAvailableByOriginalURLs(ctx context.Context, userID uint64, normalizedURLs []string) ([]model.ShortLink, error) {
	if len(normalizedURLs) == 0 {
		return nil, nil
	}

	var (
		records      []model.ShortLink
		placeholders = make([]string, 0, len(normalizedURLs))
		args         = make([]any, 0, len(normalizedURLs)+1)
		builder      strings.Builder
	)

	args = append(args, userID)
	for _, normalizedURL := range normalizedURLs {
		placeholders = append(placeholders, "?")
		args = append(args, normalizedURL)
	}

	builder.WriteString(`
SELECT
	id,
	short_code,
	original_url
FROM short_links
WHERE user_id = ?
  AND deleted_at IS NULL
  AND original_url IN (`)
	builder.WriteString(strings.Join(placeholders, ", "))
	builder.WriteString(")\n")

	if err := d.conn.QueryRowsPartialCtx(ctx, &records, builder.String(), args...); err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return records, nil
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

func (d *ShortLinkDAO) ListByUserIDWithCursor(ctx context.Context, userID uint64, shortCode, originalURL string, lastID uint64, limit int) ([]model.ShortLink, error) {
	var records []model.ShortLink
	query, args := buildUserLinkListQuery(false, userID, shortCode, originalURL, lastID)
	query += "\nORDER BY id DESC\nLIMIT ?\n"
	args = append(args, limit)

	if err := d.conn.QueryRowsPartialCtx(ctx, &records, query, args...); err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return records, nil
}

func (d *ShortLinkDAO) CountByUserID(ctx context.Context, userID uint64, shortCode, originalURL string) (int64, error) {
	var result struct {
		Total int64 `db:"total"`
	}

	query, args := buildUserLinkListQuery(true, userID, shortCode, originalURL, 0)
	if err := d.conn.QueryRowPartialCtx(ctx, &result, query, args...); err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return 0, nil
		}

		return 0, err
	}

	return result.Total, nil
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

// BatchInsert inserts multiple short link records in one statement.
func (d *ShortLinkDAO) BatchInsert(ctx context.Context, userID *uint64, records []model.ShortLink) error {
	if len(records) == 0 {
		return nil
	}

	var (
		builder strings.Builder
		args    = make([]any, 0, len(records)*4)
	)

	builder.WriteString(`
INSERT INTO short_links (
	user_id,
	short_code,
	original_url,
	url_hash
) VALUES `)

	for i, record := range records {
		if i > 0 {
			builder.WriteString(",")
		}
		builder.WriteString("(?, ?, ?, ?)")
		args = append(args, userID, record.ShortCode, record.OriginalURL, record.URLHash)
	}

	_, err := d.conn.ExecCtx(ctx, builder.String(), args...)
	return err
}

// BatchCreateWithAutoIncrement inserts multiple rows and derives short codes from generated ids in one transaction.
func (d *ShortLinkDAO) BatchCreateWithAutoIncrement(ctx context.Context, userID *uint64, records []model.ShortLink) ([]model.ShortLink, error) {
	if d == nil || d.conn == nil {
		return nil, fmt.Errorf("short link dao is not configured")
	}
	if len(records) == 0 {
		return nil, nil
	}

	created := make([]model.ShortLink, len(records))
	copy(created, records)

	err := d.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		var (
			builder strings.Builder
			args    = make([]any, 0, len(created)*3)
		)

		builder.WriteString(`
INSERT INTO short_links (
	user_id,
	short_code,
	original_url,
	url_hash
) VALUES `)

		for i, record := range created {
			if i > 0 {
				builder.WriteString(",")
			}
			builder.WriteString("(?, NULL, ?, ?)")
			args = append(args, userID, record.OriginalURL, record.URLHash)
		}

		result, err := session.ExecCtx(ctx, builder.String(), args...)
		if err != nil {
			return err
		}

		firstID, err := result.LastInsertId()
		if err != nil {
			return err
		}
		if firstID <= 0 {
			return fmt.Errorf("invalid first short link id: %d", firstID)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected != int64(len(created)) {
			return fmt.Errorf("unexpected rows affected: %d", rowsAffected)
		}

		var (
			updateBuilder strings.Builder
			updateArgs    = make([]any, 0, len(created)*2)
			idArgs        = make([]any, 0, len(created))
		)

		updateBuilder.WriteString("UPDATE short_links SET short_code = CASE id ")
		for i := range created {
			id := uint64(firstID) + uint64(i)
			code := codestrategy.BuildCodeFromID(id)
			created[i].ID = id
			created[i].ShortCode = code
			updateBuilder.WriteString("WHEN ? THEN ? ")
			updateArgs = append(updateArgs, id, code)
			idArgs = append(idArgs, id)
		}
		updateBuilder.WriteString("END WHERE id IN (")
		for i := range created {
			if i > 0 {
				updateBuilder.WriteString(",")
			}
			updateBuilder.WriteString("?")
		}
		updateBuilder.WriteString(")")
		updateArgs = append(updateArgs, idArgs...)

		if _, err := session.ExecCtx(ctx, updateBuilder.String(), updateArgs...); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return created, nil
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

		shortCode = codestrategy.BuildCodeFromID(uint64(id))
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

func buildUserLinkListQuery(countOnly bool, userID uint64, shortCode, originalURL string, lastID uint64) (string, []any) {
	var builder strings.Builder
	if countOnly {
		builder.WriteString("SELECT COUNT(1) AS total\n")
	} else {
		builder.WriteString(`SELECT
	id,
	user_id,
	short_code,
	original_url,
	visit_count,
	created_at
`)
	}

	builder.WriteString(`FROM short_links
WHERE user_id = ?
  AND deleted_at IS NULL`)

	args := []any{userID}
	if lastID > 0 {
		builder.WriteString("\n  AND id < ?")
		args = append(args, lastID)
	}

	shortCode = strings.TrimSpace(shortCode)
	if shortCode != "" {
		builder.WriteString("\n  AND short_code LIKE ?")
		args = append(args, "%"+shortCode+"%")
	}

	originalURL = strings.TrimSpace(originalURL)
	if originalURL != "" {
		builder.WriteString("\n  AND original_url LIKE ?")
		args = append(args, "%"+originalURL+"%")
	}

	return builder.String(), args
}
