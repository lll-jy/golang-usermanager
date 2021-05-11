package test

import (
	"fmt"
	"testing"

	"git.garena.com/jiayu.li/entry-task/cmd/handlers"
	"git.garena.com/jiayu.li/entry-task/cmd/paths"
	_ "github.com/go-sql-driver/mysql"
)

func Test_handlers(t *testing.T) {
	db := setupDb(t)
	handlers.PrepareTemplates("../templates/%s.html")
	paths.SetupPaths("test")

	t.Run("Restricted access", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			test_restricted_template_granted_access(t, db, i, "/view", handlers.ViewPageHandler)
			test_restricted_template_granted_access(t, db, i, "/reset", handlers.ResetPageHandler)
			test_restricted_template_granted_access(t, db, i, "/edit", handlers.EditPageHandler)
		}
		test_restricted_template_no_access(t, db, "/view", handlers.ViewPageHandler)
		test_restricted_template_no_access(t, db, "/reset", handlers.ResetPageHandler)
		test_restricted_template_no_access(t, db, "/edit", handlers.EditPageHandler)
	})

	t.Run("Login", func(t *testing.T) {
		clearEffects(db)
		for i := 0; i < 5; i++ {
			test_single_user_login(t, i, db)
		}
		for i := 0; i < 5; i++ {
			test_invalid_user_login(t, fmt.Sprintf("name=user%d&password=pass%d", i, i), db, "incorrect password", fmt.Sprintf("Wrong password for user%d not detected correctly.", i))
		}
		for i := 0; i < 5; i++ {
			test_invalid_user_login(t, fmt.Sprintf("name=useruser%d&password=pass%d", i, i), db, "user not exist", fmt.Sprintf("Non-existing user useruser%d not detected correctly.", i))
		}
	})

	t.Run("Signup", func(t *testing.T) {
		clearEffects(db)
		for i := 0; i < 5; i++ {
			test_valid_signup(t, db, i)
		}
		for i := 0; i < 5; i++ {
			test_invalid_signup(t, db, fmt.Sprintf("user%d", i), "pass", "pass", "user already exists", fmt.Sprintf("Failed to recognized existing user user%d", i))
			name := fmt.Sprintf("testuser%d", 100+i)
			pass := fmt.Sprintf("pass%d%d", 200+i, 200+i)
			test_invalid_signup(t, db, name, pass, fmt.Sprintf("pass%d", 200+i), "mismatch password", "Failed to detect password mismatch.")
			wrong_pass := fmt.Sprintf("p%d", i)
			test_invalid_signup(t, db, name, wrong_pass, wrong_pass, "wrong password format", "Failed to detect wrong password format.")
			test_invalid_signup(t, db, fmt.Sprintf("testuser,%d", i), pass, pass, "wrong username format", "Failed to detect wrong username format.")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			test_valid_delete(t, db, i)
		}
	})

	t.Run("Edit", func(t *testing.T) {
		clearEffects(db)
		for i := 0; i < 5; i++ {
			test_valid_edit(t, db, i)
		}
	})

	t.Run("Reset", func(t *testing.T) {
		clearEffects(db)
		for i := 0; i < 5; i++ {
			test_valid_reset_pass(t, db, i)
		}
		clearEffects(db)
		for i := 0; i < 5; i++ {
			test_valid_reset_name(t, db, i)
		}
		for i := 0; i < 5; i++ {
			test_invalid_reset_duplicate(t, db, i)
		}
	})

	t.Run("Upload", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			test_upload(t, db, i)
		}
	})
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

