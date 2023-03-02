package utils

import "testing"

func TestArrayCompare(t *testing.T) {
	type strCase struct {
		arr1   []string
		arr2   []string
		want   bool
		strict bool
	}

	cases := []strCase{
		{
			arr1:   []string{"a", "b", "c"},
			arr2:   []string{"a", "b", "c"},
			want:   true,
			strict: false,
		},
		{
			arr1:   []string{"b", "a", "c"},
			arr2:   []string{"c", "b", "a"},
			want:   true,
			strict: false,
		},
		{
			arr1:   []string{"b", "a", "c"},
			arr2:   []string{"c", "b", "a"},
			want:   false,
			strict: true,
		},
		{
			arr1:   []string{"a", "b", "c"},
			arr2:   []string{"a", "b", "d"},
			want:   false,
			strict: true,
		},
	}
	for _, item := range cases {
		give := ArrayCompare(item.arr1, item.arr2, item.strict)
		if give != item.want {
			t.Errorf("ArrayCompare(%v, %v, %v) give: %v, want: %v", item.arr1, item.arr2, item.strict, give, item.want)
		}
	}
}

func TestSplitWith(t *testing.T) {
	type testCase struct {
		str     string
		sep     []string
		split   []string
		exclude []string
	}

	cases := []testCase{
		{
			str:     "a.b-c.d",
			sep:     []string{".", "-"},
			split:   []string{"a", "b-c", "d"},
			exclude: []string{"b-c"},
		},
		{
			str:     "a.b-c@d",
			sep:     []string{".", "-", "@"},
			split:   []string{"a", "b", "c", "d"},
			exclude: []string{},
		},
	}
	for _, item := range cases {
		give := SplitWith(item.str, item.sep, item.exclude)
		if !ArrayCompare(item.split, give, true) {
			t.Errorf("SplitWith(%s, %v) give: %v, want: %v", item.str, item.sep, give, item.split)
		}
	}
}
