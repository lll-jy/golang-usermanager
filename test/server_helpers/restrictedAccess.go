package server_helpers

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.garena.com/jiayu.li/entry-task/cmd/handlers"
	"git.garena.com/jiayu.li/entry-task/cmd/protocol"
)

func ExecuteRestrictedTemplate(t *testing.T, db *sql.DB, cookieString string, url string, fn func(*sql.DB, http.ResponseWriter, *http.Request)) http.Header {
	response := httptest.NewRecorder()
	request := makeRequest(http.MethodGet, url, t)
	updateCookie(cookieString, response, request)
	http.HandlerFunc(makeHandler(db, fn)).ServeHTTP(response, request)
	header := response.Header()
	return header
}

func GrantedRestrictedTemplate(t *testing.T, db *sql.DB, i int, url string, fn func(*sql.DB, http.ResponseWriter, *http.Request)) {
	cookieString := handlers.SetSessionInfo(
		&protocol.User{
			Name:     fmt.Sprintf("user%d", i),
			Password: fmt.Sprintf("pass%d%d", i*2, i*2),
		},
		&protocol.User{},
		handlers.InfoErr{},
		"",
	)
	header := ExecuteRestrictedTemplate(t, db, cookieString, url, fn)
	if header["Status"][0] != "successful view" {
		t.Errorf("Failed access to restricted page for user%d", i)
	}
}

func DeniedRestrictedTemplate(t *testing.T, db *sql.DB, url string, fn func(*sql.DB, http.ResponseWriter, *http.Request)) {
	cookieString := handlers.SetSessionInfo(
		&protocol.User{},
		&protocol.User{},
		handlers.InfoErr{},
		"",
	)
	header := ExecuteRestrictedTemplate(t, db, cookieString, url, fn)
	if header["Status"][0] != "login error" {
		t.Errorf("Wrongly granted access.")
	}
}
