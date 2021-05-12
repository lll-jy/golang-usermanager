package main

import (
	"fmt"
	"git.garena.com/jiayu.li/entry-task/cmd/handlers"
	"git.garena.com/jiayu.li/entry-task/cmd/paths"
	"git.garena.com/jiayu.li/entry-task/cmd/protocol"
	"git.garena.com/jiayu.li/entry-task/test/server_helpers"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"testing"
)

func prepareClient(t *testing.T) *http.Client {
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Errorf("Cannot creat cookie jar.")
	}
	return &http.Client{
		Jar: jar,
	}
}

func Test_mixedRequests(t *testing.T) {
	for i := 0; i < 200; i++ {
		i := i
		t.Run(fmt.Sprintf("Login, view, and edit user%d", i), func(t *testing.T) {
			t.Parallel()
			client := prepareClient(t)
			resp, err := client.PostForm("http://localhost:8080/login", url.Values{
				"name": {fmt.Sprintf("user%d", i)},
				"password": {fmt.Sprintf("pass%d%d", i * 2, i * 2)},
			})
			if err != nil {
				t.Errorf("Error login, %s", err.Error())
			}
			resp.Body.Close()
			resp, err = client.Get("http://localhost:8080/view")
			if err != nil {
				t.Errorf("Error view, %s", err.Error())
			}
			resp.Body.Close()
			resp, err = client.PostForm("http://localhost:8080/edit", url.Values{
				"nickname": {fmt.Sprintf("nickname%d", i)},
			})
			if err != nil {
				t.Errorf("Error edit, %s", err.Error())
			}
			resp.Body.Close()
		})
		t.Run(fmt.Sprintf("Signup, delete testuser%d", i), func(t *testing.T) {
			t.Parallel()
			client := prepareClient(t)
			resp, err := client.PostForm("http://localhost:8080/signup", url.Values{
				"name": {fmt.Sprintf("testuser%d", i)},
				"password": {fmt.Sprintf("testpass%d%d", i * 2, i * 2)},
				"password_repeat": {fmt.Sprintf("testpass%d%d", i * 2, i * 2)},
			})
			if err != nil {
				t.Errorf("Error signup, %s", err.Error())
			}
			resp.Body.Close()
			resp, err = client.PostForm("http://localhost:8080/delete", url.Values{})
			if err != nil {
				t.Errorf("Error delete, %s", err.Error())
			}
			resp.Body.Close()
		})
	}
}

func abort_Test_mixedRequests(t *testing.T) {
	for i := 0; i < 200; i++ {
		i := i
		client := &http.Client{}
		t.Run(fmt.Sprintf("Login to user%d", i), func(t *testing.T) {
			t.Parallel()
			resp, err := client.PostForm("http://localhost:8080/login", url.Values{
				"name": {fmt.Sprintf("user%d", i)},
				"password": {fmt.Sprintf("pass%d%d", i * 2, i * 2)},
			})
			if err != nil {
				t.Errorf("Error login, %s", err.Error())
			}
			resp.Body.Close()

			/*resp, err = client.PostForm("http://localhost:8080/reset", url.Values{
				"name": {fmt.Sprintf("user%d", i)},
				"password": {fmt.Sprintf("pass%d%d", i * 2, i * 3)},
				"password_repeat": {fmt.Sprintf("pass%d%d", i * 2, i * 3)},
			})
			if err != nil {
				t.Errorf("Error reset, %s", err.Error())
			}
			resp.Body.Close()*/
		})
		/*t.Run(fmt.Sprintf("Signup testuser%d", i), func(t *testing.T) {
			t.Parallel()
			resp, err := client.PostForm("http://localhost:8080/signup", url.Values{
				"name": {fmt.Sprintf("testuser%d", i)},
				"password": {fmt.Sprintf("testpass%d%d", i * 2, i * 2)},
				"password_repeat": {fmt.Sprintf("testpass%d%d", i * 2, i * 2)},
			})
			if err != nil {
				t.Errorf("Error signup, %s", err.Error())
			}
			resp.Body.Close()
		})*/
		t.Run(fmt.Sprintf("Delete testuser%d", i), func(t *testing.T) {
			t.Parallel()
			cookieString := handlers.SetSessionInfo(
				&protocol.User{
					Name:     fmt.Sprintf("testuser%d", i),
					Password: fmt.Sprintf("testpass%d%d", i*2, i*2),
				},
				&protocol.User{
					Name:     fmt.Sprintf("testuser%d", i),
					Password: fmt.Sprintf("testpass%d%d", i*2, i*2),
				},
				handlers.InfoErr{},
				paths.PlaceholderPath,
			)
			request := server_helpers.MakeRequest(http.MethodPost, "http://localhost:8080/delete", t)
			t.Logf("header is %v", request.Header)
			request.Header.Set("Cookie", cookieString)
			//t.Logf("header is %v", request.Header)
			//request.AddCookie(&http.Cookie{Name: "session", Value: cookieString})
			t.Logf("header2 is %v", request.Header)
			//server_helpers.UpdateCookie(cookieString, httptest.NewRecorder(), request)
			resp, err := client.Do(request)
			if err != nil {
				t.Errorf("Error delete, %s", err.Error())
			}
			resp.Body.Close()
		})
	}
}
