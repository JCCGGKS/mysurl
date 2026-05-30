package utils

import "testing"

func TestEncodeBase62(t *testing.T) {
	tests := []struct {
		name  string
		input uint64
		want  string
	}{
		{name: "zero", input: 0, want: "0"},
		{name: "one", input: 1, want: "1"},
		{name: "sixty one", input: 61, want: "Z"},
		{name: "sixty two", input: 62, want: "10"},
		{name: "three eight four three", input: 3843, want: "ZZ"},
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
		{name: "zero", input: "0", want: 0},
		{name: "mixed", input: "1zZ", want: 1*62*62 + 35*62 + 61},
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
