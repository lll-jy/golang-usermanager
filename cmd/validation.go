package main

import "regexp"

func isValidUsername(username string) bool { // TODO
	var validUsername = regexp.MustCompile("^u([a-z]+)$")
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
