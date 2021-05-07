package main

import (
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

func renderTemplate(w http.ResponseWriter, tmpl string, info PageInfo) {
	err := templates.ExecuteTemplate(w, tmpl+".html", info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
