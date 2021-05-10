package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"git.garena.com/jiayu.li/entry-task/cmd/paths"
	"git.garena.com/jiayu.li/entry-task/cmd/protocol"

	//"./protocol"
	"golang.org/x/crypto/bcrypt"
)

// login handler

func decryptPhoto(url string, pass string, name string, photo *string) {
	if url == "" || url == paths.PlaceholderPath {
		*photo = paths.PlaceholderPath
	} else {
		encrypted, err := ioutil.ReadFile(url)
		if err != nil {
			log.Printf("The encrypted file is invalid.")
		}
		decrypted := decrypt(encrypted, pass)
		*photo = fmt.Sprintf("%s/user%s.jpeg", paths.TempPath, name)
		ioutil.WriteFile(*photo, decrypted, 0600)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	pass := r.FormValue("password")
	redirectTarget := "/"
	u := createUser("", "")
	tu := createUser(name, pass)
	ie := InfoErr{}
	photo := ""
	if protocol.IsExistingUsername(db, name, &u) {
		log.Printf("User %s found.", name)
		if protocol.IsCorrectPassword(pass, u.Password) {
			log.Printf("Login to %s successful!", name)
			u.Name = name
			decryptPhoto(u.PhotoUrl, pass, name, &photo)
			tu.PhotoUrl = u.PhotoUrl
			tu.Nickname = u.Nickname
			redirectTarget = "/view"
		} else {
			log.Printf("Login to %s unsuccessful due to wrong password!", name)
			tu.Password = ""
			ie.PasswordErr = "Incorrect password."
		}
	} else {
		log.Printf("User %s does not exists. Redirect to sign up page.", name)
		redirectTarget = "/signup"
	}
	setSession(&u, &tu, ie, photo, w)
	http.Redirect(w, r, redirectTarget, 302)
}

// logout handler

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	info := getPageInfo(r)
	clearSession(w)
	log.Printf("User %s logged out.", info.User.Name)
	if info.Photo != "" && info.Photo != paths.PlaceholderPath {
		os.Remove(info.Photo)
	}
	http.Redirect(w, r, "/", 302)
}

// sign up handler

func userInfoHandler(w http.ResponseWriter, r *http.Request, rt string, tgt string, query string) {
	name := r.FormValue("name")
	pass := r.FormValue("password")
	repeatPass := r.FormValue("password_repeat")
	redirectTarget := rt
	info := getPageInfo(r)
	u := info.User
	tu := createUser(name, pass)
	ie := InfoErr{}
	if protocol.IsValidUsername(name) {
		if protocol.IsExistingUsername(db, name, u) {
			log.Printf("User signup failure: duplicate user %s found.", name)
			ie.UsernameErr = fmt.Sprintf("The username %s already exists.", name)
		} else if protocol.IsValidPassword(pass) {
			if pass == repeatPass {
				if rt == "/signup" {
					log.Printf("New user %s signed up.", name)
				}
				hashed, err := bcrypt.GenerateFromPassword([]byte(pass), 3)
				if err != nil {
					log.Printf("Error: password %s cannot be hashed.", pass)
				}
				executeQuery(db, query, name, hashed, tu.Name)
				if rt == "/reset" {
					u.Name = ""
					protocol.IsExistingUsername(db, name, u)
				}
				u.Name = name
				u.Password = string(hashed)
				if info.Photo != "" && info.Photo != paths.PlaceholderPath {
					os.Remove(u.PhotoUrl)
					if err != nil {
						log.Printf("The original file path is invalid.")
					}
					fileBytes, err := ioutil.ReadFile(info.Photo)
					if err != nil {
						log.Printf("The temporary file cannot be read properly.")
					}
					err = ioutil.WriteFile(u.PhotoUrl, encrypt(fileBytes, pass), 0600)
					if err != nil {
						log.Printf("The file cannot be re-encrypted.")
					}
					log.Printf("The original file key is updated.")
				}
				decryptPhoto(u.PhotoUrl, pass, name, &info.Photo)
				tu.Nickname = u.Nickname
				redirectTarget = tgt
			} else {
				log.Printf("User signup/reset failure: password does not match.")
				ie.PasswordRepeatErr = "The password does not match."
			}
		} else {
			log.Printf("User signup/reset failure: password format invalid.")
			u.Name = name
			ie.PasswordErr = "The password is not valid."
		}
	} else {
		log.Printf("User signup/reset failture: invalid username format of %s.", name)
		ie.UsernameErr = "The username format is not valid."
	}
	setSession(u, &tu, ie, info.Photo, w)
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
		info.TempUser.Nickname = r.FormValue("nickname")
		photo := interface{}(info.TempUser.PhotoUrl)
		if info.Photo == paths.PlaceholderPath || info.Photo == "" {
			photo = nil
		}
		nickname := interface{}(info.TempUser.Nickname)
		if info.TempUser.Nickname == "" {
			nickname = nil
		}
		executeQuery(db, "UPDATE users SET photo = ?, nickname = ? WHERE username = ?", photo, nickname, info.User.Name)
		log.Printf("User information of %s updated.", info.User.Name)
		if info.User.PhotoUrl != "" && info.User.PhotoUrl != paths.PlaceholderPath {
			err := os.Remove(info.User.PhotoUrl)
			if err == nil {
				log.Printf("Removed original photo from database.")
			} else {
				log.Printf(err.Error())
			}
		}
		info.User.PhotoUrl = protocol.ConvertToString(photo)
		info.User.Nickname = info.TempUser.Nickname
		setSession(info.User, info.TempUser, info.InfoErr, info.Photo, w)
		http.Redirect(w, r, "/view", 302)
	} else {
		http.Redirect(w, r, "/", 302)
	}
}

