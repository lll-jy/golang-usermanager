package main

import (
	"fmt"
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