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

func test_restricted_template_resulting_header(t *testing.T, db *sql.DB, cookieString string, url string, fn func(*sql.DB, http.ResponseWriter, *http.Request)) http.Header {
	response := httptest.NewRecorder()
	request := makeRequest(http.MethodGet, url, t)
	updateCookie(cookieString, response, request)
	http.HandlerFunc(makeHandler(db, fn)).ServeHTTP(response, request)
	header := response.Header()
	return header
}

func test_restricted_template_granted_access(t *testing.T, db *sql.DB, i int, url string, fn func(*sql.DB, http.ResponseWriter, *http.Request)) {
	cookieString := handlers.SetSessionInfo(
		&protocol.User{
			Name:     fmt.Sprintf("user%d", i),
			Password: fmt.Sprintf("pass%d%d", i*2, i*2),
		},
		&protocol.User{},
		handlers.InfoErr{},
		"",
	)
	header := test_restricted_template_resulting_header(t, db, cookieString, url, fn)
	if header["Status"][0] != "successful view" {
		t.Errorf("Failed access to restricted page for user%d", i)
	}
}

func test_restricted_template_no_access(t *testing.T, db *sql.DB, url string, fn func(*sql.DB, http.ResponseWriter, *http.Request)) {
	cookieString := handlers.SetSessionInfo(
		&protocol.User{},
		&protocol.User{},
		handlers.InfoErr{},
		"",
	)
	header := test_restricted_template_resulting_header(t, db, cookieString, url, fn)
	if header["Status"][0] != "login error" {
		t.Errorf("Wrongly granted access.")
	}
}
