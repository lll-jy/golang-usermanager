package test

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"git.garena.com/jiayu.li/entry-task/cmd/handlers"
	"git.garena.com/jiayu.li/entry-task/cmd/protocol"
	"google.golang.org/protobuf/proto"
)

func makeHandler(db *sql.DB, fn func(*sql.DB, http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(db, w, r)
	}
}

func makeRequest(method string, url string, t *testing.T) *http.Request {
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		t.Errorf("Request failure :%v", err.Error())
	}
	return request
}
func setupDb(t *testing.T) *sql.DB {
	db, err := sql.Open("mysql", "root:@/entryTask")
	if err != nil {
		t.Errorf("Database connection failed: %v", err.Error())
	}
	return db
}

func getUser(h http.Header) *protocol.User {
	user := &protocol.User{}
	proto.Unmarshal([]uint8(h["User"][0]), user)
	return user
}

func formSetup(body string, t *testing.T, db *sql.DB, url string) (*httptest.ResponseRecorder, *http.Request) {
	response := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Errorf("Request error %s.", err)
	}
	return response, request
}

func updateCookie(cookieString string, response *httptest.ResponseRecorder, request *http.Request) {
	cookie := &http.Cookie{
		Name:  "session",
		Value: cookieString,
		Path:  "/",
	}
	http.SetCookie(response, cookie)
	request.Header["Cookie"] = response.HeaderMap["Set-Cookie"]
}

// https://stackoverflow.com/questions/29505089/how-can-i-compare-two-files-in-golang
func areIdenticalFiles(file1 string, file2 string) bool {
	chunk := 64000
	f1, err := os.Open(file1)
	if err != nil {
		return false
	}
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		return false
	}
	defer f2.Close()
	for {
		b1 := make([]byte, chunk)
		_, err1 := f1.Read(b1)

		b2 := make([]byte, chunk)
		_, err2 := f2.Read(b2)

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true
			} else {
				return false
			}
		}
		if !bytes.Equal(b1, b2) {
			return false
		}
	}
}

func clearEffects(db *sql.DB) {
	handlers.ExecuteQuery(db, "DELETE FROM users WHERE username LIKE 'test%'")
	for i := 0; i < 5; i++ {
		handlers.ExecuteQuery(db, "UPDATE users SET photo = ?, nickname = ? WHERE username = ?", nil, fmt.Sprintf("nick%d", i), fmt.Sprintf("user%d", i))
	}
}
