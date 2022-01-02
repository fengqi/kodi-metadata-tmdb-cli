package utils

import (
	"net/url"
	"strings"
)

func StringMapToQuery(m map[string]string) string {
	if len(m) == 0 {
		return ""
	}

	s := ""
	for k, v := range m {
		s += k + "=" + url.QueryEscape(v) + "&"
	}

	return strings.TrimRight(s, "&")
}
