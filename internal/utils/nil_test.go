package utils

import "testing"

func TestIsNil(t *testing.T) {
	var nilPointer *int
	var nilSlice []string
	var nilMap map[string]string
	var nilFunc func()
	var nilChan chan int
	var nilInterface any

	value := 1
	nonNilPointer := &value
	nonNilSlice := []string{}
	nonNilMap := map[string]string{}
	nonNilChan := make(chan int)

	tests := []struct {
		name  string
		input any
		want  bool
	}{
		{name: "nil literal", input: nil, want: true},
		{name: "nil pointer", input: nilPointer, want: true},
		{name: "nil slice", input: nilSlice, want: true},
		{name: "nil map", input: nilMap, want: true},
		{name: "nil func", input: nilFunc, want: true},
		{name: "nil chan", input: nilChan, want: true},
		{name: "nil interface", input: nilInterface, want: true},
		{name: "non nil pointer", input: nonNilPointer, want: false},
		{name: "non nil slice", input: nonNilSlice, want: false},
		{name: "non nil map", input: nonNilMap, want: false},
		{name: "non nil chan", input: nonNilChan, want: false},
		{name: "zero int", input: 0, want: false},
		{name: "empty string", input: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNil(tt.input); got != tt.want {
				t.Fatalf("IsNil(%T) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestSuccessTreatsTypedNilAsNil(t *testing.T) {
	var nilPointer *int

	resp := Success(nilPointer)
	if resp.Data != nil {
		t.Fatalf("Success() data = %#v, want nil", resp.Data)
	}
}