/*func Test_server(t *testing.T) {
t.Run("Test run", func(t *testing.T) {
	db := setupDb(t)

	t.Run("Single user test", func(t *testing.T) {
		test_single_user(t, 0, db)
	})

	/*t.Run("Login handler", func(t *testing.T) {
		response := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader("name=user2&password=pass44"))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if err != nil {
			t.Errorf("Request error %s.", err)
		}
		http.HandlerFunc(makeHandler(db, handlers.LoginHandler)).ServeHTTP(response, request)
		fmt.Printf("Loggin: %s\n", response.Header()["Status"])
		fmt.Printf("Check %t", response.Header()["Status"][0] == "successful login")
		// fmt.Printf("Result is\n%s\n", response.Body.String())

		request, err = http.NewRequest(http.MethodPost, "/login", strings.NewReader("name=user2&password=pass4"))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if err != nil {
			t.Errorf("Request error %s.", err)
		}
		http.HandlerFunc(makeHandler(db, handlers.LoginHandler)).ServeHTTP(response, request)
	})*/

/*t.Run("Index page", func(t *testing.T) {
	response := httptest.NewRecorder()
	request := makeRequest(http.MethodGet, "/", t)
	makeHandler(db, handlers.IndexPageHandler)(response, request)
	// t.Logf("Index body is %s", response.Body.String())
	// fmt.Printf("Index body is %s\n", response.Body.String())
})*/

/*t.Run("Cookie handler", func(t *testing.T) {
	response := httptest.NewRecorder()
	request := makeRequest(http.MethodGet, "/view", t)

	// makeHandler(db, handlers.LoginHandler)(response, request)
	// fmt.Printf("1 Login Body is %s\n", response.Body.String())

	cookieString := handlers.SetSessionInfo(
		&protocol.User{
			Name:     "user3",
			Password: "pass66",
		},
		&protocol.User{
			Name:     "user3",
			Password: "pass66",
		},
		handlers.InfoErr{},
		"",
	)
	cookie := &http.Cookie{
		Name:  "session",
		Value: cookieString,
		Path:  "/",
	}
	http.SetCookie(response, cookie)
	request.Header = http.Header{"Cookie": response.HeaderMap["Set-Cookie"]}
	// fmt.Printf("2 Login Body is %s\n", response.Body.String())
	cookie, err = request.Cookie("session")
	// fmt.Printf("%v\n", cookie)
	fmt.Printf("%v\n", handlers.GetPageInfo(request))
	// fmt.Printf("%v\n", response.Result().Cookies()[0].Value)

	//ctx := request.Context()
	//ctx = context.WithValue(ctx, "session", cookieString)
	//request = request.WithContext(ctx)

	http.HandlerFunc(makeHandler(db, handlers.ViewPageHandler)).ServeHTTP(response, request)
	fmt.Printf("Login body is \n%s\n", response.Body.String())
})*/

/*t.Run("Cookie test", func(t *testing.T) {
	cookieString := handlers.SetSessionInfo(
		&protocol.User{
			Name:     "user3",
			Password: "pass66",
		},
		&protocol.User{
			Name:     "user3",
			Password: "pass66",
		},
		handlers.InfoErr{},
		"",
	)
	cookie := &http.Cookie{
		Name:  "session",
		Value: cookieString,
		Path:  "/",
	}
	tt := New("/login", makeHandler(db, handlers.LoginHandler), t)
	tt.AddCookies(cookie)
	tt.Do()
	fmt.Printf(tt.GetBody())
	//fmt.Printf("%s", cookie)
})*/
/*	})
}*/

/*func Test_dummy2(t *testing.T) {
	f, err := os.Open("../templates/index.html")
	if err != nil {
		t.Logf("Cannot find file %s", err.Error())
	} else {
		t.Logf("Opened %v", f)
	}
}*/

/*func Test_dummy(t *testing.T) {
	t.Run("hello world", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		PlayerServer(response, request)
		// got := response.Body.String()
		// fmt.Println(got)
		// t.Log(got)
	})
}

func PlayerServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "20")
}

func TestSetCookie(t *testing.T) {
	response := httptest.NewRecorder()
	http.SetCookie(response, &http.Cookie{
		Name:  "test",
		Value: "hello",
	})
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	request.Header = http.Header{"Cookie": response.HeaderMap["Set-Cookie"]}

	//cookie, _ := request.Cookie("test")

	//fmt.Println(cookie.Value)
}
*/
