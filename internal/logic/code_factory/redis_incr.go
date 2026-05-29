package codefactory

import (
	"context"
	"errors"
)

type RedisIncrGenerator struct{}

func NewRedisIncrGenerator() *RedisIncrGenerator {
	return &RedisIncrGenerator{}
}

func (g *RedisIncrGenerator) Provider() string {
	return ProviderRedisIncr
}

func (g *RedisIncrGenerator) NextCode(_ context.Context) (string, error) {
	return "", errors.New("redis incr generator is not implemented")
}
