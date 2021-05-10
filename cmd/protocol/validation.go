package protocol

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"

	"git.garena.com/jiayu.li/entry-task/cmd/paths"
	"golang.org/x/crypto/bcrypt"
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

func IsExistingUsername(db *sql.DB, username string, user *User) bool {
	query, err := db.Prepare("SELECT password, photo, nickname FROM users WHERE username = ? AND username <> ?")
	if err != nil {
		log.Printf("Cannot parse query: %s", err.Error())
		return false
	}
	defer query.Close()
	var p, pu, nn interface{}
	query.QueryRow(username, user.Name).Scan(&p, &pu, &nn)
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

func IsCorrectPassword(userpass string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(userpass)) == nil
}
