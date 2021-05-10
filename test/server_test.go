package test

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"git.garena.com/jiayu.li/entry-task/cmd/handlers"
	"git.garena.com/jiayu.li/entry-task/cmd/protocol"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/protobuf/proto"
)

func makeHandler(db *sql.DB, fn func(*sql.DB, http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(db, w, r)
	}
}

func makeRequest(method string, url string, t *testing.T) *http.Request {
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		t.Errorf("Request failure :%v", err.Error())
	}
	return request
}

func setupDb(t *testing.T) *sql.DB {
	db, err := sql.Open("mysql", "root:@/entryTask")
	if err != nil {
		t.Errorf("Database connection failed: %v", err.Error())
	}
	return db
}

func getUser(h http.Header) *protocol.User {
	user := &protocol.User{}
	proto.Unmarshal([]uint8(h["User"][0]), user)
	return user
}

func form_setup(body string, t *testing.T, db *sql.DB, url string) (*httptest.ResponseRecorder, *http.Request) {
	response := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Errorf("Request error %s.", err)
	}
	return response, request
}

func login_setup(body string, t *testing.T, db *sql.DB) *httptest.ResponseRecorder {
	response, request := form_setup(body, t, db, "/login")
	http.HandlerFunc(makeHandler(db, handlers.LoginHandler)).ServeHTTP(response, request)
	return response
}

func test_single_user_login(t *testing.T, i int, db *sql.DB) {
	response := login_setup(fmt.Sprintf("name=user%d&password=pass%d%d", i, i*2, i*2), t, db)
	header := response.Header()
	user := getUser(header)
	if user.Name != fmt.Sprintf("user%d", i) || header["Status"][0] != "successful login" {
		t.Errorf("Login to user%d unsuccessful.", i)
	} else if user.Nickname != fmt.Sprintf("nick%d", i) {
		t.Errorf("Wrong information retrieved.")
	}
}

func test_invalid_user_login(t *testing.T, info string, db *sql.DB, expectedErr string, errString string) {
	response := login_setup(info, t, db)
	header := response.Header()
	if header["Status"][0] != expectedErr {
		t.Errorf(errString)
	}
}

func updateCookie(cookieString string, response *httptest.ResponseRecorder, request *http.Request) {
	cookie := &http.Cookie{
		Name:  "session",
		Value: cookieString,
		Path:  "/",
	}
	http.SetCookie(response, cookie)
	request.Header["Cookie"] = response.HeaderMap["Set-Cookie"]
}

func test_restricted_template_resulting_header(t *testing.T, db *sql.DB, cookieString string, url string, fn func(*sql.DB, http.ResponseWriter, *http.Request)) http.Header {
	response := httptest.NewRecorder()
	request := makeRequest(http.MethodGet, url, t)
	updateCookie(cookieString, response, request)
	http.HandlerFunc(makeHandler(db, fn)).ServeHTTP(response, request)
	header := response.Header()
	return header
}

func test_restricted_template_granted_access(t *testing.T, db *sql.DB, i int, url string, fn func(*sql.DB, http.ResponseWriter, *http.Request)) {
	/*response := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(fmt.Sprintf("name=user%d&password=pass%d%d", i, i*2, i*2)))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Errorf("Request error %s.", err)
	}
	http.HandlerFunc(makeHandler(db, handlers.LoginHandler)).ServeHTTP(response, request)
	header := response.Header()
	user := getUser(header)
	if user.Name != fmt.Sprintf("user%d", i) || header["Status"][0] != "successful login" {
		t.Errorf("Login to user%d unsuccessful.", i)
	} else if user.Nickname != fmt.Sprintf("nick%d", i) {
		t.Errorf("Wrong information retrieved.")
	}*/

	cookieString := handlers.SetSessionInfo(
		&protocol.User{
			Name:     fmt.Sprintf("user%d", i),
			Password: fmt.Sprintf("pass%d%d", i*2, i*2),
		},
		&protocol.User{},
		handlers.InfoErr{},
		"",
	)
	header := test_restricted_template_resulting_header(t, db, cookieString, url, fn)
	if header["Status"][0] != "successful view" {
		t.Errorf("Failed access to restricted page for user%d", i)
	}
}

func test_restricted_template_no_access(t *testing.T, db *sql.DB, url string, fn func(*sql.DB, http.ResponseWriter, *http.Request)) {
	cookieString := handlers.SetSessionInfo(
		&protocol.User{},
		&protocol.User{},
		handlers.InfoErr{},
		"",
	)
	header := test_restricted_template_resulting_header(t, db, cookieString, url, fn)
	if header["Status"][0] != "login error" {
		t.Errorf("Wrongly granted access.")
	}
}

