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

type stubBatchGenerator struct {
	stubGenerator
	codes []string
}

func (g *stubBatchGenerator) NextCodes(context.Context, int) ([]string, error) {
	if g.err != nil {
		return nil, g.err
	}

	return g.codes, nil
}

func TestCodeManagerRegisterOverwrite(t *testing.T) {
	manager := NewCodeManager()
	manager.Register(&stubGenerator{provider: ProviderRedisIncr, code: "old"})
	manager.Register(&stubGenerator{provider: ProviderRedisIncr, code: "new"})

	code, err := manager.NextCode(context.Background(), ProviderRedisIncr)
	if err != nil {
		t.Fatalf("NextCode() unexpected error: %v", err)
	}
	if code != "new" {
		t.Fatalf("NextCode() = %q, want %q", code, "new")
	}
}

func TestCodeManagerGetMissingProvider(t *testing.T) {
	manager := NewCodeManager()

	if _, err := manager.Get(ProviderRedisIncr); err == nil {
		t.Fatal("Get() expected error for missing provider")
	}
}

func TestCodeManagerNextCodePropagatesError(t *testing.T) {
	manager := NewCodeManager()
	wantErr := errors.New("boom")
	manager.Register(&stubGenerator{provider: ProviderRedisIncr, err: wantErr})

	_, err := manager.NextCode(context.Background(), ProviderRedisIncr)
	if !errors.Is(err, wantErr) {
		t.Fatalf("NextCode() error = %v, want %v", err, wantErr)
	}
}

func TestCodeManagerRejectsMySQLAutoIncrementNextCode(t *testing.T) {
	manager := NewCodeManager()

	if _, err := manager.NextCode(context.Background(), ProviderMySQLAutoIncrement); err == nil {
		t.Fatal("NextCode() expected error for mysql auto increment provider")
	}
}

func TestCodeManagerNextCodesUsesBatchGenerator(t *testing.T) {
	manager := NewCodeManager()
	manager.Register(&stubBatchGenerator{
		stubGenerator: stubGenerator{provider: ProviderRedisIncr},
		codes:         []string{"a", "b", "c"},
	})

	codes, err := manager.NextCodes(context.Background(), ProviderRedisIncr, 3)
	if err != nil {
		t.Fatalf("NextCodes() unexpected error: %v", err)
	}
	if len(codes) != 3 || codes[0] != "a" || codes[2] != "c" {
		t.Fatalf("NextCodes() = %v, want [a b c]", codes)
	}
}

func TestCodeManagerNextCodesFallsBackToSingleCalls(t *testing.T) {
	manager := NewCodeManager()
	manager.Register(&stubGenerator{provider: ProviderSnowflake, code: "x"})

	codes, err := manager.NextCodes(context.Background(), ProviderSnowflake, 2)
	if err != nil {
		t.Fatalf("NextCodes() unexpected error: %v", err)
	}
	if len(codes) != 2 || codes[0] != "x" || codes[1] != "x" {
		t.Fatalf("NextCodes() = %v, want [x x]", codes)
	}
}

func TestBuildCodeFromID(t *testing.T) {
	if code := BuildCodeFromID(62); code != "RO" {
		t.Fatalf("BuildCodeFromID() = %q, want %q", code, "RO")
	}
}
