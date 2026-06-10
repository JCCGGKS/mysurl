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

func (g *stubGenerator) NextCode(context.Context) (string, error) {
	if g.err != nil {
		return "", g.err
	}

	return g.code, nil
}

func TestCodeManagerRegisterOverwrite(t *testing.T) {
	manager := NewCodeManager(ProviderRedisIncr)
	manager.Register(&stubGenerator{provider: ProviderRedisIncr, code: "old"})
	manager.Register(&stubGenerator{provider: ProviderRedisIncr, code: "new"})

	code, err := manager.NextCode(context.Background())
	if err != nil {
		t.Fatalf("NextCode() unexpected error: %v", err)
	}
	if code != "new" {
		t.Fatalf("NextCode() = %q, want %q", code, "new")
	}
}

func TestCodeManagerGetMissingProvider(t *testing.T) {
	manager := NewCodeManager(ProviderMySQLAutoIncrement)

	if _, err := manager.Get(ProviderRedisIncr); err == nil {
		t.Fatal("Get() expected error for missing provider")
	}
}

func TestCodeManagerNextCodePropagatesError(t *testing.T) {
	manager := NewCodeManager(ProviderRedisIncr)
	wantErr := errors.New("boom")
	manager.Register(&stubGenerator{provider: ProviderRedisIncr, err: wantErr})

	_, err := manager.NextCode(context.Background())
	if !errors.Is(err, wantErr) {
		t.Fatalf("NextCode() error = %v, want %v", err, wantErr)
	}
}

func TestCodeManagerRejectsMySQLAutoIncrementNextCode(t *testing.T) {
	manager := NewCodeManager(ProviderMySQLAutoIncrement)

	if _, err := manager.NextCode(context.Background()); err == nil {
		t.Fatal("NextCode() expected error for mysql auto increment provider")
	}
}

func TestBuildCodeFromID(t *testing.T) {
	if code := BuildCodeFromID(62); code != "10" {
		t.Fatalf("BuildCodeFromID() = %q, want %q", code, "10")
	}
}