func test_valid_signup(t *testing.T, db *sql.DB, i int) {
	handlers.ExecuteQuery(db, "DELETE FROM users WHERE username LIKE 'test%'")
	response, request := form_setup(fmt.Sprintf("name=testuser%d&password=testpass%d%d&password_repeat=testpass%d%d", i, i*2, i*2, i*2, i*2), t, db, "/signup")
	http.HandlerFunc(makeHandler(db, handlers.SignupHandler)).ServeHTTP(response, request)
	header := response.Header()
	user := getUser(header)
	if user.Name != fmt.Sprintf("testuser%d", i) || header["Status"][0] != "successful signup" {
		t.Errorf("Sign up for testuser%d unssucessful.", i)
	}
}

func test_invalid_signup(t *testing.T, db *sql.DB, name string, pass string, repeat string, status string, errString string) {
	response, request := form_setup(fmt.Sprintf("name=%s&password=%s&password_repeat=%s", name, pass, repeat), t, db, "/signup")
	t.Log("here here")
	t.Log(request.FormValue("name"))
	t.Log(request.ParseForm())
	http.HandlerFunc(makeHandler(db, handlers.SignupHandler)).ServeHTTP(response, request)
	header := response.Header()
	if header["Status"][0] != status {
		t.Errorf(errString)
	}
}

func delete_execute(t *testing.T, db *sql.DB, name string, pass string) http.Header {
	cookieString := handlers.SetSessionInfo(
		&protocol.User{
			Name:     name,
			Password: pass,
		},
		&protocol.User{},
		handlers.InfoErr{},
		"",
	)
	response := httptest.NewRecorder()
	request := makeRequest(http.MethodPost, "/delete", t)
	updateCookie(cookieString, response, request)
	http.HandlerFunc(makeHandler(db, handlers.DeleteHandler)).ServeHTTP(response, request)
	return response.Header()
}

func test_valid_delete(t *testing.T, db *sql.DB, i int) {
	name := fmt.Sprintf("testuser%d", i)
	pass := fmt.Sprintf("pass%d%d", i*2, i*2)
	handlers.ExecuteQuery(db, "INSERT INTO users VALUES (?, ?, NULL, NULL) ON DUPLICATE KEY UPDATE username = ?", name, pass, name)
	header := delete_execute(t, db, name, pass)
	if header["Status"][0] != fmt.Sprintf("delete %s", name) {
		t.Errorf("Deletion of %s failed.", name)
	} else {
		user := &protocol.User{}
		flag := protocol.IsExistingUsername(db, name, user)
		if flag {
			t.Errorf("Deletion not taking expected effect.")
		}
	}
}

func test_valid_edit(t *testing.T, db *sql.DB, i int) {
	name := fmt.Sprintf("user%d", i)
	pass := fmt.Sprintf("pass%d%d", i*2, i*2)
	user := &protocol.User{}
	nickname := fmt.Sprintf("nick%d", i)
	nicknew := fmt.Sprintf("mick%d", i)
	response, request := form_setup(fmt.Sprintf("nickname=%s", nicknew), t, db, "/edit")
	cookieString := handlers.SetSessionInfo(
		&protocol.User{
			Name:     name,
			Password: pass,
			Nickname: nickname,
		},
		&protocol.User{
			Name:     name,
			Password: pass,
			Nickname: nickname,
		},
		handlers.InfoErr{},
		"assets/placeholder.jpeg",
	)
	updateCookie(cookieString, response, request)
	http.HandlerFunc(makeHandler(db, handlers.EditHandler)).ServeHTTP(response, request)
	flag := protocol.IsExistingUsername(db, name, user)
	if !flag {
		t.Errorf("Wrongly deleted/updated primary key of %s.", name)
	} else if user.Nickname != nicknew {
		t.Errorf("Update of %s failed.", name)
	}
}

func test_upload(t *testing.T, db *sql.DB, i int) {
	
}

func clearEffects(db *sql.DB) {
	handlers.ExecuteQuery(db, "DELETE FROM users WHERE username LIKE 'test%'")
	for i := 0; i < 5; i++ {
		handlers.ExecuteQuery(db, "UPDATE users SET photo = ?, nickname = ? WHERE username = ?", nil, fmt.Sprintf("nick%d", i), fmt.Sprintf("user%d", i))
	}
}

func Test_handlers(t *testing.T) {
	db := setupDb(t)

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

	t.Run("Upload", func(t *testing.T) {

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
