package logic

import (
	"log"
	"regexp"
	"strconv"
)

const maxStringLength = 32

func ValidateID(id string) bool {

	if !DefaultValidation(id) {
		return false
	}
	ok, _ := regexp.Match("^[a-zA-Z0-9_]*$", []byte(id))
	return ok
}

func ValidatePostParams(title, desc, cookieCheck, multiOptions string, options []string) bool {

	if !DefaultValidation(title) || !DefaultValidation(desc) {
		return false
	}

	if _, err := strconv.ParseBool(cookieCheck); err != nil {
		log.Println("parse bool cookie")
		return false
	}

	if _, err := strconv.ParseBool(multiOptions); err != nil {
		log.Println("parse bool multiOptions")
		return false
	}

	if options == nil || len(options) < 2 {
		log.Println("options array too small")
		return false
	}

	for _, elem := range options {
		if !DefaultValidation(elem) {
			log.Println("element in array '" + elem + "' is not default valid")
			return false
		}
	}

	return true
}

func DefaultValidation(str string) bool {
	l := len(str)
	if l > maxStringLength || l < 1 {
		return false
	}
	return true
}
