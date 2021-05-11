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

func Test_speed(t *testing.T) {
	db := server_helpers.SetupDb(t)
	handlers.PrepareTemplates("../templates/%s.html")
	paths.SetupPaths("test")
	//handlers.Initialize(db)
	start := time.Now()
	t.Run("Test speed to handle requests", func(t *testing.T) {
		for j := 0; j < 3; j++ {
			j := j
			for i := 0; i < 200; i++ {
				i := i
				t.Run(fmt.Sprintf("Login to user%d#%d", i, j), func(t *testing.T) {
					//server_helpers.ValidLogin(t, i, db)
					server_helpers.LoginExecute(fmt.Sprintf("name=user%d&password=pass%d%d", i, i*2, i*2), t, db)
				})
			}
		}
	})

	end := time.Now()
	dur := end.Sub(start).Seconds()
	t.Errorf("Used %v seconds", dur)
}
