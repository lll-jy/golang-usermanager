package main

import (
	"fmt"
	"log"
	"net/http"
)

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
	setSession(&u, ie, w)
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
	setSession(&u, ie, w)
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

// delete handler

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	info := getPageInfo(r)
	if info.User.Password != "" {
		// TODO delete user
		clearSession(w)
	}
	http.Redirect(w, r, "/", 302)
}
