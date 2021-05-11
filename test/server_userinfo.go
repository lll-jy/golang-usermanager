package test

import (
	"database/sql"
	"fmt"
	"net/http"
	"testing"

	"git.garena.com/jiayu.li/entry-task/cmd/handlers"
	"git.garena.com/jiayu.li/entry-task/cmd/protocol"
)

func signupExecute(t *testing.T, db *sql.DB, name string, pass string) http.Header {
	handlers.ExecuteQuery(db, "DELETE FROM users WHERE username LIKE 'test%'")
	response, request := formSetup(fmt.Sprintf("name=%s&password=%s&password_repeat=%s", name, pass, pass), t, db, "/signup")
	http.HandlerFunc(makeHandler(db, handlers.SignupHandler)).ServeHTTP(response, request)
	return response.Header()
}

func test_valid_signup(t *testing.T, db *sql.DB, i int) {
	header := signupExecute(t, db, fmt.Sprintf("testuser%d", i), fmt.Sprintf("testpass%d%d", i*2, i*2))
	user := getUser(header)
	if user.Name != fmt.Sprintf("testuser%d", i) || header["Status"][0] != "successful signup" {
		t.Errorf("Sign up for testuser%d unssucessful.", i)
	}
}

func test_invalid_signup(t *testing.T, db *sql.DB, name string, pass string, repeat string, status string, errString string) {
	response, request := formSetup(fmt.Sprintf("name=%s&password=%s&password_repeat=%s", name, pass, repeat), t, db, "/signup")
	http.HandlerFunc(makeHandler(db, handlers.SignupHandler)).ServeHTTP(response, request)
	header := response.Header()
	if header["Status"][0] != status {
		t.Errorf(errString)
	}
}

func resetExecute(t *testing.T, db *sql.DB, old_name string, old_pass string, new_name string, new_pass string) http.Header {
	signupExecute(t, db, old_name, old_pass)
	user := &protocol.User{}
	protocol.IsExistingUsername(db, old_name, user)
	user.Name = old_name
	cookieString := handlers.SetSessionInfo(
		user,
		&protocol.User{
			Name:     old_name,
			Password: old_pass,
		},
		handlers.InfoErr{},
		"",
	)
	response, request := formSetup(fmt.Sprintf("name=%s&password=%s&password_repeat=%s", new_name, new_pass, new_pass), t, db, "/reset")
	updateCookie(cookieString, response, request)
	http.HandlerFunc(makeHandler(db, handlers.ResetHandler)).ServeHTTP(response, request)
	return response.Header()
}

func test_valid_reset_pass(t *testing.T, db *sql.DB, i int) {
	name := fmt.Sprintf("testuser%d", i)
	pass := fmt.Sprintf("testpass%d%d", i*2, i*2)
	newPass := fmt.Sprintf("testnewpass%d%d", i*2, i*2)
	header := resetExecute(t, db, name, pass, name, newPass)
	user := getUser(header)
	if header["Status"][0] != "successful signup" {
		t.Errorf("Failed to reset password for %s due to %s.", name, header["Status"][0])
	} else if user.Name != name {
		t.Errorf("Wrongly reset username.")
	} else {
		user.Name = ""
		flag := protocol.IsExistingUsername(db, name, user)
		if !flag {
			t.Errorf("Wrongly changed username in database.")
		} else if !protocol.IsCorrectPassword(newPass, user.Password) {
			t.Errorf("Failed to update password for %s.", name)
		}
	}
}

func test_valid_reset_name(t *testing.T, db *sql.DB, i int) {
	name := fmt.Sprintf("testuser%d", i)
	newName := fmt.Sprintf("testnewuser%d", i)
	pass := fmt.Sprintf("testpass%d%d", i*2, i*2)
	header := resetExecute(t, db, name, pass, newName, pass)
	user := getUser(header)
	if header["Status"][0] != "successful signup" {
		t.Errorf("Failed to reset username for %s to %s.", name, newName)
	} else if user.Name != newName {
		t.Errorf("Failed to update username.")
	} else {
		user.Name = ""
		flag := protocol.IsExistingUsername(db, name, user)
		if !flag {
			t.Errorf("Wrongly changed username in database.")
		} else if !protocol.IsCorrectPassword(pass, user.Password) {
			t.Errorf("Wrongly updated password when changing username for %s to %s.", name, newName)
		}
	}
}

func test_invalid_reset_duplicate(t *testing.T, db *sql.DB, i int) {
	name := fmt.Sprintf("testuser%d", i)
	newName := fmt.Sprintf("user%d", i)
	pass := fmt.Sprintf("testpass%d%d", i*2, i*2)
	header := resetExecute(t, db, name, pass, newName, pass)
	if header["Status"][0] != "user already exists" {
		t.Errorf("Wrongly allowed reset username to existing user")
	}
}
