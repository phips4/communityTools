package logic

import "regexp"

const (
	MaxIdStringLength   = 64
	MaxLongStringLength = 1024
	DeleteTokenLength   = 32
)

func IdLengthValidation(str string) bool {
	l := len(str)
	if l > MaxIdStringLength || l < 1 {
		return false
	}
	return true
}

func ValidateID(id string) bool {
	if !IdLengthValidation(id) {
		return false
	}
	ok, _ := regexp.Match("^[a-zA-Z0-9_]*$", []byte(id))
	return ok
}