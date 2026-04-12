package utils

import "testing"

func TestNormalizePath(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "trim and trailing slash",
			input:    "  /data/shows/  ",
			expected: "/data/shows",
		},
		{
			name:     "mixed separators on windows style path",
			input:    `C:\data\shows\Season 1\`,
			expected: "C:/data/shows/Season 1",
		},
		{
			name:     "clean dot and dotdot",
			input:    "/data//shows/./A/../B/",
			expected: "/data/shows/B",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := NormalizePath(tc.input)
			if got != tc.expected {
				t.Errorf("NormalizePath(%q) = %q; expected %q", tc.input, got, tc.expected)
			}
		})
	}
}

func TestPathEqual(t *testing.T) {
	testCases := []struct {
		name  string
		a     string
		b     string
		equal bool
	}{
		{
			name:  "slash style equivalent",
			a:     `C:\data\shows\A`,
			b:     "C:/data/shows/A",
			equal: true,
		},
		{
			name:  "trailing slash equivalent",
			a:     "/data/shows/A/",
			b:     "/data/shows/A",
			equal: true,
		},
		{
			name:  "cleaned path equivalent",
			a:     "/data/shows/../shows/A",
			b:     "/data/shows/A",
			equal: true,
		},
		{
			name:  "different path",
			a:     "/data/shows/A",
			b:     "/data/shows/B",
			equal: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := PathEqual(tc.a, tc.b)
			if got != tc.equal {
				t.Errorf("PathEqual(%q, %q) = %v; expected %v", tc.a, tc.b, got, tc.equal)
			}
		})
	}
}
