package codestrategy

import (
	"context"
	"errors"

	"mysurl1/internal/utils"

	goredis "github.com/redis/go-redis/v9"
)

type RedisIncrGenerator struct {
	redis *goredis.Client
}

const redisShortCodeSequenceKey = "shortlink:code:seq"

func NewRedisIncrGenerator(redis *goredis.Client) *RedisIncrGenerator {
	return &RedisIncrGenerator{
		redis: redis,
	}
}

func (g *RedisIncrGenerator) Provider() string {
	return ProviderRedisIncr
}

func (g *RedisIncrGenerator) NextCode(ctx context.Context) (string, error) {
	if g == nil || g.redis == nil {
		return "", errors.New("redis incr generator client is not configured")
	}

	sequenceID, err := g.redis.Incr(ctx, redisShortCodeSequenceKey).Uint64()
	if err != nil {
		return "", err
	}

	shortCode := utils.EncodeBase62(sequenceID)
	return shortCode, nil
}
