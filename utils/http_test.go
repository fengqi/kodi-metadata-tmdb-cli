package utils

import (
	"testing"
)

func TestStringMapToQuery(t *testing.T) {
	tests := []struct {
		name string
		m    map[string]string
		want string
	}{
		{
			name: "empty map",
			m:    map[string]string{},
			want: "",
		},
		{
			name: "single key-value pair",
			m:    map[string]string{"key": "value"},
			want: "key=value",
		},
		{
			name: "multiple key-value pairs",
			m:    map[string]string{"key1": "value1", "key2": "value2"},
			want: "key1=value1&key2=value2",
		},
		{
			name: "escaped values",
			m:    map[string]string{"key": "value with spaces"},
			want: "key=value+with+spaces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StringMapToQuery(tt.m)
			if got != tt.want {
				t.Errorf("StringMapToQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
