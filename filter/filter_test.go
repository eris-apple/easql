package filter_test

import (
	"github.com/eris-apple/easql/filter"
	"testing"
)

func TestFilter_String(t *testing.T) {
	tests := []struct {
		name     string
		input    filter.Filter
		expected string
	}{
		{"Empty struct", filter.Filter{}, "limit=0 offset=0 order="},
		{"Only limit", filter.Filter{Limit: 2}, "limit=2 offset=0 order="},
		{"Limit and offset", filter.Filter{Limit: 2, Offset: 1}, "limit=2 offset=1 order="},
		{"Full", filter.Filter{Limit: 2, Offset: 1, Order: "asc"}, "limit=2 offset=1 order=asc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.String()
			if result != tt.expected {
				t.Errorf("expected: %v, got: %v", tt.expected, result)
			}
		})
	}
}
