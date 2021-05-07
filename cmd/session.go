package main

import (
	"log"
	"net/http"

	"git.garena.com/jiayu.li/entry-task/cmd/protocol"
	"github.com/gorilla/securecookie"
	"google.golang.org/protobuf/proto"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32),
)

func createUser(name string, pass string) protocol.User {
	return protocol.User{
		Name:     name,
		Password: pass,
		PhotoUrl: "",
		Nickname: "",
	}
}

type InfoErr struct {
	UsernameErr       string
	PasswordErr       string
	PasswordRepeatErr string
}

type PageInfo struct {
	User         *protocol.User
	InfoErr      InfoErr
	DisplayName  string
	Action       string
	Title        string
	CancelAction string
}

func getPageInfo(r *http.Request) (info PageInfo) {
	var user string
	var nameErr string
	var passErr string
	var repeatPassErr string
	if cookie, err := r.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			user = cookieValue["user"]
			nameErr = cookieValue["nameErr"]
			passErr = cookieValue["passErr"]
			repeatPassErr = cookieValue["repeatPassErr"]
		}
	}
	u := &protocol.User{}
	if err := proto.Unmarshal([]uint8(user), u); err != nil {
		log.Printf("Error: wrong format! %s cannot be parsed as a user.", user)
	}
	ie := InfoErr{
		UsernameErr:       nameErr,
		PasswordErr:       passErr,
		PasswordRepeatErr: repeatPassErr,
	}
	return PageInfo{
		User:    u,
		InfoErr: ie,
	}
}

func setSession(u *protocol.User, uie InfoErr, w http.ResponseWriter) {
	user, err := proto.Marshal(u)
	if err != nil {
		log.Printf("Error: wrong format! %s cannot be parsed as a user.", u)
	}
	value := map[string]string{
		"user":          string(user),
		"nameErr":       uie.UsernameErr,
		"passErr":       uie.PasswordErr,
		"repeatPassErr": uie.PasswordRepeatErr,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}
}

func clearSession(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}
