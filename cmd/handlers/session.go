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
	}
}

// InfoErr keeps the error messages of the page.
type InfoErr struct {
	NameErr           string
	PasswordErr       string
	PasswordRepeatErr string
}

// PageInfo keeps all the information that cookie may need to know.
type PageInfo struct {
	User         *protocol.User
	TempUser     *protocol.User
	InfoErr      *InfoErr
	DisplayName  string
	Action       string
	Title        string
	CancelAction string
	Photo        string
}

// generatePageInfo constructs the PageInfo type struct with given information.
func generatePageInfo(user, tempUser, nameErr, passErr, repeatPassErr, photo string) PageInfo {
	u := &protocol.User{}
	if err := proto.Unmarshal([]uint8(user), u); err != nil {
		log.Printf("Error: wrong format! %s cannot be parsed as a user: %s", user, err.Error())
	}
	tu := &protocol.User{}
	if err := proto.Unmarshal([]uint8(tempUser), tu); err != nil {
		log.Printf("Error: wrong format! %s cannot be parsed as a user (temp user): %s", tempUser, err.Error())
	}
	ie := &InfoErr{
		NameErr:           nameErr,
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

// GetPageInfo constructs the PageInfo struct passed in the cookie of given request.
func GetPageInfo(r *http.Request) PageInfo {
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

// SetSessionInfo generates the cookie encoded string with the given information.
func SetSessionInfo(u *protocol.User, tu *protocol.User, ie *InfoErr, photo string) string {
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

// setSession updates the cookie of the writer with the given information.
func setSession(u *protocol.User, tu *protocol.User, ie *InfoErr, photo string, w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:  "session",
		Value: SetSessionInfo(u, tu, ie, photo),
		Path:  "/",
	}
	http.SetCookie(w, cookie)
}

// clearSession clears session cookies.
func clearSession(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}
