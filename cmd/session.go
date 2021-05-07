package main

import (
	"net/http"

	"git.garena.com/jiayu.li/entry-task/cmd/protocol"
	"github.com/gorilla/securecookie"
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
	User         protocol.User
	InfoErr      InfoErr
	DisplayName  string
	Action       string
	Title        string
	CancelAction string
}

func getPageInfo(r *http.Request) (info PageInfo) {
	var username string
	var password string
	var photo string
	var nickname string
	var nameErr string
	var passErr string
	var repeatPassErr string
	if cookie, err := r.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			username = cookieValue["name"]
			password = cookieValue["pass"]
			photo = cookieValue["photo"]
			nickname = cookieValue["nickname"]
			nameErr = cookieValue["nameErr"]
			passErr = cookieValue["passErr"]
			repeatPassErr = cookieValue["repeatPassErr"]
		}
	}
	u := protocol.User{
		Name:     username,
		Password: password,
		PhotoUrl: photo,
		Nickname: nickname,
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

func setSession(u protocol.User, uie InfoErr, w http.ResponseWriter) {
	value := map[string]string{
		"name":          u.Name,
		"pass":          u.Password,
		"photo":         u.PhotoUrl,
		"nickname":      u.Nickname,
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
