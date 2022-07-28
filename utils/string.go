package utils

import "strings"

// EndsWith 字符以xx结尾
func EndsWith(str, subStr string) bool {
	index := strings.LastIndex(str, subStr)

	return index > 0 && str[index:] == subStr
}
