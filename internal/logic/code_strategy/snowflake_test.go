package codestrategy

import (
	"context"
	"testing"
)

func TestNewSnowflakeGeneratorRejectsInvalidWorkerID(t *testing.T) {
	_, err := NewSnowflakeGenerator(-1, nil)
	if err == nil {
		t.Fatal("NewSnowflakeGenerator() expected error for invalid worker id")
	}
}

func TestSnowflakeGeneratorRequiresDAO(t *testing.T) {
	generator, err := NewSnowflakeGenerator(1, nil)
	if err != nil {
		t.Fatalf("NewSnowflakeGenerator() unexpected error: %v", err)
	}

	_, err = generator.NextCode(context.Background(), nil, "https://example.com", "hash")
	if err == nil {
		t.Fatal("NextCode() expected error when dao is nil")
	}
}
