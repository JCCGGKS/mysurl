package dao

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"

	goredis "github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

const bloomNormalizedURLKey = "shortlink:bloom:normalized_long_url"

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

func (c *ShortLinkCache) GetLongToShort(ctx context.Context, normalizedURL string) (string, error) {
	if c == nil || c.redis == nil {
		return "", nil
	}

	shortCode, err := c.redis.Get(ctx, longToShortCacheKey(normalizedURL)).Result()
	if errors.Is(err, goredis.Nil) {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return shortCode, nil
}

func (c *ShortLinkCache) SetLongToShort(ctx context.Context, normalizedURL, shortCode string) error {
	if c == nil || c.redis == nil {
		return nil
	}

	return c.redis.Set(ctx, longToShortCacheKey(normalizedURL), shortCode, 0).Err()
}

func (c *ShortLinkCache) BloomExists(ctx context.Context, normalizedURL string) (bool, error) {
	if c == nil || c.redis == nil {
		return true, nil
	}
	if c.bloomDisabled.Load() {
		return true, nil
	}

	result, err := c.redis.Do(ctx, "BF.EXISTS", bloomNormalizedURLKey, normalizedURL).Bool()
	if err != nil {
		if c.handleBloomUnavailable(err) {
			return true, nil
		}

		return false, err
	}

	return result, nil
}

func (c *ShortLinkCache) BloomAdd(ctx context.Context, normalizedURL string) error {
	if c == nil || c.redis == nil {
		return nil
	}
	if c.bloomDisabled.Load() {
		return nil
	}

	err := c.redis.Do(ctx, "BF.ADD", bloomNormalizedURLKey, normalizedURL).Err()
	if err == nil {
		return nil
	}
	if c.handleBloomUnavailable(err) {
		return nil
	}

	return err
}

func shortToLongCacheKey(shortCode string) string {
	return fmt.Sprintf("shortlink:code:%s", shortCode)
}

func longToShortCacheKey(normalizedURL string) string {
	return fmt.Sprintf("shortlink:long:%s", normalizedURL)
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
