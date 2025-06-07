package encoder

import (
	"testing"
)

func TestEncodeBase62(t *testing.T) {
	tests := []struct {
		id       uint64
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{9, "9"},
		{10, "a"},
		{35, "z"},
		{36, "A"},
		{61, "Z"},
		{62, "10"},
		{63, "11"},
		{124, "20"},
		{3843, "ZZ"},
		{238327, "ZZZ"},
		{1234567890, "1ly7vk"},
	}

	for _, tt := range tests {
		got := EncodeBase62(tt.id)
		if got != tt.expected {
			t.Errorf("EncodeBase62(%d) = %q; want %q", tt.id, got, tt.expected)
		}
	}
}

func TestDecodeBase62(t *testing.T) {
	tests := []struct {
		id       string
		expected uint64
	}{
		{"0", 0},
		{"1", 1},
		{"9", 9},
		{"a", 10},
		{"z", 35},
		{"A", 36},
		{"Z", 61},
		{"10", 62},
		{"11", 63},
		{"20", 124},
		{"ZZ", 3843},
		{"ZZZ", 238327},
		{"1LY7VK", 1624950792},
	}

	for _, tt := range tests {
		got := DecodeBase62(tt.id)
		if got != tt.expected {
			t.Errorf("DecodeBase62(%s) = %d; want %d", tt.id, got, tt.expected)
		}
	}
}
