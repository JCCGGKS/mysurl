package codestrategy

import (
	"context"
	"errors"
)

type SnowflakeGenerator struct{}



func (g *SnowflakeGenerator) Provider() string {
	return ProviderSnowflake
}

func (g *SnowflakeGenerator) NextCode(_ context.Context, _ NextCodeInput) (string, error) {
	return "", errors.New("snowflake generator is not implemented")
}
