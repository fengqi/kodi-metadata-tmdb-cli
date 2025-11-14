package utils

import (
	"testing"
)

func TestStringMapToQuery(t *testing.T) {
	testCases := []struct {
		name     string
		input    map[string]string
		expected string
	}{
		{
			name:     "empty map",
			input:    map[string]string{},
			expected: "",
		},
		{
			name:     "single key-value pair",
			input:    map[string]string{"key1": "value1"},
			expected: "key1=value1",
		},
		{
			name:     "multiple key-value pairs",
			input:    map[string]string{"key1": "value1", "key2": "value2"},
			expected: "key1=value1&key2=value2",
		},
		{
			name:     "values with special characters",
			input:    map[string]string{"key1": "val&ue", "key2": "v@lue"},
			expected: "key1=val%26ue&key2=v%40lue",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := StringMapToQuery(tc.input)
			if result != tc.expected {
				t.Errorf("StringMapToQuery(%v) = %v; expected %v", tc.input, result, tc.expected)
			}
		})
	}
}
