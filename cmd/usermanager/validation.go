package main

import (
	"fmt"
	"log"

	"git.garena.com/jiayu.li/entry-task/cmd/paths"
	"git.garena.com/jiayu.li/entry-task/cmd/protocol"
	"golang.org/x/crypto/bcrypt"
)

func convertToString(i interface{}) string {
	s := fmt.Sprintf("%s", i)
	if s == "%!s(<nil>)" {
		s = ""
	}
	return s
}

func isExistingUsername(username string, user *protocol.User) bool {
	query, err := db.Prepare("SELECT password, photo, nickname FROM users WHERE username = ? AND username <> ?")
	if err != nil {
		log.Printf("Cannot parse query: %s", err.Error())
		return false
	}
	defer query.Close()
	var p, pu, nn interface{}
	query.QueryRow(username, user.Name).Scan(&p, &pu, &nn)
	user.Password = convertToString(p)
	user.PhotoUrl = convertToString(pu)
	user.Nickname = convertToString(nn)
	if user.Password != "" {
		if user.PhotoUrl == "" {
			user.PhotoUrl = paths.PlaceholderPath // EXTEND: maybe some cloud space
		}
		return true
	} else {
		return false
	}
}

func isCorrectPassword(userpass string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(userpass)) == nil
}
