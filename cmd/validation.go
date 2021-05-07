package main

import (
	"log"
	"regexp"

	"git.garena.com/jiayu.li/entry-task/cmd/protocol"
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

func isExistingUsername(username string, user *protocol.User) bool {
	query, err := db.Prepare("SELECT password, photo, nickname FROM users WHERE username = ? AND username <> ?")
	if err != nil {
		log.Printf("Cannot parse query: %s", err.Error())
		return false
	}
	defer query.Close()
	user.Password = ""
	query.QueryRow(username, user.Name).Scan(&user.Password, &user.PhotoUrl, &user.Nickname)
	if user.Password != "" {
		if user.PhotoUrl == "" {
			user.PhotoUrl = "assets/placeholder.jpeg" // EXTEND: maybe some cloud space
		}
		return true
	} else {
		return false
	}
}

func isCorrectPassword(userpass string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(userpass)) == nil
}
