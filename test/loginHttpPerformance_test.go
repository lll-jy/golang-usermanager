package main

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

func Test_loginRequests(t *testing.T) {
	for j := 0; j < 5; j++ {
		j := j
		for i := 0; i < 200; i++ {
			i := i
			t.Run(fmt.Sprintf("Login to user%d#%d", i, j), func(t *testing.T) {
				t.Parallel()
				resp, err := http.PostForm("http://localhost:8080/login", url.Values{
					"name": {"user3"},
					"password": {"pass66"},
				})
				if err != nil {
					t.Errorf("Error login, %s", err.Error())
				}
				//t.Errorf("%v", resp.Header)
				resp.Body.Close()
			})
		}
	}
}