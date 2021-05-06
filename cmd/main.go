package main

// https://gist.github.com/mschoebel/9398202
// https://golang.org/doc/articles/wiki/

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

// cookie handling

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32),
)

type User struct {
	Name     string
	Password string
	PhotoUrl string
	Nickname string
}

type InfoErr struct {
	UsernameErr string
	PasswordErr string
}

type PageInfo struct {
	User    User
	InfoErr InfoErr
}

func getPageInfo(r *http.Request) (info PageInfo) {
	var username string
	var password string
	var nameErr string
	var passErr string
	if cookie, err := r.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			username = cookieValue["name"]
			password = cookieValue["pass"]
			nameErr = cookieValue["nameErr"]
			passErr = cookieValue["passErr"]
		}
	}
	u := User{Name: username, Password: password}
	ie := InfoErr{UsernameErr: nameErr, PasswordErr: passErr}
	return PageInfo{
		User:    u,
		InfoErr: ie,
	}
}

func setSession(u User, uie InfoErr, w http.ResponseWriter) {
	value := map[string]string{
		"name":    u.Name,
		"pass":    u.Password,
		"nameErr": uie.UsernameErr,
		"passErr": uie.PasswordErr,
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

func renderUserinfoTemplate(w http.ResponseWriter, tmpl string, info PageInfo) {
	err := templates.ExecuteTemplate(w, tmpl+".html", info)
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
	// .. check credentials .. TODO
	u := User{Name: name, Password: pass}
	ie := InfoErr{}
	if isExistingUsername(name) {
		log.Printf("User %s found.", name)
		if isCorrectPassword(name, pass) {
			log.Printf("Login to %s successful!", name)
			redirectTarget = "/view"
		} else {
			log.Printf("Login to %s unsuccessful due to wrong password!", name)
			u.Password = ""
			ie.PasswordErr = "Incorrect password." // TODO err msg
		}
	} else {
		log.Printf("User %s does not exists. Redirect to sign up page.", name)
	}
	setSession(u, ie, w)
	http.Redirect(w, r, redirectTarget, 302)
}

// logout handler

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	clearSession(w)
	http.Redirect(w, r, "/", 302)
}

// index page

func indexPageHandler(w http.ResponseWriter, r *http.Request) {
	info := getPageInfo(r)
	renderUserinfoTemplate(w, "index", info)
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
	info := getPageInfo(r)
	username := info.User.Name
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
	router.HandleFunc("/view", internalPageHandler)

	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/logout", logoutHandler).Methods("POST")

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
