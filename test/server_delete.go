package test

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.garena.com/jiayu.li/entry-task/cmd/handlers"
	"git.garena.com/jiayu.li/entry-task/cmd/protocol"
)

func deleteExecute(t *testing.T, db *sql.DB, name string, pass string) http.Header {
	cookieString := handlers.SetSessionInfo(
		&protocol.User{
			Name:     name,
			Password: pass,
		},
		&protocol.User{},
		handlers.InfoErr{},
		"",
	)
	response := httptest.NewRecorder()
	request := makeRequest(http.MethodPost, "/delete", t)
	updateCookie(cookieString, response, request)
	http.HandlerFunc(makeHandler(db, handlers.DeleteHandler)).ServeHTTP(response, request)
	return response.Header()
}

func test_valid_delete(t *testing.T, db *sql.DB, i int) {
	name := fmt.Sprintf("testuser%d", i)
	pass := fmt.Sprintf("pass%d%d", i*2, i*2)
	handlers.ExecuteQuery(db, "INSERT INTO users VALUES (?, ?, NULL, NULL) ON DUPLICATE KEY UPDATE username = ?", name, pass, name)
	header := deleteExecute(t, db, name, pass)
	if header["Status"][0] != fmt.Sprintf("delete %s", name) {
		t.Errorf("Deletion of %s failed.", name)
	} else {
		user := &protocol.User{}
		flag := protocol.IsExistingUsername(db, name, user)
		if flag {
			t.Errorf("Deletion not taking expected effect.")
		}
	}
}
