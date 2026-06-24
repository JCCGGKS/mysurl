package dao

import (
	"context"
	"errors"
	"strings"

	"mysurl1/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type VisitStatDAO struct {
	conn sqlx.SqlConn
}

func NewVisitStatDAO(conn sqlx.SqlConn) *VisitStatDAO {
	return &VisitStatDAO{conn: conn}
}

func (d *VisitStatDAO) GetVisitCount(ctx context.Context, shortLinkID uint64) (uint64, error) {
	if d == nil || d.conn == nil {
		return 0, nil
	}

	var record model.VisitStat
	query := `
SELECT
	short_link_id,
	visit_count
FROM visit_stats
WHERE short_link_id = ?
LIMIT 1
`

	if err := d.conn.QueryRowPartialCtx(ctx, &record, query, shortLinkID); err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return 0, nil
		}
		return 0, err
	}

	return record.VisitCount, nil
}

func (d *VisitStatDAO) GetVisitCountsByIDs(ctx context.Context, shortLinkIDs []uint64) (map[uint64]uint64, error) {
	results := make(map[uint64]uint64, len(shortLinkIDs))
	if d == nil || d.conn == nil || len(shortLinkIDs) == 0 {
		return results, nil
	}

	var (
		records      []model.VisitStat
		placeholders = make([]string, 0, len(shortLinkIDs))
		args         = make([]any, 0, len(shortLinkIDs))
		builder      strings.Builder
	)

	builder.WriteString(`
SELECT
	short_link_id,
	visit_count
FROM visit_stats
WHERE short_link_id IN (`)
	for _, shortLinkID := range shortLinkIDs {
		placeholders = append(placeholders, "?")
		args = append(args, shortLinkID)
	}
	builder.WriteString(strings.Join(placeholders, ", "))
	builder.WriteString(")\n")

	if err := d.conn.QueryRowsPartialCtx(ctx, &records, builder.String(), args...); err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return results, nil
		}
		return nil, err
	}

	for _, record := range records {
		results[record.ShortLinkID] = record.VisitCount
	}

	return results, nil
}

func (d *VisitStatDAO) UpsertVisitCounts(ctx context.Context, counts map[uint64]uint64) error {
	if d == nil || d.conn == nil || len(counts) == 0 {
		return nil
	}

	var (
		builder strings.Builder
		args    = make([]any, 0, len(counts)*2)
	)

	builder.WriteString(`
INSERT INTO visit_stats (
	short_link_id,
	visit_count
) VALUES `)

	index := 0
	for shortLinkID, visitCount := range counts {
		if index > 0 {
			builder.WriteString(",")
		}
		builder.WriteString("(?, ?)")
		args = append(args, shortLinkID, visitCount)
		index++
	}

	builder.WriteString(`
ON DUPLICATE KEY UPDATE
	visit_count = VALUES(visit_count)
`)

	_, err := d.conn.ExecCtx(ctx, builder.String(), args...)
	return err
}
