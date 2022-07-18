package utils

import (
	"testing"
)

type endsWith struct {
	str    string
	subStr string
	with   bool
}

func TestEndsWith(t *testing.T) {
	cases := []endsWith{
		{"china", "na", true},
		{"china", "ch", false},
		{"string", "ring", true},
		{"string", "g", true},
		{"string", "s", false},
	}
	for _, item := range cases {
		give := EndsWith(item.str, item.subStr)
		if give != item.with {
			t.Errorf("EndsWith(%s, %s) give: %v, want: %v", item.str, item.subStr, give, item.with)
		}
	}
}
