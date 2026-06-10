package codestrategy

import "testing"

func TestNewSnowflakeGeneratorRejectsInvalidWorkerID(t *testing.T) {
	_, err := NewSnowflakeGenerator(-1)
	if err == nil {
		t.Fatal("NewSnowflakeGenerator() expected error for invalid worker id")
	}
}

func TestSnowflakeGeneratorRequiresDAO(t *testing.T) {
	generator, err := NewSnowflakeGenerator(1)
	if err != nil {
		t.Fatalf("NewSnowflakeGenerator() unexpected error: %v", err)
	}

	code, err := generator.NextCode(t.Context())
	if err != nil {
		t.Fatalf("NextCode() unexpected error: %v", err)
	}
	if code == "" {
		t.Fatal("NextCode() expected non-empty short code")
	}
}
