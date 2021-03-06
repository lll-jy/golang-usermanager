package server_helpers

import (
	"database/sql"
	"fmt"
	"net/http"
	"testing"

	"git.garena.com/jiayu.li/entry-task/cmd/handlers"
	"git.garena.com/jiayu.li/entry-task/cmd/paths"
	"git.garena.com/jiayu.li/entry-task/cmd/protocol"
)

func editExecute(t *testing.T, db *sql.DB, name, pass, photo, tempPhoto,
	nick, nickNew string) (header http.Header, usr *protocol.User) {
	response, request := formSetup(fmt.Sprintf("nickname=%s", nickNew), t, db, "/edit")
	user := &protocol.User{}
	protocol.IsExistingUsername(db, name, user)
	cookieString := handlers.SetSessionInfo(
		&protocol.User{
			Name:     name,
			Password: user.Password,
			PhotoUrl: user.PhotoUrl,
			Nickname: nick,
		},
		&protocol.User{
			Name:     name,
			Password: pass,
			PhotoUrl: photo,
			Nickname: nick,
		},
		&handlers.InfoErr{},
		tempPhoto,
	)
	UpdateCookie(cookieString, response, request)
	http.HandlerFunc(makeHandler(db, handlers.EditHandler)).ServeHTTP(response, request)
	return response.Header(), user
}

func ValidEditPhoto(t *testing.T, db *sql.DB, i int) string {
	name := fmt.Sprintf("user%d", i)
	pass := fmt.Sprintf("pass%d%d", i*2, i*2)
	nick := fmt.Sprintf("nick%d", i)
	tempPhoto := fmt.Sprintf("%s/user%s.jpeg", paths.TempPath, name)
	photo := Upload(t, db, i)
	editExecute(t, db, name, pass, photo, tempPhoto, nick, nick)
	user := &protocol.User{}
	flag := protocol.IsExistingUsername(db, name, user)
	if !flag {
		t.Errorf("Wrongly deleted user from database.")
	} else if user.PhotoUrl != photo {
		t.Errorf("The photo uploaded wrongly.")
	}
	return photo
}

func ValidEditNickname(t *testing.T, db *sql.DB, i int) {
	name := fmt.Sprintf("user%d", i)
	pass := fmt.Sprintf("pass%d%d", i*2, i*2)
	nickname := fmt.Sprintf("nick%d", i)
	nickNew := fmt.Sprintf("mick%d", i)
	_, user := editExecute(t, db, name, pass, paths.PlaceholderPath, paths.PlaceholderPath, nickname, nickNew)
	flag := protocol.IsExistingUsername(db, name, user)
	if !flag {
		t.Errorf("Wrongly deleted/updated primary key of %s.", name)
	} else if user.Nickname != nickNew {
		t.Errorf("Update of %s failed.", name)
	}
}
