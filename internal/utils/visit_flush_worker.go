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

type VisitCountCache interface {
	// ScanVisitCountKeys scans visit count keys by cursor with a bounded count.
	ScanVisitCountKeys(ctx context.Context, cursor uint64, count int64) ([]string, uint64, error)
	// GetVisitCountDelta loads the current visit count delta from a Redis key.
	GetVisitCountDelta(ctx context.Context, key string) (uint64, error)
	// DeleteVisitCountKey removes a flushed visit count key from Redis.
	DeleteVisitCountKey(ctx context.Context, key string) error
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
	keys, _, err := visitCountCache.ScanVisitCountKeys(ctx, 0, batch)
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
