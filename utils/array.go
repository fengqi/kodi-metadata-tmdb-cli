package utils

import "strings"

// ArrayCompare 比较数组是否一致, 非严格模式下顺序可不一致
func ArrayCompare[T comparable](arr1, arr2 []T, strict bool) bool {
	if len(arr1) != len(arr2) {
		return false
	}

	for k, v := range arr1 {
		if strict {
			if arr2[k] != v {
				return false
			}
		} else {
			if !InArray(arr2, v) {
				return false
			}
		}
	}

	return true
}

// InArray 在数组内
func InArray[T comparable](arr []T, noodle T) bool {
	for _, item := range arr {
		if item == noodle {
			return true
		}
	}
	return false
}

// HasArrayPrefix 数组中的字符以prefix开头
func HasArrayPrefix(prefix string, arr []string) bool {
	if arr == nil {
		return false
	}

	for _, s := range arr {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}

// ArrayToUpper 字符串数组转为大小
func ArrayToUpper(arr []string) []string {
	for k, v := range arr {
		arr[k] = strings.ToUpper(v)
	}
	return arr
}

// SplitWith 以指定的多个字符串切割string，且跳过部分特殊字符
// 灵感来源：strings.FieldsFunc()
func SplitWith(s string, sep []string, exclude []string) []string {
	type span struct {
		start int
		end   int
	}

	spans := make([]span, 0, 32)
	start := -1
	for end, item := range s {
		flag := false
		if InArray(sep, string(item)) {
			flag = true
		}

		if flag {
			if start >= 0 {
				spans = append(spans, span{start, end})
				start = ^start
			}
		} else {
			if start < 0 {
				start = end
			}
		}
	}

	if start >= 0 {
		spans = append(spans, span{start, len(s)})
	}

	exclude = ArrayToUpper(exclude)
	a := make([]string, len(spans))
	flag := false
	jump := 0
	for i, span := range spans {
		if span.start < jump {
			continue
		}

		a[i] = s[span.start:span.end]

		if HasArrayPrefix(strings.ToUpper(a[i]), exclude) {
			flag = true
			continue
		}

		if flag && exclude != nil {
			for _, item := range exclude {
				guest := s[spans[i-1].start : spans[i-1].start+len(item)]
				if InArray(exclude, strings.ToUpper(guest)) {
					a[i] = guest
					a[i-1] = ""
					jump = spans[i-1].start + len(item)
					flag = false
					break
				}
			}
		}
	}

	k := 0
	for _, v := range a {
		if v != "" {
			a[k] = v
			k++
		}
	}

	return a[:k]
}
