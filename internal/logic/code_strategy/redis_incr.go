package codestrategy

import (
	"context"
	"errors"

	goredis "github.com/redis/go-redis/v9"
)

type RedisIncrGenerator struct {
	redis *goredis.Client
}

func NewRedisIncrGenerator(redis *goredis.Client) *RedisIncrGenerator {
	return &RedisIncrGenerator{redis: redis}
}

func (g *RedisIncrGenerator) Provider() string {
	return ProviderRedisIncr
}

func (g *RedisIncrGenerator) NextCode(_ context.Context) (string, error) {
	return "", errors.New("redis incr generator is not implemented")
}
