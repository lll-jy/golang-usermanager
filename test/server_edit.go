package test

import (
	"database/sql"
	"fmt"
	"net/http"
	"testing"

	"git.garena.com/jiayu.li/entry-task/cmd/handlers"
	"git.garena.com/jiayu.li/entry-task/cmd/protocol"
)

func test_valid_edit(t *testing.T, db *sql.DB, i int) {
	name := fmt.Sprintf("user%d", i)
	pass := fmt.Sprintf("pass%d%d", i*2, i*2)
	nickname := fmt.Sprintf("nick%d", i)
	nicknew := fmt.Sprintf("mick%d", i)
	response, request := formSetup(fmt.Sprintf("nickname=%s", nicknew), t, db, "/edit")
	cookieString := handlers.SetSessionInfo(
		&protocol.User{
			Name:     name,
			Password: pass,
			Nickname: nickname,
		},
		&protocol.User{
			Name:     name,
			Password: pass,
			Nickname: nickname,
		},
		handlers.InfoErr{},
		"assets/placeholder.jpeg",
	)
	updateCookie(cookieString, response, request)
	http.HandlerFunc(makeHandler(db, handlers.EditHandler)).ServeHTTP(response, request)
	user := &protocol.User{}
	flag := protocol.IsExistingUsername(db, name, user)
	if !flag {
		t.Errorf("Wrongly deleted/updated primary key of %s.", name)
	} else if user.Nickname != nicknew {
		t.Errorf("Update of %s failed.", name)
	}
}
