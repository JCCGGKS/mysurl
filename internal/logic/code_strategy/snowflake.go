package codestrategy

import (
	"context"
	"errors"
)

type SnowflakeGenerator struct{}

func NewSnowflakeGenerator() *SnowflakeGenerator {
	return &SnowflakeGenerator{}
}

func (g *SnowflakeGenerator) Provider() string {
	return ProviderSnowflake
}

func (g *SnowflakeGenerator) NextCode(_ context.Context, _, _ string) (string, error) {
	return "", errors.New("snowflake generator is not implemented")
}
