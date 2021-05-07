package main

import (
	"log"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func isValidUsername(username string) bool {
	var validUsername = regexp.MustCompile("^[a-zA-Z](([a-zA-Z0-9]|-|_){3})(([a-zA-Z0-9]|-|_){0,16})$")
	// start with letter, at least 4 characters, at most 20 characters, string with alphanumerical or '-' or '_'
	return validUsername.MatchString(username)
}

func isValidPassword(password string) bool {
	// 4-20 characters, case sensitive
	length := len(password)
	return length <= 20 && length >= 4
}

func isExistingUsername(username string, password *string) bool {
	query, err := db.Prepare("SELECT password FROM users WHERE username = ?")
	if err != nil {
		log.Printf("Cannot parse query: %s", err.Error())
		return false
	}
	defer query.Close()
	err = query.QueryRow(username).Scan(password)
	log.Printf("%s pass: %s", username, *password)
	return err == nil
}

func isCorrectPassword(userpass string, password string) bool {
	hashed, err := bcrypt.GenerateFromPassword([]byte(userpass), 3)
	if err != nil {
		log.Printf("Cannot hash user password %s", userpass)
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)) == nil
}
