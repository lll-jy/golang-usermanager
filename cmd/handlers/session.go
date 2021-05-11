package handlers

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
	NameErr       string
	PasswordErr       string
	PasswordRepeatErr string
}

type PageInfo struct {
	User         *protocol.User
	TempUser     *protocol.User
	InfoErr      InfoErr
	DisplayName  string
	Action       string
	Title        string
	CancelAction string
	Photo        string
}

func generatePageInfo(user, tempUser, nameErr, passErr, repeatPassErr, photo string) (info PageInfo) {
	u := &protocol.User{}
	if err := proto.Unmarshal([]uint8(user), u); err != nil {
		log.Printf("Error: wrong format! %s cannot be parsed as a user.", user)
	}
	tu := &protocol.User{}
	if err := proto.Unmarshal([]uint8(tempUser), tu); err != nil {
		log.Printf("Error: wrong format! %s cannot be parsed as a user (temp user)", tempUser)
	}
	ie := InfoErr{
		NameErr:       nameErr,
		PasswordErr:       passErr,
		PasswordRepeatErr: repeatPassErr,
	}
	return PageInfo{
		User:     u,
		TempUser: tu,
		InfoErr:  ie,
		Photo:    photo,
	}
}

func GetPageInfo(r *http.Request) (info PageInfo) {
	var user string
	var tempUser string
	var nameErr string
	var passErr string
	var repeatPassErr string
	var photo string
	if cookie, err := r.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			user = cookieValue["user"]
			tempUser = cookieValue["temp"]
			nameErr = cookieValue["nameErr"]
			passErr = cookieValue["passErr"]
			repeatPassErr = cookieValue["repeatPassErr"]
			photo = cookieValue["photo"]
		}
	}
	return generatePageInfo(user, tempUser, nameErr, passErr, repeatPassErr, photo)
}

func SetSessionInfo(u *protocol.User, tu *protocol.User, ie InfoErr, photo string) string {
	user, err := proto.Marshal(u)
	if err != nil {
		log.Printf("Error: wrong format! %s cannot be parsed as a user.", u)
	}
	tempUser, err := proto.Marshal(tu)
	if err != nil {
		log.Printf("Error: wrong format! %s cannot be parsed as a user (temp user).", tu)
	}
	value := map[string]string{
		"user":          string(user),
		"temp":          string(tempUser),
		"nameErr":       ie.NameErr,
		"passErr":       ie.PasswordErr,
		"repeatPassErr": ie.PasswordRepeatErr,
		"photo":         photo,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		return encoded
	} else {
		log.Printf("Session cannot be parsed and encoded.")
		return ""
	}
}

func setSession(u *protocol.User, tu *protocol.User, ie InfoErr, photo string, w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:  "session",
		Value: SetSessionInfo(u, tu, ie, photo),
		Path:  "/",
	}
	http.SetCookie(w, cookie)
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
