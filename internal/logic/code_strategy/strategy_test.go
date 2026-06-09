package codestrategy

import (
	"context"
	"errors"
	"testing"
)

type stubGenerator struct {
	provider string
	code     string
	err      error
}

func (g *stubGenerator) Provider() string {
	return g.provider
}

func (g *stubGenerator) NextCode(context.Context, *uint64, string, string) (string, error) {
	if g.err != nil {
		return "", g.err
	}

	return g.code, nil
}

func TestCodeManagerRegisterOverwrite(t *testing.T) {
	manager := NewCodeManager(ProviderRedisIncr)
	manager.Register(&stubGenerator{provider: ProviderRedisIncr, code: "old"})
	manager.Register(&stubGenerator{provider: ProviderRedisIncr, code: "new"})

	code, err := manager.GenerateShortCode(context.Background(), nil, "https://example.com", "hash")
	if err != nil {
		t.Fatalf("GenerateShortCode() unexpected error: %v", err)
	}
	if code != "new" {
		t.Fatalf("GenerateShortCode() = %q, want %q", code, "new")
	}
}

func TestCodeManagerGetMissingProvider(t *testing.T) {
	manager := NewCodeManager(ProviderMySQLAutoIncrement)

	if _, err := manager.Get(ProviderRedisIncr); err == nil {
		t.Fatal("Get() expected error for missing provider")
	}
}

func TestCodeManagerGenerateShortCodePropagatesError(t *testing.T) {
	manager := NewCodeManager(ProviderMySQLAutoIncrement)
	wantErr := errors.New("boom")
	manager.Register(&stubGenerator{provider: ProviderMySQLAutoIncrement, err: wantErr})

	_, err := manager.GenerateShortCode(context.Background(), nil, "https://example.com", "hash")
	if !errors.Is(err, wantErr) {
		t.Fatalf("GenerateShortCode() error = %v, want %v", err, wantErr)
	}
}
