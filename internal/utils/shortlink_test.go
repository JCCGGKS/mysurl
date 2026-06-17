package utils

import (
	"strings"
	"testing"
)

func TestShuffleBase62Alphabet(t *testing.T) {
	got := shuffleBase62Alphabet(42)
	if got != "juZ0ILhQRpNVa7PE9WrGzge3KBtXsqmD6fYx1o8OdMHFTwSv5ibCykUJl2cA4n" {
		t.Fatalf("shuffleBase62Alphabet(42) = %q", got)
	}

	if len(got) != len(base62AlphabetSource) {
		t.Fatalf("shuffleBase62Alphabet(42) length = %d, want %d", len(got), len(base62AlphabetSource))
	}

	for _, ch := range base62AlphabetSource {
		if strings.Count(got, string(ch)) != 1 {
			t.Fatalf("shuffleBase62Alphabet(42) invalid count for %q", ch)
		}
	}
}

func TestEncodeBase62(t *testing.T) {
	tests := []struct {
		name  string
		input uint64
		want  string
	}{
		{name: "zero", input: 0, want: "j"},
		{name: "one", input: 1, want: "u"},
		{name: "sixty one", input: 61, want: "n"},
		{name: "sixty two", input: 62, want: "uj"},
		{name: "three eight four three", input: 3843, want: "nn"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodeBase62(tt.input); got != tt.want {
				t.Fatalf("EncodeBase62(%d) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestDecodeBase62(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    uint64
		wantErr bool
	}{
		{name: "zero", input: "j", want: 0},
		{name: "mixed", input: "uxn", want: 1*62*62 + 35*62 + 61},
		{name: "invalid", input: "!", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeBase62(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("DecodeBase62(%q) expected error", tt.input)
				}

				return
			}

			if err != nil {
				t.Fatalf("DecodeBase62(%q) unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("DecodeBase62(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}
