package dao

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"

	goredis "github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

const bloomShortCodeKey = "shortlink:bloom:short_code"

type ShortToLongCacheValue struct {
	ID          uint64 `json:"id"`
	OriginalURL string `json:"original_url"`
}

type ShortLinkCache struct {
	redis         *goredis.Client
	bloomDisabled atomic.Bool
	bloomWarned   atomic.Bool
}

func NewShortLinkCache(redis *goredis.Client) *ShortLinkCache {
	return &ShortLinkCache{redis: redis}
}

func (c *ShortLinkCache) GetShortToLong(ctx context.Context, shortCode string) (*ShortToLongCacheValue, error) {
	if c == nil || c.redis == nil {
		return nil, nil
	}

	raw, err := c.redis.Get(ctx, shortToLongCacheKey(shortCode)).Result()
	if errors.Is(err, goredis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var value ShortToLongCacheValue
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		return nil, err
	}

	return &value, nil
}

func (c *ShortLinkCache) SetShortToLong(ctx context.Context, shortCode string, value ShortToLongCacheValue) error {
	if c == nil || c.redis == nil {
		return nil
	}

	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.redis.Set(ctx, shortToLongCacheKey(shortCode), payload, 0).Err()
}

func (c *ShortLinkCache) GetLongToShort(ctx context.Context, userID uint64, normalizedURL string) (string, error) {
	if c == nil || c.redis == nil {
		return "", nil
	}

	shortCode, err := c.redis.Get(ctx, longToShortCacheKey(userID, normalizedURL)).Result()
	if errors.Is(err, goredis.Nil) {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return shortCode, nil
}

func (c *ShortLinkCache) SetLongToShort(ctx context.Context, userID uint64, normalizedURL, shortCode string) error {
	if c == nil || c.redis == nil {
		return nil
	}

	return c.redis.Set(ctx, longToShortCacheKey(userID, normalizedURL), shortCode, 0).Err()
}

func (c *ShortLinkCache) ShortCodeBloomExists(ctx context.Context, shortCode string) (bool, error) {
	if c == nil || c.redis == nil {
		return true, nil
	}
	if c.bloomDisabled.Load() {
		return true, nil
	}

	result, err := c.redis.Do(ctx, "BF.EXISTS", bloomShortCodeKey, shortCode).Bool()
	if err != nil {
		if c.handleBloomUnavailable(err) {
			return true, nil
		}

		return false, err
	}

	return result, nil
}

func (c *ShortLinkCache) ShortCodeBloomAdd(ctx context.Context, shortCode string) error {
	if c == nil || c.redis == nil {
		return nil
	}
	if c.bloomDisabled.Load() {
		return nil
	}

	err := c.redis.Do(ctx, "BF.ADD", bloomShortCodeKey, shortCode).Err()
	if err == nil {
		return nil
	}
	if c.handleBloomUnavailable(err) {
		return nil
	}

	return err
}

func (c *ShortLinkCache) IncrVisitCount(ctx context.Context, id uint64) error {
	if c == nil || c.redis == nil {
		return nil
	}

	return c.redis.Incr(ctx, visitCountKey(id)).Err()
}

func (c *ShortLinkCache) ScanVisitCountKeys(ctx context.Context, cursor uint64, count int64) ([]string, uint64, error) {
	if c == nil || c.redis == nil {
		return nil, 0, nil
	}

	keys, nextCursor, err := c.redis.Scan(ctx, cursor, "shortlink:visit:*", count).Result()
	if err != nil {
		return nil, 0, err
	}

	return keys, nextCursor, nil
}

func (c *ShortLinkCache) GetVisitCountDelta(ctx context.Context, key string) (uint64, error) {
	if c == nil || c.redis == nil {
		return 0, nil
	}

	raw, err := c.redis.Get(ctx, key).Result()
	if errors.Is(err, goredis.Nil) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	value, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, err
	}

	return value, nil
}

func (c *ShortLinkCache) DeleteVisitCountKey(ctx context.Context, key string) error {
	if c == nil || c.redis == nil {
		return nil
	}

	return c.redis.Del(ctx, key).Err()
}

func shortToLongCacheKey(shortCode string) string {
	return fmt.Sprintf("shortlink:code:%s", shortCode)
}

func longToShortCacheKey(userID uint64, normalizedURL string) string {
	return fmt.Sprintf("shortlink:user:%d:long:%s", userID, normalizedURL)
}

func visitCountKey(id uint64) string {
	return fmt.Sprintf("shortlink:visit:%d", id)
}

func (c *ShortLinkCache) handleBloomUnavailable(err error) bool {
	if !isBloomUnavailableError(err) {
		return false
	}

	c.bloomDisabled.Store(true)
	if c.bloomWarned.CompareAndSwap(false, true) {
		logx.Errorf("redis bloom unavailable, fallback to cache+mysql path: %v", err)
	}

	return true
}

func isBloomUnavailableError(err error) bool {
	if err == nil {
		return false
	}

	message := strings.ToLower(err.Error())
	if strings.Contains(message, "unknown command") && strings.Contains(message, "bf.") {
		return true
	}
	if strings.Contains(message, "module") && strings.Contains(message, "not loaded") {
		return true
	}

	return false
}
