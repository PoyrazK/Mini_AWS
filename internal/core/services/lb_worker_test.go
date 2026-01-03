package services

import (
	"testing"
)

func TestParsePorts(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]int
	}{
		{"8080:80", map[string]int{"8080": 80}},
		{"3000:3000,5000:5000", map[string]int{"3000": 3000, "5000": 5000}},
		{"", map[string]int{}},
	}

	for _, tt := range tests {
		result := parsePorts(tt.input)
		for k, v := range tt.expected {
			if result[k] != v {
				t.Errorf("parsePorts(%q): expected %d for key %s, got %d", tt.input, v, k, result[k])
			}
		}
	}
}

func TestSplitCommas(t *testing.T) {
	result := splitCommas("a,b,c")
	if len(result) != 3 || result[0] != "a" || result[1] != "b" || result[2] != "c" {
		t.Errorf("splitCommas failed: %v", result)
	}
}

func TestSplitColons(t *testing.T) {
	result := splitColons("host:80")
	if len(result) != 2 || result[0] != "host" || result[1] != "80" {
		t.Errorf("splitColons failed: %v", result)
	}
}
