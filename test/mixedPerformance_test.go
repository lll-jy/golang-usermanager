package main

import (
	"fmt"
	"git.garena.com/jiayu.li/entry-task/cmd/handlers"
	"testing"
	"time"

	"git.garena.com/jiayu.li/entry-task/test/server_helpers"
	_ "github.com/go-sql-driver/mysql"
)
func Test_massiveRequests(t *testing.T) {
	db := server_helpers.Setup(t)
	handlers.Initialize(db)
	start := time.Now()

	t.Run("Test mixed requests handling speed", func(t *testing.T) {
		t.Run("Login", func(t *testing.T) {
			t.Parallel()
			for i := 0; i < 200; i++ {
				go t.Run(fmt.Sprintf("user%d", i), func(t *testing.T) {
					t.Parallel()
					server_helpers.ValidLogin(t, i, db)
				})
			}
		})
		/*go t.Run("View", func(t *testing.T) {
			t.Parallel()
			for i := 0; i < 200; i++ {
				go t.Run(fmt.Sprintf("user%d", i), func(t *testing.T) {
					t.Parallel()
					server_helpers.GrantedRestrictedTemplate(t, db, i, "/view", handlers.ViewPageHandler)
				})
			}
		})
		/*t.Run("Edit nickname", func(t *testing.T) {
			t.Parallel()
			for i := 0; i < 200; i++ {
				t.Run(fmt.Sprintf("user%d", i), func(t *testing.T) {
					t.Parallel()
					go server_helpers.ValidEditNickname(t, db, i)
				})
			}
		})
		t.Run("Reset password", func(t *testing.T) {
			t.Parallel()
			for i := 0; i < 200; i++ {
				t.Run(fmt.Sprintf("user%d", i), func(t *testing.T) {
					t.Parallel()
					go server_helpers.ValidResetPass(t, db, i) // signup + reset, 2 requests
				})
			}
		})*/
		//server_helpers.ValidEditNickname(t, db, 0)
	})

	end := time.Now()
	dur := end.Sub(start).Seconds()
	if dur > 1 {
		t.Errorf("Time running 1000 mixed requests exceeds 1 second. Used %v seconds. %v", dur, end.Sub(start))
	}
}
