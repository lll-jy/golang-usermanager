// Package protocol handles helper functions concerning user data structure.
package protocol

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"regexp"

	"git.garena.com/jiayu.li/entry-task/cmd/paths"
)

// IsValidUsername checks whether a given username start with letter, at least 4 characters, at most 20 characters,
// string with alphanumerical or '-' or '_'.
func IsValidUsername(username string) bool {
	var validUsername = regexp.MustCompile("^[a-zA-Z](([a-zA-Z0-9]|-|_){3})(([a-zA-Z0-9]|-|_){0,16})$")
	return validUsername.MatchString(username)
}

// IsValidPassword checks whether a given password is between 4-20 characters.
func IsValidPassword(password string) bool {
	length := len(password)
	return length <= 20 && length >= 4
}

// ConvertToString converts a given interface read from database to a string.
func ConvertToString(i interface{}) string {
	s := fmt.Sprintf("%s", i)
	if s == "%!s(<nil>)" {
		s = ""
	}
	return s
}

// IsExistingUsername checks whether the given username has duplicate in the database.
func IsExistingUsername(db *sql.DB, username string, user *User) bool {
	query, err := db.Prepare("SELECT password, photo, nickname FROM users WHERE username = ? AND username <> ?")
	if err != nil {
		log.Printf("Cannot parse query: %s", err.Error())
		return false
	}
	var p, pu, nn interface{}
	query.QueryRow(username, user.Name).Scan(&p, &pu, &nn)
	query.Close()
	user.Password = ConvertToString(p)
	user.PhotoUrl = ConvertToString(pu)
	user.Nickname = ConvertToString(nn)
	if user.Password != "" {
		if user.PhotoUrl == "" {
			user.PhotoUrl = paths.PlaceholderPath
		}
		return true
	} else {
		return false
	}
}

// IsCorrectPassword checks whether the password user provides hashes to the password stored in the database.
func IsCorrectPassword(userPass string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(userPass)) == nil
}
