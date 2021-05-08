package protocol

import (
	"fmt"
	"regexp"
)

func IsValidUsername(username string) bool {
	var validUsername = regexp.MustCompile("^[a-zA-Z](([a-zA-Z0-9]|-|_){3})(([a-zA-Z0-9]|-|_){0,16})$")
	// start with letter, at least 4 characters, at most 20 characters, string with alphanumerical or '-' or '_'
	return validUsername.MatchString(username)
}

func IsValidPassword(password string) bool {
	// 4-20 characters, case sensitive
	length := len(password)
	return length <= 20 && length >= 4
}

func ConvertToString(i interface{}) string {
	s := fmt.Sprintf("%s", i)
	if s == "%!s(<nil>)" {
		s = ""
	}
	return s
}
