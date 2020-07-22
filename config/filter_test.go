package config_test

import (
	"testing"

	"github.com/rebel-l/mp3sync/config"
)

func TestFilter_Contains(t *testing.T) {
	testCases := []struct {
		name     string
		filter   config.Filter
		content  string
		expected bool
	}{
		{
			name:     "match exact",
			filter:   config.Filter{"exact"},
			content:  "exact",
			expected: true,
		},
		{
			name:     "first match exact",
			filter:   config.Filter{"exact", "contains"},
			content:  "exact",
			expected: true,
		},
		{
			name:     "second match exact",
			filter:   config.Filter{"contains", "exact"},
			content:  "exact",
			expected: true,
		},
		{
			name:     "match not",
			filter:   config.Filter{"chains"},
			content:  "contains",
			expected: false,
		},
		{
			name:     "contains",
			filter:   config.Filter{"contains"},
			content:  "exact contains this",
			expected: true,
		},
		{
			name:     "first contains",
			filter:   config.Filter{"contains", "something"},
			content:  "exact contains this",
			expected: true,
		},
		{
			name:     "second contains",
			filter:   config.Filter{"something", "contains"},
			content:  "exact contains this",
			expected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := testCase.filter.Contains(testCase.content)

			if testCase.expected != got {
				t.Errorf("expected %t but got %t", testCase.expected, got)
			}
		})
	}
}
