package utils

import (
	"net/url"
	"strings"
)

func StringMapToQuery(m map[string]string) string {
	if len(m) == 0 {
		return ""
	}

	var s strings.Builder
	for k, v := range m {
		s.WriteString(k + "=" + url.QueryEscape(v) + "&")
	}

	return strings.TrimRight(s.String(), "&")
}
