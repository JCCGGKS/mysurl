package dao

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	goredis "github.com/redis/go-redis/v9"
)

const bloomNormalizedURLKey = "shortlink:bloom:normalized_long_url"

type ShortToLongCacheValue struct {
	ID          uint64 `json:"id"`
	OriginalURL string `json:"original_url"`
}

type ShortLinkCache struct {
	redis *goredis.Client
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

	result, err := c.redis.Do(ctx, "BF.EXISTS", bloomNormalizedURLKey, normalizedURL).Int64()
	if err != nil {
		return false, err
	}

	return result == 1, nil
}

func (c *ShortLinkCache) BloomAdd(ctx context.Context, normalizedURL string) error {
	if c == nil || c.redis == nil {
		return nil
	}

	return c.redis.Do(ctx, "BF.ADD", bloomNormalizedURLKey, normalizedURL).Err()
}

func shortToLongCacheKey(shortCode string) string {
	return fmt.Sprintf("shortlink:code:%s", shortCode)
}

func longToShortCacheKey(normalizedURL string) string {
	return fmt.Sprintf("shortlink:long:%s", normalizedURL)
}
