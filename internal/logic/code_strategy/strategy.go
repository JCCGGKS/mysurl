package codestrategy

import (
	"context"
	"fmt"

	"mysurl1/internal/config"
	"mysurl1/internal/dao"
)

const (
	ProviderMySQLAutoIncrement = "mysql_auto_increment"
	ProviderRedisIncr          = "redis_incr"
	ProviderSnowflake          = "snowflake"
)

type NextCodeInput struct {
	DAO *dao.ShortLinkDAO
}

type CodeGenerator interface {
	Provider() string
	NextCode(ctx context.Context, input NextCodeInput) (string, error)
}

type GenerateShortCodeFunc func(ctx context.Context) (string, error)

var strategyRegistry = map[string]CodeGenerator{
	ProviderMySQLAutoIncrement: &MySQLAutoIncrementGenerator{},
	ProviderRedisIncr:          &RedisIncrGenerator{},
	ProviderSnowflake:          &SnowflakeGenerator{},
}

func NewCodeGenerator(provider string) (CodeGenerator, error) {
	generator, ok := strategyRegistry[provider]
	if !ok {
		return nil, fmt.Errorf("unsupported short code provider: %s", provider)
	}

	return generator, nil
}

func NewCodeGenService(cfg config.ShortConf) (CodeGenerator, error) {
	provider := cfg.Provider
	if provider == "" {
		provider = ProviderMySQLAutoIncrement
	}

	generator, err := NewCodeGenerator(provider)
	if err != nil {
		return nil, fmt.Errorf("init short code generator failed: %w", err)
	}

	return generator, nil
}

func GenerateShortCode(ctx context.Context, generator CodeGenerator, input NextCodeInput) (string, error) {
	if generator == nil {
		return "", fmt.Errorf("short code generator is nil")
	}

	return generator.NextCode(ctx, input)
}

func BuildGenerateShortCodeFunc(generator CodeGenerator, input NextCodeInput) GenerateShortCodeFunc {
	return func(ctx context.Context) (string, error) {
		return GenerateShortCode(ctx, generator, input)
	}
}
