package codestrategy

import (
	"context"
	"errors"

	"mysurl1/internal/dao"
	"mysurl1/internal/utils"

	goredis "github.com/redis/go-redis/v9"
)

type RedisIncrGenerator struct {
	redis *goredis.Client
	dao   *dao.ShortLinkDAO
}

const redisShortCodeSequenceKey = "shortlink:code:seq"

func NewRedisIncrGenerator(redis *goredis.Client, shortLinkDAO *dao.ShortLinkDAO) *RedisIncrGenerator {
	return &RedisIncrGenerator{
		redis: redis,
		dao:   shortLinkDAO,
	}
}

func (g *RedisIncrGenerator) Provider() string {
	return ProviderRedisIncr
}

func (g *RedisIncrGenerator) NextCode(ctx context.Context, normalizedURL, urlHash string) (string, error) {
	if g == nil || g.redis == nil {
		return "", errors.New("redis incr generator client is not configured")
	}
	if g.dao == nil {
		return "", errors.New("redis incr generator dao is not configured")
	}

	sequenceID, err := g.redis.Incr(ctx, redisShortCodeSequenceKey).Uint64()
	if err != nil {
		return "", err
	}

	shortCode := utils.EncodeBase62(sequenceID)
	return g.dao.CreateWithShortCode(ctx, shortCode, normalizedURL, urlHash)
}
