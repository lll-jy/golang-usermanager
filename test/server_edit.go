package test

import (
	"database/sql"
	"fmt"
	"net/http"
	"testing"

	"git.garena.com/jiayu.li/entry-task/cmd/handlers"
	"git.garena.com/jiayu.li/entry-task/cmd/paths"
	"git.garena.com/jiayu.li/entry-task/cmd/protocol"
)

func editExecute(t *testing.T, db *sql.DB, name, pass, photo, tempPhoto, nick, nicknew string) (http.Header, *protocol.User) {
	response, request := formSetup(fmt.Sprintf("nickname=%s", nicknew), t, db, "/edit")
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
		handlers.InfoErr{},
		tempPhoto,
	)
	updateCookie(cookieString, response, request)
	http.HandlerFunc(makeHandler(db, handlers.EditHandler)).ServeHTTP(response, request)
	return response.Header(), user
}

func test_valid_edit_photo_uploaded(t *testing.T, db *sql.DB, i int) {
	name := fmt.Sprintf("user%d", i)
	pass := fmt.Sprintf("pass%d%d", i*2, i*2)
	nick := fmt.Sprintf("nick%d", i)
	tempPhoto := fmt.Sprintf("%s/user%s.jpeg", paths.TempPath, name)
	photo := test_upload(t, db, i)
	editExecute(t, db, name, pass, photo, tempPhoto, nick, nick)
	user := &protocol.User{}
	//user = getUser(header)
	//t.Errorf("See here: %v", user)
	flag := protocol.IsExistingUsername(db, name, user)
	//t.Errorf("See there: %v", user)
	if !flag {
		t.Errorf("Wrongly deleted user from database.")
	} else if user.PhotoUrl != photo {
		t.Errorf("The photo uploaded wrongly.")
	}
}

func test_valid_edit_nickname(t *testing.T, db *sql.DB, i int) {
	name := fmt.Sprintf("user%d", i)
	pass := fmt.Sprintf("pass%d%d", i*2, i*2)
	nickname := fmt.Sprintf("nick%d", i)
	nicknew := fmt.Sprintf("mick%d", i)
	_, user := editExecute(t, db, name, pass, paths.PlaceholderPath, paths.PlaceholderPath, nickname, nicknew)
	flag := protocol.IsExistingUsername(db, name, user)
	if !flag {
		t.Errorf("Wrongly deleted/updated primary key of %s.", name)
	} else if user.Nickname != nicknew {
		t.Errorf("Update of %s failed.", name)
	}
}
