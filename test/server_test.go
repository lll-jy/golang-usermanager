package main

import (
	"fmt"
	"os"
	"testing"

	"git.garena.com/jiayu.li/entry-task/cmd/handlers"
	"git.garena.com/jiayu.li/entry-task/cmd/paths"
	"git.garena.com/jiayu.li/entry-task/test/server_helpers"
	_ "github.com/go-sql-driver/mysql"
)

func Test_handlers(t *testing.T) {
	db := server_helpers.SetupDb(t)
	handlers.PrepareTemplates("../templates/%s.html")
	paths.SetupPaths("test")

	t.Run("Restricted access", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			server_helpers.GrantedRestrictedTemplate(t, db, i, "/view", handlers.ViewPageHandler)
			server_helpers.GrantedRestrictedTemplate(t, db, i, "/reset", handlers.ResetPageHandler)
			server_helpers.GrantedRestrictedTemplate(t, db, i, "/edit", handlers.EditPageHandler)
		}
		server_helpers.DeniedRestrictedTemplate(t, db, "/view", handlers.ViewPageHandler)
		server_helpers.DeniedRestrictedTemplate(t, db, "/reset", handlers.ResetPageHandler)
		server_helpers.DeniedRestrictedTemplate(t, db, "/edit", handlers.EditPageHandler)
	})

	t.Run("Login", func(t *testing.T) {
		server_helpers.ClearEffects(db)
		for i := 0; i < 5; i++ {
			server_helpers.ValidLogin(t, i, db)
		}
		for i := 0; i < 5; i++ {
			server_helpers.InvalidLogin(t, fmt.Sprintf("name=user%d&password=pass%d", i, i), db, "incorrect password", fmt.Sprintf("Wrong password for user%d not detected correctly.", i))
		}
		for i := 0; i < 5; i++ {
			server_helpers.InvalidLogin(t, fmt.Sprintf("name=useruser%d&password=pass%d", i, i), db, "user not exist", fmt.Sprintf("Non-existing user useruser%d not detected correctly.", i))
		}
	})

	t.Run("Signup", func(t *testing.T) {
		server_helpers.ClearEffects(db)
		for i := 0; i < 5; i++ {
			server_helpers.ValidSignup(t, db, i)
		}
		for i := 0; i < 5; i++ {
			server_helpers.InvalidSignup(t, db, fmt.Sprintf("user%d", i), "pass", "pass", "user already exists", fmt.Sprintf("Failed to recognized existing user user%d", i))
			name := fmt.Sprintf("testuser%d", 100+i)
			pass := fmt.Sprintf("pass%d%d", 200+i, 200+i)
			server_helpers.InvalidSignup(t, db, name, pass, fmt.Sprintf("pass%d", 200+i), "mismatch password", "Failed to detect password mismatch.")
			wrong_pass := fmt.Sprintf("p%d", i)
			server_helpers.InvalidSignup(t, db, name, wrong_pass, wrong_pass, "wrong password format", "Failed to detect wrong password format.")
			server_helpers.InvalidSignup(t, db, fmt.Sprintf("testuser,%d", i), pass, pass, "wrong username format", "Failed to detect wrong username format.")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			server_helpers.ValidDelete(t, db, i)
		}
	})

	t.Run("Upload", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			server_helpers.Upload(t, db, i)
		}
	})

	t.Run("Edit", func(t *testing.T) {
		server_helpers.ClearEffects(db)
		for i := 0; i < 5; i++ {
			server_helpers.ValidEditNickname(t, db, i)
		}
		server_helpers.ClearEffects(db)
		for i := 0; i < 5; i++ {
			server_helpers.ValidEditPhoto(t, db, i)
		}
	})

	t.Run("Reset", func(t *testing.T) {
		server_helpers.ClearEffects(db)
		for i := 0; i < 5; i++ {
			server_helpers.ValidResetPass(t, db, i)
		}
		server_helpers.ClearEffects(db)
		for i := 0; i < 5; i++ {
			server_helpers.ValidResetName(t, db, i)
		}
		server_helpers.ClearEffects(db)
		for i := 0; i < 5; i++ {
			server_helpers.ValidResetPassWithPhoto(t, db, i)
		}
		for i := 0; i < 5; i++ {
			server_helpers.InvalidResetDuplicate(t, db, i)
		}
	})

	err := os.RemoveAll(paths.TempPath)
	if err != nil {
		t.Errorf("Cannot remove directory, %v", err.Error())
	}
	os.MkdirAll(paths.TempPath, 0777)
	err = os.RemoveAll(paths.FileBaseRelativePath)
	if err != nil {
		t.Errorf("Cannot remove directory, %v", err.Error())
	}
	os.MkdirAll(paths.FileBaseRelativePath, 0777)
}

/*func Test_login(t *testing.T) {
	t.Run("Test simultaneous login", func(t *testing.T) {
		db := setupDb(t)

		start := time.Now()
		for j := 0; j < 3; j++ {
			for i := 100; i < 110; i++ {
				i := i
				t.Run(fmt.Sprintf("Test login to user%d", 0), func(t *testing.T) {
					t.Parallel()
					test_single_user_login(t, i, db)
				})
			}
		}
		end := time.Now()
		dur := end.Sub(start).Seconds()
		t.Logf("It takes %v seconds to complete.", dur)
		if dur > 1 {
			t.Errorf("Time out")
		}
	})
}*/
