package main

import "regexp"

func isValidUsername(username string) bool { // TODO
	// var validUsername = regexp.MustCompile("^u([a-z]+)$")
	var validUsername = regexp.MustCompile("^[a-zA-Z](([a-zA-Z0-9]|-|_){3})(([a-zA-Z0-9]|-|_){0,16})$")
	// start with letter, at least 4 characters, at most 20 characters, string with alphanumerical or '-' or '_'
	return validUsername.MatchString(username)
}

func isValidPassword(password string) bool { // TODO
	var validPassword = regexp.MustCompile("^uu([a-z]+)$")
	return validPassword.MatchString(password)
}

func isExistingUsername(username string) bool { // TODO
	var validUsername = regexp.MustCompile("^uu([a-z]+)$")
	return validUsername.MatchString(username)
}

func isCorrectPassword(username string, password string) bool { // TODO
	return username == password
}
