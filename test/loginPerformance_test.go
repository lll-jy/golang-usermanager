package main

import (
	"fmt"
	"git.garena.com/jiayu.li/entry-task/cmd/handlers"
	"testing"
	"time"

	"git.garena.com/jiayu.li/entry-task/test/server_helpers"
	_ "github.com/go-sql-driver/mysql"
)

func Test_massiveLogins(t *testing.T) {
	db := server_helpers.Setup(t)
	handlers.Initialize(db)
	var start time.Time
	t.Run("Test login requests handling speed", func(t *testing.T) {
		start = time.Now()
		for j := 0; j < 5; j++ {
			j := j
			for i := 0; i < 200; i++ {
				i := i
				go t.Run(fmt.Sprintf("Login to user%d#%d", i, j), func(t *testing.T) {
					t.Parallel()
					server_helpers.ValidLogin(t, i, db)
				})
			}
		}
	})
	defer db.Close()

	end := time.Now()
	dur := end.Sub(start).Seconds()
	if dur > 1 {
		t.Errorf("Time running 1000 login requests exceeds 1 second. Used %v seconds.", dur)
	}
}

// sudo chown -R <username> .config