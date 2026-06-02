package utils

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const (
	visitFlushInterval = 5 * time.Second
	visitFlushBatch    = 100
	visitKeyPrefix     = "shortlink:visit:"
)

type VisitCountCache interface {
	// ScanVisitCountKeys scans visit count keys by cursor with a bounded count.
	ScanVisitCountKeys(ctx context.Context, cursor uint64, count int64) ([]string, uint64, error)
	// GetVisitCountDelta loads the current visit count delta from a Redis key.
	GetVisitCountDelta(ctx context.Context, key string) (uint64, error)
	// DeleteVisitCountKey removes a flushed visit count key from Redis.
	DeleteVisitCountKey(ctx context.Context, key string) error
}

func StartVisitFlushWorker(db sqlx.SqlConn, visitCountCache VisitCountCache) {
	if db == nil || IsNil(visitCountCache) {
		return
	}

	go func() {
		ticker := time.NewTicker(visitFlushInterval)
		defer ticker.Stop()

		for range ticker.C {
			if err := flushVisitCounts(context.Background(), db, visitCountCache); err != nil {
				logx.Errorf("flush visit counts failed: %v", err)
			}
		}
	}()
}

func flushVisitCounts(ctx context.Context, db sqlx.SqlConn, visitCountCache VisitCountCache) error {
	keys, _, err := visitCountCache.ScanVisitCountKeys(ctx, 0, visitFlushBatch)
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}

	type visitDelta struct {
		id    uint64
		delta uint64
		key   string
	}

	items := make([]visitDelta, 0, len(keys))
	for _, key := range keys {
		delta, err := visitCountCache.GetVisitCountDelta(ctx, key)
		if err != nil {
			logx.Errorf("get visit count delta failed, key=%s err=%v", key, err)
			continue
		}
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

	if len(items) == 0 {
		return nil
	}

	err = db.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		for _, item := range items {
			if _, err := session.ExecCtx(ctx, "UPDATE short_links SET visit_count = visit_count + ? WHERE id = ?", item.delta, item.id); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	var totalDelta uint64
	for _, item := range items {
		totalDelta += item.delta
		if err := visitCountCache.DeleteVisitCountKey(ctx, item.key); err != nil {
			logx.Errorf("delete visit count key failed, key=%s err=%v", item.key, err)
		}
	}

	logx.Infof("flush visit counts success, keys=%d total_delta=%d", len(items), totalDelta)
	return nil
}

func parseVisitCountKey(key string) (uint64, error) {
	raw := strings.TrimPrefix(key, visitKeyPrefix)
	return strconv.ParseUint(raw, 10, 64)
}
