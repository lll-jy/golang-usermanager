package test

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"git.garena.com/jiayu.li/entry-task/cmd/handlers"
	_ "github.com/go-sql-driver/mysql"
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

func test_single_user(t *testing.T, i int, db *sql.DB) {
	response := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(fmt.Sprintf("name=user%d&password=pass%d%d", i, i*2, i*2)))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Errorf("Request error %s.", err)
	}
	http.HandlerFunc(makeHandler(db, handlers.LoginHandler)).ServeHTTP(response, request)
	header := response.Header()
	if header["User"][0] != fmt.Sprintf("user%d", i) || header["Status"][0] != "successful login" {
		t.Errorf("Login to user%d unsuccessful.", i)
	}
	//fmt.Printf("Loggin: %s\n", response.Header()["Status"])
	//fmt.Printf("Check %t", response.Header()["Status"][0] == "successful login")
	// fmt.Printf("Result is\n%s\n", response.Body.String())
}

func Test_server(t *testing.T) {
	t.Run("Test run", func(t *testing.T) {
		db, err := sql.Open("mysql", "root:@/entryTask")
		if err != nil {
			t.Errorf("Database connection failed: %v", err.Error())
		}

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
	})
}

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
