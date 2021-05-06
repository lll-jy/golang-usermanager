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

func createUser(name string, pass string) User {
	return User{
		Name:     name,
		Password: pass,
		PhotoUrl: "",
		Nickname: "",
	}
}

type InfoErr struct {
	UsernameErr       string
	PasswordErr       string
	PasswordRepeatErr string
}

type PageInfo struct {
	User         User
	InfoErr      InfoErr
	DisplayName  string
	Action       string
	Title        string
	CancelAction string
}

func getPageInfo(r *http.Request) (info PageInfo) {
	var username string
	var password string
	var photo string
	var nickname string
	var nameErr string
	var passErr string
	var repeatPassErr string
	if cookie, err := r.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			username = cookieValue["name"]
			password = cookieValue["pass"]
			photo = cookieValue["photo"]
			nickname = cookieValue["nickname"]
			nameErr = cookieValue["nameErr"]
			passErr = cookieValue["passErr"]
			repeatPassErr = cookieValue["repeatPassErr"]
		}
	}
	u := User{
		Name:     username,
		Password: password,
		PhotoUrl: photo,
		Nickname: nickname,
	}
	ie := InfoErr{
		UsernameErr:       nameErr,
		PasswordErr:       passErr,
		PasswordRepeatErr: repeatPassErr,
	}
	return PageInfo{
		User:    u,
		InfoErr: ie,
	}
}

func setSession(u User, uie InfoErr, w http.ResponseWriter) {
	value := map[string]string{
		"name":          u.Name,
		"pass":          u.Password,
		"photo":         u.PhotoUrl,
		"nickname":      u.Nickname,
		"nameErr":       uie.UsernameErr,
		"passErr":       uie.PasswordErr,
		"repeatPassErr": uie.PasswordRepeatErr,
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
var templates = template.Must(template.ParseFiles(
	"templates/index.html",
	"templates/view.html",
	"templates/signup.html",
	"templates/profile.html",
))

func renderTemplate(w http.ResponseWriter, tmpl string, info PageInfo) {
	err := templates.ExecuteTemplate(w, tmpl+".html", info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// validation helper functions

func isValidUsername(username string) bool { // TODO
	var validUsername = regexp.MustCompile("^u([a-z]+)$")
	return validUsername.MatchString(username)
}

func isValidPassword(password string) bool { // TODO
	var validPassword = regexp.MustCompile("^uu([a-z]+)$")
	return validPassword.MatchString(password)
}

func isExistingUsername(username string) bool { // TODO
	var validUsername = regexp.MustCompile("^uu([a-z]+)$")
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
	u := createUser(name, pass)
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
		redirectTarget = "/signup"
	}
	setSession(u, ie, w)
	http.Redirect(w, r, redirectTarget, 302)
}

// logout handler

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	clearSession(w)
	http.Redirect(w, r, "/", 302)
}

// sign up handler

func userInfoHandler(w http.ResponseWriter, r *http.Request, rt string, tgt string) {
	name := r.FormValue("name")
	pass := r.FormValue("password")
	repeatPass := r.FormValue("password_repeat")
	redirectTarget := rt
	u := createUser(name, pass)
	ie := InfoErr{}
	if isValidUsername(name) {
		if isExistingUsername(name) {
			log.Printf("User signup failure: duplicate user %s found.", name)
			ie.UsernameErr = fmt.Sprintf("The username %s already exists.", name)
		} else if isValidPassword(pass) {
			if pass == repeatPass {
				log.Printf("New user %s signed up.", name)
				redirectTarget = tgt // TODO: insert
			} else {
				log.Printf("User signup failure: password does not match.")
				ie.PasswordRepeatErr = "The password does not match."
			}
		} else {
			log.Printf("User signup failure: password format invalid.")
			ie.PasswordErr = "The password is not valid." // TODO requirement
		}
	} else {
		log.Printf("User signup failture: invalid username format of %s.", name)
		ie.UsernameErr = "The username format is not valid." // TODO requirement
	}
	setSession(u, ie, w)
	http.Redirect(w, r, redirectTarget, 302)
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	userInfoHandler(w, r, "/signup", "/edit")
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
	info := getPageInfo(r)
	if info.User.Password != "" {
		userInfoHandler(w, r, "/reset", "view")
	} else {
		http.Redirect(w, r, "/", 302)
	}
}

// edit handler

func editHandler(w http.ResponseWriter, r *http.Request) {
	info := getPageInfo(r)
	if info.User.Password != "" {
		info.User.PhotoUrl = r.FormValue("photo")
		info.User.Nickname = r.FormValue("nickname") // TODO
		setSession(info.User, info.InfoErr, w)
		http.Redirect(w, r, "/view", 302)
	} else {
		http.Redirect(w, r, "/", 302)
	}
}

// index page

func indexPageHandler(w http.ResponseWriter, r *http.Request) {
	info := getPageInfo(r)
	renderTemplate(w, "index", info)
}

func signupPageHandler(w http.ResponseWriter, r *http.Request) {
	info := getPageInfo(r)
	info.Action = "/signup"
	info.Title = "Sign up"
	info.CancelAction = "/"
	renderTemplate(w, "signup", info)
}

// view page

func setDisplayName(info *PageInfo) {
	if info.User.Nickname != "" {
		info.DisplayName = info.User.Nickname
	} else {
		info.DisplayName = fmt.Sprintf("user %s", info.User.Name)
	}
}

func viewPageHandler(w http.ResponseWriter, r *http.Request) {
	info := getPageInfo(r)
	if info.User.Password != "" {
		setDisplayName(&info)
		renderTemplate(w, "view", info)
	} else {
		http.Redirect(w, r, "/", 302)
	}
}

// edit page

func editPageHandler(w http.ResponseWriter, r *http.Request) {
	info := getPageInfo(r)
	if info.User.Password != "" {
		renderTemplate(w, "profile", info)
	} else {
		http.Redirect(w, r, "/", 302)
	}
}

// reset page

func resetPageHandler(w http.ResponseWriter, r *http.Request) {
	info := getPageInfo(r)
	if info.User.Password != "" {
		info.Action = "/reset"
		info.Title = "Reset"
		info.CancelAction = "/view"
		renderTemplate(w, "signup", info)
	} else {
		http.Redirect(w, r, "/", 302)
	}
}

// server main method

var router = mux.NewRouter()

func main() {

	router.HandleFunc("/", indexPageHandler)
	router.HandleFunc("/view", viewPageHandler)
	router.HandleFunc("/signup", signupPageHandler).Methods("GET")
	router.HandleFunc("/edit", editPageHandler).Methods("GET")
	router.HandleFunc("/reset", resetPageHandler).Methods("GET")

	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/logout", logoutHandler).Methods("POST")
	router.HandleFunc("/signup", signupHandler).Methods("POST")
	router.HandleFunc("/edit", editHandler).Methods("POST")
	router.HandleFunc("/reset", resetHandler).Methods("POST")

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
