package utils

import (
	"math/rand"
	"strings"
	"time"
)

var charset = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var seedSet = false

func HasPrefix(s string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if s == "" || prefix == "" {
			continue
		}
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}

func RandomStringFrom(n int, chars []rune) string {
	if n <= 0 || len(chars) <= 0 {
		return ""
	}
	if !seedSet {
		rand.Seed(time.Now().UnixNano())
		seedSet = true
	}
	b := make([]rune, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func RandomString(n int) string {
	return RandomStringFrom(n, charset)
}