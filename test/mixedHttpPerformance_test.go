package main

import (
	"fmt"
	"git.garena.com/jiayu.li/entry-task/test/server_helpers"
	"net/url"
	"testing"
)

func Test_mixedRequests(t *testing.T) {
	for i := 0; i < 200; i++ {
		i := i
		t.Run(fmt.Sprintf("Login, view, and edit user%d", i), func(t *testing.T) {
			t.Parallel()
			client := server_helpers.PrepareClient(t)
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
			client := server_helpers.PrepareClient(t)
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

// command for generating profile
// go test test/l*.go -parallel 100 -cpuprofile test/profiles/login_cpu.prof -memprofile test/profiles/login_mem.prof -bench .
// go tool pprof test/profiles/...prof