package utils

import (
	"path/filepath"
	"strings"
)

// NormalizePath 标准化路径用于比较
func NormalizePath(path string) string {
	cleaned := filepath.ToSlash(filepath.Clean(strings.TrimSpace(path)))
	return strings.TrimRight(cleaned, "/")
}

// PathEqual 比较两个路径是否相等
func PathEqual(pathA, pathB string) bool {
	return NormalizePath(pathA) == NormalizePath(pathB)
}
