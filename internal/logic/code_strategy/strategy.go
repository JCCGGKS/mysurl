package codestrategy

import (
	"context"
	"fmt"
	"sync"
)

const (
	ProviderMySQLAutoIncrement = "mysql_auto_increment"
	ProviderRedisIncr          = "redis_incr"
	ProviderSnowflake          = "snowflake"
)

type CodeGenerator interface {
	Provider() string
	NextCode(ctx context.Context, normalizedURL, urlHash string) (string, error)
}

type CodeManager struct {
	provider   string
	mu         sync.RWMutex
	generators map[string]CodeGenerator
}

func NewCodeManager(provider string) *CodeManager {
	if provider == "" {
		provider = ProviderMySQLAutoIncrement
	}

	return &CodeManager{
		provider:   provider,
		generators: make(map[string]CodeGenerator),
	}
}

func (m *CodeManager) Register(generator CodeGenerator) {
	if generator == nil {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.generators[generator.Provider()] = generator
}

func (m *CodeManager) Get(provider string) (CodeGenerator, error) {
	if provider == "" {
		provider = ProviderMySQLAutoIncrement
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	generator, ok := m.generators[provider]
	if !ok {
		return nil, fmt.Errorf("unsupported short code provider: %s", provider)
	}

	return generator, nil
}

func (m *CodeManager) Provider() string {
	if m == nil || m.provider == "" {
		return ProviderMySQLAutoIncrement
	}

	return m.provider
}

func (m *CodeManager) GenerateShortCode(ctx context.Context, normalizedURL, urlHash string) (string, error) {
	if m == nil {
		return "", fmt.Errorf("code manager is nil")
	}

	generator, err := m.Get(m.provider)
	if err != nil {
		return "", err
	}

	return generator.NextCode(ctx, normalizedURL, urlHash)
}
