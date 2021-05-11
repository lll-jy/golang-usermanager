package server_helpers

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.garena.com/jiayu.li/entry-task/cmd/handlers"
)

func LoginExecute(body string, t *testing.T, db *sql.DB) *httptest.ResponseRecorder {
	response, request := formSetup(body, t, db, "/login")
	http.HandlerFunc(makeHandler(db, handlers.LoginHandler)).ServeHTTP(response, request)
	return response
}

func ValidLogin(t *testing.T, i int, db *sql.DB) {
	response := LoginExecute(fmt.Sprintf("name=user%d&password=pass%d%d", i, i*2, i*2), t, db)
	header := response.Header()
	user := getUser(header)
	if user.Name != fmt.Sprintf("user%d", i) || header["Status"][0] != "successful login" {
		t.Errorf("Login to user%d unsuccessful.", i)
	} else if user.Nickname != fmt.Sprintf("nick%d", i) {
		t.Errorf("Wrong information retrieved.")
	}
}

func InvalidLogin(t *testing.T, info string, db *sql.DB, expectedErr string, errString string) {
	response := LoginExecute(info, t, db)
	header := response.Header()
	if header["Status"][0] != expectedErr {
		t.Errorf(errString)
	}
}
