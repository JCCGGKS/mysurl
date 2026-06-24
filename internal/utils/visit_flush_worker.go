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
	count uint64
}

type VisitCountCache interface {
	// ScanVisitCountKeys scans visit count keys by cursor with a bounded count.
	ScanVisitCountKeys(ctx context.Context, cursor uint64, count int64) ([]string, uint64, error)
	// GetVisitCounts loads full visit counts for the given Redis keys in batch.
	GetVisitCounts(ctx context.Context, keys []string) (map[string]uint64, error)
}

type VisitCountWriter interface {
	UpsertVisitCounts(ctx context.Context, counts map[uint64]uint64) error
}

func StartVisitFlushWorker(db sqlx.SqlConn, visitCountCache VisitCountCache, visitCountWriter VisitCountWriter, conf config.VisitFlushConf) {
	if db == nil || IsNil(visitCountCache) || IsNil(visitCountWriter) {
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
			if err := flushVisitCounts(context.Background(), db, visitCountCache, visitCountWriter, conf.Batch); err != nil {
				logx.Errorf("flush visit counts failed: %v", err)
			}
		}
	}()
}

func flushVisitCounts(ctx context.Context, db sqlx.SqlConn, visitCountCache VisitCountCache, visitCountWriter VisitCountWriter, batch int64) error {
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
			counts, err := visitCountCache.GetVisitCounts(ctx, keys)
			if err != nil {
				return err
			}

			for _, key := range keys {
				count := counts[key]
				if count == 0 {
					continue
				}

				id, err := parseVisitCountKey(key)
				if err != nil {
					logx.Errorf("parse visit count key failed, key=%s err=%v", key, err)
					continue
				}

				items = append(items, visitDelta{
					id:    id,
					count: count,
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

	counts := make(map[uint64]uint64, len(items))
	var totalCount uint64
	for _, item := range items {
		counts[item.id] = item.count
		totalCount += item.count
	}

	if err := db.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		_ = session
		return visitCountWriter.UpsertVisitCounts(ctx, counts)
	}); err != nil {
		return err
	}

	logx.Infof("flush visit counts success, keys=%d total_count=%d", len(items), totalCount)
	return nil
}

func parseVisitCountKey(key string) (uint64, error) {
	raw := strings.TrimPrefix(key, visitKeyPrefix)
	return strconv.ParseUint(raw, 10, 64)
}
