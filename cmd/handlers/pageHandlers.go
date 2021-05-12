package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

// templates that the web app uses
var templates *template.Template

// PrepareTemplates parses the template files and store in templates variable
func PrepareTemplates(templateFileNameFormat string) {
	templates = template.Must(template.ParseFiles(
		fmt.Sprintf(templateFileNameFormat, "index"),
		fmt.Sprintf(templateFileNameFormat, "view"),
		fmt.Sprintf(templateFileNameFormat, "signup"),
		fmt.Sprintf(templateFileNameFormat, "profile"),
	))
}

// renderTemplate renders the template for given template with given info and write to the response writer
func renderTemplate(w http.ResponseWriter, tmpl string, info *PageInfo) {
	err := templates.ExecuteTemplate(w, tmpl+".html", info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// renderRestrictedTemplate renders the template with restricted access with given info.
func renderRestrictedTemplate(w http.ResponseWriter, r *http.Request, tmpl string, fn func(*PageInfo)) {
	info := GetPageInfo(r)
	header := w.Header()

	if info.User.Password != "" {
		fn(&info)
		renderTemplate(w, tmpl, &info)
		header.Set("status", "successful view")
		log.Printf("Opened %s page.", tmpl)
	} else {
		http.Redirect(w, r, "/", 302)
		header.Set("status", "login error")
		log.Printf("Access denied. Redirect to homepage.")
	}
}

func IndexPageHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	info := GetPageInfo(r)
	renderTemplate(w, "index", &info)
}

func SignupPageHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	info := GetPageInfo(r)
	info.Action = "/signup"
	info.Title = "Sign up"
	info.CancelAction = "/"
	renderTemplate(w, "signup", &info)
}

func ViewPageHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	renderRestrictedTemplate(w, r, "view", func(info *PageInfo) {
		if info.User.Nickname != "" {
			info.DisplayName = info.User.Nickname
		} else {
			info.DisplayName = fmt.Sprintf("user %s", info.User.Name)
		}
	})
}

func EditPageHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	renderRestrictedTemplate(w, r, "profile", func(info *PageInfo) {})
}

func ResetPageHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	renderRestrictedTemplate(w, r, "signup", func(info *PageInfo) {
		info.Action = "/reset"
		info.Title = "Reset"
		info.CancelAction = "/view"
		info.User.Password = ""
	})
}
