package utils

import (
	"context"
	"strconv"
	"strings"
	"time"

	"mysurl1/internal/config"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const (
	visitKeyPrefix = "shortlink:visit:"
)

type visitDelta struct {
	id    uint64
	delta uint64
	key   string
}

type VisitCountCache interface {
	// ScanVisitCountKeys scans visit count keys by cursor with a bounded count.
	ScanVisitCountKeys(ctx context.Context, cursor uint64, count int64) ([]string, uint64, error)
	// GetVisitCountDeltas loads visit count deltas for the given Redis keys in batch.
	GetVisitCountDeltas(ctx context.Context, keys []string) (map[string]uint64, error)
	// DeleteVisitCountKeys removes flushed visit count keys from Redis in batch.
	DeleteVisitCountKeys(ctx context.Context, keys []string) error
}

func StartVisitFlushWorker(db sqlx.SqlConn, visitCountCache VisitCountCache, conf config.VisitFlushConf) {
	if db == nil || IsNil(visitCountCache) {
		return
	}
	if conf.Interval <= 0 {
		conf.Interval = 5 * time.Second
	}
	if conf.Batch <= 0 {
		conf.Batch = 100
	}

	go func() {
		ticker := time.NewTicker(conf.Interval)
		defer ticker.Stop()

		for range ticker.C {
			if err := flushVisitCounts(context.Background(), db, visitCountCache, conf.Batch); err != nil {
				logx.Errorf("flush visit counts failed: %v", err)
			}
		}
	}()
}

func flushVisitCounts(ctx context.Context, db sqlx.SqlConn, visitCountCache VisitCountCache, batch int64) error {
	var (
		cursor uint64
		items  []visitDelta
	)

	for {
		keys, nextCursor, err := visitCountCache.ScanVisitCountKeys(ctx, cursor, batch)
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			deltas, err := visitCountCache.GetVisitCountDeltas(ctx, keys)
			if err != nil {
				return err
			}

			for _, key := range keys {
				delta := deltas[key]
				if delta == 0 {
					continue
				}

				id, err := parseVisitCountKey(key)
				if err != nil {
					logx.Errorf("parse visit count key failed, key=%s err=%v", key, err)
					continue
				}

				items = append(items, visitDelta{
					id:    id,
					delta: delta,
					key:   key,
				})
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	if len(items) == 0 {
		return nil
	}

	updateQuery, updateArgs := buildVisitCountBatchUpdate(items)
	if err := db.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		_, err := session.ExecCtx(ctx, updateQuery, updateArgs...)
		return err
	}); err != nil {
		return err
	}

	var (
		totalDelta   uint64
		keysToDelete = make([]string, 0, len(items))
	)
	for _, item := range items {
		totalDelta += item.delta
		keysToDelete = append(keysToDelete, item.key)
	}

	if err := visitCountCache.DeleteVisitCountKeys(ctx, keysToDelete); err != nil {
		logx.Errorf("delete visit count keys failed, keys=%d err=%v", len(keysToDelete), err)
	}

	logx.Infof("flush visit counts success, keys=%d total_delta=%d", len(items), totalDelta)
	return nil
}

func parseVisitCountKey(key string) (uint64, error) {
	raw := strings.TrimPrefix(key, visitKeyPrefix)
	return strconv.ParseUint(raw, 10, 64)
}

func buildVisitCountBatchUpdate(items []visitDelta) (string, []any) {
	var builder strings.Builder
	args := make([]any, 0, len(items)*3)

	builder.WriteString("UPDATE short_links SET visit_count = visit_count + CASE id")
	for _, item := range items {
		builder.WriteString(" WHEN ? THEN ?")
		args = append(args, item.id, item.delta)
	}

	builder.WriteString(" ELSE 0 END WHERE id IN (")
	for i, item := range items {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString("?")
		args = append(args, item.id)
	}
	builder.WriteString(")")

	return builder.String(), args
}
