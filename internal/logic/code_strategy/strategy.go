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
	provider   string
	mu         sync.RWMutex
	generators map[string]Generator
}

func NewCodeManager(provider string) *CodeManager {
	if provider == "" {
		provider = ProviderMySQLAutoIncrement
	}

	return &CodeManager{
		provider:   provider,
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

func (m *CodeManager) IsMySQLAutoIncrement() bool {
	return m.Provider() == ProviderMySQLAutoIncrement
}

func (m *CodeManager) NextCode(ctx context.Context) (string, error) {
	if m == nil {
		return "", fmt.Errorf("code manager is nil")
	}
	if m.IsMySQLAutoIncrement() {
		return "", fmt.Errorf("provider %s does not support pre-generated short codes", ProviderMySQLAutoIncrement)
	}

	generator, err := m.Get(m.provider)
	if err != nil {
		return "", err
	}

	return generator.NextCode(ctx)
}

func BuildCodeFromID(id uint64) string {
	return utils.EncodeBase62(id)
}
