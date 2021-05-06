package main

// https://gist.github.com/mschoebel/9398202
// https://golang.org/doc/articles/wiki/

import (
	"fmt"
	"net/http"
	"regexp"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

// cookie handling

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

type User struct {
	Name     string
	Password string
	PhotoUrl string
	Nickname string
}

func getUser(r *http.Request) (u User) {
	var username string
	var password string
	if cookie, err := r.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			username = cookieValue["name"]
			password = cookieValue["pass"]
		}
	}
	return User{Name: username, Password: password}
}

func setSession(u User, w http.ResponseWriter) {
	value := map[string]string{
		"name": u.Name,
		"pass": u.Password,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}
}

func clearSession(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

// templates
var templates = template.Must(template.ParseFiles("templates/index.html"))

func renderTemplate(w http.ResponseWriter, tmpl string) {
	err := templates.ExecuteTemplate(w, tmpl+".html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func isExistingUsername(username string) bool { // TODO
	var validUsername = regexp.MustCompile("^u([a-z]+)$")
	return validUsername.MatchString(username)
}

func isCorrectPassword(username string, password string) bool { // TODO
	return username == password
}

// login handler

func loginHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	pass := r.FormValue("password")
	redirectTarget := "/"
	if name != "" && pass != "" {
		// .. check credentials .. TODO
		if isExistingUsername(name) {
			fmt.Println("exist!!")
			if isCorrectPassword(name, pass) {
				fmt.Println("login successful!")
			} else {
				fmt.Println("Fail~")
			}
		} else {
			fmt.Println("not exist!!")
		}
		setSession(User{Name: name, Password: pass}, w)
		redirectTarget = "/internal"
	}
	http.Redirect(w, r, redirectTarget, 302)
}

// logout handler

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	clearSession(w)
	http.Redirect(w, r, "/", 302)
}

// index page

func indexPageHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index")
}

// internal page

const internalPage = `
<h1>Internal</h1>
<hr>
<small>User: %s</small>
<form method="post" action="/logout">
    <button type="submit">Logout</button>
</form>
`

func internalPageHandler(w http.ResponseWriter, r *http.Request) {
	u := getUser(r)
	username := u.Name
	if username != "" {
		fmt.Fprintf(w, internalPage, username)
	} else {
		http.Redirect(w, r, "/", 302)
	}
}

// server main method

var router = mux.NewRouter()

func main() {

	router.HandleFunc("/", indexPageHandler)
	router.HandleFunc("/internal", internalPageHandler)

	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/logout", logoutHandler).Methods("POST")

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
