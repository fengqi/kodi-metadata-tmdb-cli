package utils

import (
	"strconv"
	"strings"
)

// EndsWith 字符以xx结尾
func EndsWith(str, subStr string) bool {
	index := strings.LastIndex(str, subStr)

	return index > 0 && str[index:] == subStr
}

func StrToInt(str string) int {
	if n, err := strconv.Atoi(str); err == nil {
		return n
	}
	return 0
}
