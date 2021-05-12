package main

import (
	"fmt"
	"git.garena.com/jiayu.li/entry-task/test/server_helpers"
	"net/url"
	"testing"
)

func Test_loginRequests(t *testing.T) {
	for j := 0; j < 5; j++ {
		j := j
		for i := 0; i < 200; i++ {
			i := i
			client := server_helpers.PrepareClient(t)
			t.Run(fmt.Sprintf("Login to user%d#%d", i, j), func(t *testing.T) {
				t.Parallel()
				resp, err := client.PostForm("http://localhost:8080/login", url.Values{
					"name": {fmt.Sprintf("user%d", i)},
					"password": {fmt.Sprintf("pass%d%d", i * 2, i * 2)},
				})
				if err != nil {
					t.Errorf("Error login, %s", err.Error())
				}
				resp.Body.Close()
			})
		}
	}
}