// https://tutorialedge.net/golang/go-file-upload-tutorial/
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	info := getPageInfo(r)
	if info.User.Password != "" {
		r.ParseMultipartForm(10 << 20) // < 10 MB files
		file, handler, err := r.FormFile("photo_file")
		if err != nil {
			log.Println("Error retrieving file.")
		} else {
			defer file.Close()
			log.Printf("Photo %s uploaded for user %s. The file size is %+v. MIME header is %+v.", handler.Filename, info.User.Name, handler.Size, handler.Header)
			targetDir := paths.FileBasePath
			tempFile, err := ioutil.TempFile(targetDir, "upload-*.jpeg")
			if err != nil {
				log.Println("Error generating temporary file.")
			}
			defer tempFile.Close()
			fileBytes, err := ioutil.ReadAll(file)
			if err != nil {
				log.Println("Error reading file.")
			}
			tempFile.Write(encrypt(fileBytes, info.TempUser.Password))
			dirs := strings.Split(tempFile.Name(), "/")
			info.TempUser.PhotoUrl = fmt.Sprintf("%s/%s", paths.FileBaseRelativePath, dirs[len(dirs)-1]) // EXTEND: same as above
			decryptPhoto(info.TempUser.PhotoUrl, info.TempUser.Password, info.TempUser.Name, &info.Photo)
			setSession(info.User, info.TempUser, info.InfoErr, info.Photo, w)
			log.Println("Successfully uploaded file")
		}
		http.Redirect(w, r, "/edit", 302)
	} else {
		http.Redirect(w, r, "/", 302)
	}
}

// discard handler

func discardHandler(w http.ResponseWriter, r *http.Request) {
	info := getPageInfo(r)
	if info.User.Password != "" {
		if info.User.PhotoUrl != "" && info.User.PhotoUrl != paths.PlaceholderPath {
			os.Remove(info.Photo)
			decryptPhoto(info.User.PhotoUrl, info.TempUser.Password, info.User.Name, &info.Photo)
			log.Printf("Temporary file removed.")
		} else {
			log.Printf("No file to remove.")
		}
		info.TempUser.PhotoUrl = info.User.PhotoUrl
		setSession(info.User, info.TempUser, info.InfoErr, info.Photo, w)
		http.Redirect(w, r, "/view", 302)
	} else {
		http.Redirect(w, r, "/", 302)
	}
}

// remove handler

func removeHandler(w http.ResponseWriter, r *http.Request) {
	info := getPageInfo(r)
	if info.User.Password != "" {
		if info.Photo != "" && info.Photo != paths.PlaceholderPath {
			os.Remove(info.Photo)
		}
		setSession(info.User, info.TempUser, info.InfoErr, paths.PlaceholderPath, w)
		executeQuery(db, "UPDATE users SET photo = NULL WHERE username = ?", info.User.Name)
		http.Redirect(w, r, "/edit", 302)
		log.Printf("Removed profile photo for user %s", info.User.Name)
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
		if info.Photo != "" && info.Photo != paths.PlaceholderPath {
			os.Remove(info.Photo)
		}
		if info.Photo != "" && info.Photo != paths.PlaceholderPath {
			os.Remove(info.User.PhotoUrl)
		}
		clearSession(w)
	}
	http.Redirect(w, r, "/", 302)
}
