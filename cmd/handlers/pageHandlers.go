package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"
)

// templates
var templates = template.Must(template.ParseFiles(
	"templates/index.html",
	"templates/view.html",
	"templates/signup.html",
	"templates/profile.html",
))

func renderTemplate(w http.ResponseWriter, tmpl string, info *PageInfo) {
	err := templates.ExecuteTemplate(w, tmpl+".html", info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderRestrictedTemplate(w http.ResponseWriter, r *http.Request, tmpl string, fn func(*PageInfo)) {
	info := getPageInfo(r)
	if info.User.Password != "" {
		fn(&info)
		renderTemplate(w, tmpl, &info)
	} else {
		http.Redirect(w, r, "/", 302)
	}
}

// index page

func IndexPageHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	info := getPageInfo(r)
	renderTemplate(w, "index", &info)
}

func SignupPageHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	info := getPageInfo(r)
	info.Action = "/signup"
	info.Title = "Sign up"
	info.CancelAction = "/"
	renderTemplate(w, "signup", &info)
}

// view page

func setDisplayName(info *PageInfo) {
	if info.User.Nickname != "" {
		info.DisplayName = info.User.Nickname
	} else {
		info.DisplayName = fmt.Sprintf("user %s", info.User.Name)
	}
}

func ViewPageHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	renderRestrictedTemplate(w, r, "view", setDisplayName)
}

// edit page

func EditPageHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	renderRestrictedTemplate(w, r, "profile", func(info *PageInfo) {})
}

// reset page

func ResetPageHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	renderRestrictedTemplate(w, r, "signup", func(info *PageInfo) {
		info.Action = "/reset"
		info.Title = "Reset"
		info.CancelAction = "/view"
		info.User.Password = ""
	})
}
