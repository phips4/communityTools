package utils

import (
	"strings"
)

func HasPrefix(s string,prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}