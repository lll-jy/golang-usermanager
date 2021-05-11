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

func Test_massiveLogin(t *testing.T) {
	db := server_helpers.SetupDb(t)
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(100)
	db.SetConnMaxLifetime(time.Minute * 3)
	handlers.PrepareTemplates("../templates/%s.html")
	paths.SetupPaths("test")
	t.Log("Start logging in.")
	//handlers.Initialize(db)
	var start time.Time
	t.Run("Test speed to handle requests", func(t *testing.T) {
		start = time.Now()
		for j := 0; j < 5; j++ {
			j := j
			for i := 0; i < 200; i++ {
				i := i
				t.Run(fmt.Sprintf("Login to user%d#%d", i, j), func(t *testing.T) {
					t.Parallel()
					server_helpers.Login(t, db, i)
				})
			}
		}
	})

	end := time.Now()
	dur := end.Sub(start).Seconds()
	if dur > 1 {
		t.Errorf("Time exceeds 1 seconds. Used %v seconds.", dur)
	}
}
