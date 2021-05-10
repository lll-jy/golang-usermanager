package test

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.garena.com/jiayu.li/entry-task/cmd/handlers"
	"git.garena.com/jiayu.li/entry-task/cmd/protocol"
	_ "github.com/go-sql-driver/mysql"
)

func makeHandler(db *sql.DB, fn func(*sql.DB, http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(db, w, r)
	}
}

func Test_server(t *testing.T) {
	t.Run("Test run", func(t *testing.T) {
		db, err := sql.Open("mysql", "root:@/entryTask")
		if err != nil {
			t.Errorf("Database connection failed: %v", err.Error())
		}
		response := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Errorf("Request failure :%v", err.Error())
		}
		makeHandler(db, handlers.IndexPageHandler)(response, request)
		// t.Logf("Index body is %s", response.Body.String())
		// fmt.Printf("Index body is %s\n", response.Body.String())
		makeHandler(db, handlers.LoginHandler)(response, request)
		// fmt.Printf("1 Login Body is %s\n", response.Body.String())

		cookie := &http.Cookie{
			Name: "session",
			Value: handlers.SetSessionInfo(
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
			),
			Path: "/",
		}
		http.SetCookie(response, cookie)
		request.Header = http.Header{"Cookie": response.HeaderMap["Set-Cookie"]}
		makeHandler(db, handlers.LoginHandler)(response, request)
		// fmt.Printf("2 Login Body is %s\n", response.Body.String())
		cookie, err = request.Cookie("session")
		fmt.Printf("%v\n", cookie)
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
