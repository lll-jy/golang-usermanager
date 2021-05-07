package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// login handler

func loginHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	pass := r.FormValue("password")
	redirectTarget := "/"
	u := createUser("", "")
	ie := InfoErr{}
	if isExistingUsername(name, &u) {
		log.Printf("User %s found.", name)
		if isCorrectPassword(pass, u.Password) {
			log.Printf("Login to %s successful!", name)
			u.Password = "correct"
			redirectTarget = "/view"
		} else {
			log.Printf("Login to %s unsuccessful due to wrong password!", name)
			ie.PasswordErr = "Incorrect password."
		}
	} else {
		log.Printf("User %s does not exists. Redirect to sign up page.", name)
		u.Password = pass
		redirectTarget = "/signup"
	}
	u.Name = name
	setSession(&u, ie, "", w)
	http.Redirect(w, r, redirectTarget, 302)
}

// logout handler

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	clearSession(w)
	log.Printf("User %s logged out.", getPageInfo(r).User.Name)
	http.Redirect(w, r, "/", 302)
}

// sign up handler

func userInfoHandler(w http.ResponseWriter, r *http.Request, rt string, tgt string, query string) {
	name := r.FormValue("name")
	pass := r.FormValue("password")
	repeatPass := r.FormValue("password_repeat")
	redirectTarget := rt
	u := createUser(getPageInfo(r).User.Name, pass)
	ie := InfoErr{}
	if isValidUsername(name) {
		if isExistingUsername(name, &u) {
			log.Printf("User signup failure: duplicate user %s found.", name)
			u = createUser(name, pass)
			ie.UsernameErr = fmt.Sprintf("The username %s already exists.", name)
		} else if isValidPassword(pass) {
			if pass == repeatPass {
				log.Printf("New user %s signed up.", name)
				hashed, err := bcrypt.GenerateFromPassword([]byte(pass), 3)
				if err != nil {
					log.Printf("Error: password %s cannot be hashed.", pass)
				}
				executeQuery(db, query, name, hashed, u.Name)
				u.Name = name
				u.Password = "correct"
				redirectTarget = tgt
			} else {
				log.Printf("User signup failure: password does not match.")
				u.Name = name
				u.Password = pass
				ie.PasswordRepeatErr = "The password does not match."
			}
		} else {
			log.Printf("User signup failure: password format invalid.")
			u.Name = name
			u.Password = ""
			ie.PasswordErr = "The password is not valid."
		}
	} else {
		log.Printf("User signup failture: invalid username format of %s.", name)
		ie.UsernameErr = "The username format is not valid."
	}
	setSession(&u, ie, "", w)
	http.Redirect(w, r, redirectTarget, 302)
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	userInfoHandler(w, r, "/signup", "/edit", "INSERT INTO users VALUES (?, ?, NULL, NULL) ON DUPLICATE KEY UPDATE username = ?")
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
	info := getPageInfo(r)
	if info.User.Password != "" {
		userInfoHandler(w, r, "/reset", "/view", "UPDATE users SET username = ?, password = ? WHERE username = ?")
	} else {
		http.Redirect(w, r, "/", 302)
	}
}

// edit handler

func editHandler(w http.ResponseWriter, r *http.Request) {
	info := getPageInfo(r)
	if info.User.Password != "" {
		info.User.Nickname = r.FormValue("nickname")
		photo := interface{}(info.User.PhotoUrl)
		if info.User.PhotoUrl == "assets/placeholder.jpeg" || info.User.PhotoUrl == "" {
			photo = nil
		}
		nickname := interface{}(info.User.Nickname)
		if info.User.Nickname == "" {
			nickname = nil
		}
		executeQuery(db, "UPDATE users SET photo = ?, nickname = ? WHERE username = ?", photo, nickname, info.User.Name)
		log.Printf("User information of %s updated.", info.User.Name)
		log.Printf("Initial photo at %s", info.InitialPhoto)
		if info.InitialPhoto != "" && info.InitialPhoto != "/assets/placeholder.jpeg" {
			err := os.Remove(info.InitialPhoto)
			if err == nil {
				log.Printf("Removed original photo from database.")
			} else {
				log.Printf(err.Error())
			}
		}
		setSession(info.User, info.InfoErr, "", w)
		http.Redirect(w, r, "/view", 302)
	} else {
		http.Redirect(w, r, "/", 302)
	}
}

// https://tutorialedge.net/golang/go-file-upload-tutorial/
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	info := getPageInfo(r)
	initial := info.User.PhotoUrl
	if info.User.Password != "" {
		r.ParseMultipartForm(10 << 20) // < 10 MB files
		file, handler, err := r.FormFile("photo_file")
		if err != nil {
			log.Println("Error retrieving file.")
		} else {
			defer file.Close()
			log.Printf("Photo %s uploaded for user %s. The file size is %+v. MIME header is %+v.", handler.Filename, info.User.Name, handler.Size, handler.Header)
			targetDir := "../../../Desktop/EntryTask/entry-task/test/data/upload" // EXTEND: May set to some cloud space
			tempFile, err := ioutil.TempFile(targetDir, "upload-*.jpeg")
			if err != nil {
				log.Println("Error generating temporary file.")
				log.Println(err)
			}
			defer tempFile.Close()
			fileBytes, err := ioutil.ReadAll(file)
			if err != nil {
				log.Println("Error reading file.")
			}
			tempFile.Write(fileBytes)
			dirs := strings.Split(tempFile.Name(), "/")
			info.User.PhotoUrl = fmt.Sprintf("test/data/upload/%s", dirs[len(dirs)-1]) // EXTEND: same as above
			setSession(info.User, info.InfoErr, initial, w)
			log.Println("Successfully uploaded file")
		}
		http.Redirect(w, r, "/edit", 302)
	} else {
		http.Redirect(w, r, "/", 302)
	}
}

// delete handler

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	info := getPageInfo(r)
	if info.User.Password != "" {
		name := info.User.Name
		executeQuery(db, "DELETE FROM users WHERE username = ?", name)
		log.Printf("User %s deleted.", name)
		clearSession(w)
	}
	http.Redirect(w, r, "/", 302)
}
