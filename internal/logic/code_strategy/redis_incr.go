package codestrategy

import (
	"context"
	"errors"
)

type RedisIncrGenerator struct{}



func (g *RedisIncrGenerator) Provider() string {
	return ProviderRedisIncr
}

func (g *RedisIncrGenerator) NextCode(_ context.Context, _ NextCodeInput) (string, error) {
	return "", errors.New("redis incr generator is not implemented")
}
