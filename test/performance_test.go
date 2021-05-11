package main

import (
	"fmt"
	"testing"
	"time"

	"git.garena.com/jiayu.li/entry-task/cmd/handlers"
	"git.garena.com/jiayu.li/entry-task/cmd/paths"
	"git.garena.com/jiayu.li/entry-task/test/server_helpers"
	_ "github.com/go-sql-driver/mysql"
)

func Test_massiveRequests(t *testing.T) {
	db := server_helpers.SetupDb(t)
	db.SetMaxOpenConns(500)
	db.SetMaxIdleConns(400)
	db.SetConnMaxLifetime(time.Minute * 3)
	handlers.PrepareTemplates("../templates/%s.html")
	paths.SetupPaths("test")
	t.Log("Start logging in.")
	//handlers.Initialize(db)
	var start time.Time
	/*t.Run("Test login requests handling speed", func(t *testing.T) {
		start = time.Now()
		for j := 0; j < 5; j++ {
			j := j
			for i := 0; i < 200; i++ {
				i := i
				t.Run(fmt.Sprintf("Login to user%d#%d", i, j), func(t *testing.T) {
					t.Parallel()
					go server_helpers.ValidLogin(t, i, db)
				})
			}
		}
	})
	defer db.Close()

	end := time.Now()
	dur := end.Sub(start).Seconds()
	if dur > 1 {
		t.Errorf("Time running 1000 login requests exceeds 1 second. Used %v seconds.", dur)
	}*/

	t.Run("Test mixed requests handling speed", func(t *testing.T) {
		start = time.Now()
		for i := 0; i < 200; i++ {
			/*t.Run(fmt.Sprintf("Login to user%d", i), func(t *testing.T) {
				t.Parallel()
				go server_helpers.ValidLogin(t, i, db)
			})
			t.Run(fmt.Sprintf("View user%d", i), func(t *testing.T) {
				t.Parallel()
				go server_helpers.GrantedRestrictedTemplate(t, db, i, "/view", handlers.ViewPageHandler)
			})
			t.Run(fmt.Sprintf("Edit user%d", i), func(t *testing.T) {
				t.Parallel()
				go server_helpers.ValidEditNickname(t, db, i)
			})*/
			t.Run(fmt.Sprintf("Edit user photo of user%d", i), func(t *testing.T) {
				t.Parallel()
				go server_helpers.ValidEditPhoto(t, db, i) // contains 2 requests
			})
		}
	})

	end := time.Now()
	dur := end.Sub(start).Seconds()
	if dur > 1 {
		t.Errorf("Time running 1000 mixed requests exceeds 1 second. Used %v seconds.", dur)
	}
}
