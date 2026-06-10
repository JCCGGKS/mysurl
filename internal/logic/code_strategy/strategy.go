package codestrategy

import (
	"context"
	"fmt"
	"sync"

	"mysurl1/internal/utils"
)

const (
	ProviderMySQLAutoIncrement = "mysql_auto_increment"
	ProviderRedisIncr          = "redis_incr"
	ProviderSnowflake          = "snowflake"
)

type Generator interface {
	Provider() string
	NextCode(ctx context.Context) (string, error)
}

type CodeManager struct {
	mu         sync.RWMutex
	generators map[string]Generator
}

func NewCodeManager() *CodeManager {
	return &CodeManager{
		generators: make(map[string]Generator),
	}
}

func (m *CodeManager) Register(generator Generator) {
	if generator == nil {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.generators[generator.Provider()] = generator
}

func (m *CodeManager) Get(provider string) (Generator, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	generator, ok := m.generators[provider]
	if !ok {
		return nil, fmt.Errorf("unsupported short code provider: %s", provider)
	}

	return generator, nil
}

func (m *CodeManager) NextCode(ctx context.Context, provider string) (string, error) {
	if m == nil {
		return "", fmt.Errorf("code manager is nil")
	}
	if provider == "" {
		provider = ProviderMySQLAutoIncrement
	}
	if provider == ProviderMySQLAutoIncrement {
		return "", fmt.Errorf("provider %s does not support pre-generated short codes", ProviderMySQLAutoIncrement)
	}

	generator, err := m.Get(provider)
	if err != nil {
		return "", err
	}

	return generator.NextCode(ctx)
}

func BuildCodeFromID(id uint64) string {
	return utils.EncodeBase62(id)
}
