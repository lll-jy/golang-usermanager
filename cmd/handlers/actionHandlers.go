package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"git.garena.com/jiayu.li/entry-task/cmd/paths"
	"git.garena.com/jiayu.li/entry-task/cmd/protocol"
	"google.golang.org/protobuf/proto"

	"golang.org/x/crypto/bcrypt"
)

// login handler

func DecryptPhoto(url string, pass string, name string, photo *string) error {
	if url == "" || url == paths.PlaceholderPath {
		*photo = paths.PlaceholderPath
	} else {
		encrypted, err := ioutil.ReadFile(url)
		if err != nil {
			return errors.New(fmt.Sprintf("The encrypted file %s is invalid because %s.", url, err.Error()))
		}
		decrypted := decrypt(encrypted, pass)
		*photo = fmt.Sprintf("%s/user%s.jpeg", paths.TempPath, name)
		err = ioutil.WriteFile(*photo, decrypted, 0600)
		if err != nil {
			return errors.New(fmt.Sprintf("Cannot write file. %s.", err.Error()))
		}
	}
	return nil
}

func LoginHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	header.Set("user", "")
	header.Set("status", "")
	name := r.FormValue("name")
	pass := r.FormValue("password")
	redirectTarget := "/"
	u := protocol.User{}
	tu := createUser(name, pass)
	ie := InfoErr{}
	photo := ""
	if protocol.IsExistingUsername(db, name, &u) {
		log.Printf("User %s found.", name)
		if protocol.IsCorrectPassword(pass, u.Password) {
			log.Printf("Login to %s successful!", name)
			u.Name = name
			DecryptPhoto(u.PhotoUrl, pass, name, &photo)
			tu.PhotoUrl = u.PhotoUrl
			tu.Nickname = u.Nickname
			redirectTarget = "/view"
			user, err := proto.Marshal(&u)
			if err != nil {
				log.Printf("Error: wrong format! %v cannot be parsed as a user.", &u)
			}
			header.Set("user", string(user))
			header.Set("status", "successful login")
		} else {
			log.Printf("Login to %s unsuccessful due to wrong password!", name)
			tu.Password = ""
			ie.PasswordErr = "Incorrect password."
			header.Set("status", "incorrect password")
		}
	} else {
		log.Printf("User %s does not exists. Redirect to sign up page.", name)
		redirectTarget = "/signup"
		header.Set("status", "user not exist")
	}
	setSession(&u, &tu, ie, photo, w)
	go http.Redirect(w, r, redirectTarget, 302)
}

// logout handler

func LogoutHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	info := GetPageInfo(r)
	clearSession(w)
	log.Printf("User %s logged out.", info.User.Name)
	if info.Photo != "" && info.Photo != paths.PlaceholderPath {
		os.Remove(info.Photo)
	}
	http.Redirect(w, r, "/", 302)
}

// sign up handler

func userInfoHandler(db *sql.DB, w http.ResponseWriter, r *http.Request, rt string, tgt string, query string) {
	header := w.Header()
	header.Set("user", "")
	header.Set("status", "")
	name := r.FormValue("name")
	pass := r.FormValue("password")
	repeatPass := r.FormValue("password_repeat")
	redirectTarget := rt
	info := GetPageInfo(r)
	u := info.User
	tu := createUser(name, pass)
	ie := InfoErr{}
	if protocol.IsValidUsername(name) {
		if protocol.IsExistingUsername(db, name, u) {
			log.Printf("User signup failure: duplicate user %s found.", name)
			ie.NameErr = fmt.Sprintf("The username %s already exists.", name)
			header.Set("status", "user already exists")
		} else if protocol.IsValidPassword(pass) {
			if pass == repeatPass {
				if rt == "/signup" {
					log.Printf("New user %s signed up.", name)
				}
				hashed, err := bcrypt.GenerateFromPassword([]byte(pass), 3)
				if err != nil {
					log.Printf("Error: password %s cannot be hashed.", pass)
				}
				ExecuteQuery(db, query, name, hashed, tu.Name)
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
				DecryptPhoto(u.PhotoUrl, pass, name, &info.Photo)
				tu.Nickname = u.Nickname
				redirectTarget = tgt
				user, err := proto.Marshal(u)
				if err != nil {
					log.Printf("Error: wrong format! %v cannot be parsed as a user.", &u)
				}
				header.Set("user", string(user))
				header.Set("status", "successful signup")
			} else {
				log.Printf("User signup/reset failure: password does not match.")
				ie.PasswordRepeatErr = "The password does not match."
				header.Set("status", "mismatch password")
			}
		} else {
			log.Printf("User signup/reset failure: password format invalid.")
			u.Name = name
			ie.PasswordErr = "The password is not valid."
			header.Set("status", "wrong password format")
		}
	} else {
		log.Printf("User signup/reset failture: invalid username format of %s.", name)
		ie.NameErr = "The username format is not valid."
		header.Set("status", "wrong username format")
	}
	setSession(u, &tu, ie, info.Photo, w)
	http.Redirect(w, r, redirectTarget, 302)
}

func SignupHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	userInfoHandler(db, w, r, "/signup", "/edit", "INSERT INTO users VALUES (?, ?, NULL, NULL) ON DUPLICATE KEY UPDATE username = ?")
}

func ResetHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	info := GetPageInfo(r)
	if info.User.Password != "" {
		userInfoHandler(db, w, r, "/reset", "/view", "UPDATE users SET username = ?, password = ? WHERE username = ?")
	} else {
		http.Redirect(w, r, "/", 302)
	}
}

// edit handler

func EditHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	info := GetPageInfo(r)
	if info.User.Password != "" {
		info.TempUser.Nickname = r.FormValue("nickname")
		if len(info.TempUser.Nickname) > 20 {
			info.InfoErr.NameErr = "Nickname too long."
			setSession(info.User, info.TempUser, info.InfoErr, info.Photo, w)
			http.Redirect(w, r, "/edit", 302)
		} else {
			photo := interface{}(info.TempUser.PhotoUrl)
			if info.Photo == paths.PlaceholderPath || info.Photo == "" {
				photo = nil
			}
			nickname := interface{}(info.TempUser.Nickname)
			if info.TempUser.Nickname == "" {
				nickname = nil
			}
			ExecuteQuery(db, "UPDATE users SET photo = ?, nickname = ? WHERE username = ?", photo, nickname, info.User.Name)
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
		}
	} else {
		http.Redirect(w, r, "/", 302)
	}
}

// https://tutorialedge.net/golang/go-file-upload-tutorial/
func UploadHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	info := GetPageInfo(r)
	header.Set("tempPhoto", info.Photo)
	header.Set("status", "")
	header.Set("photo", "")
	if info.User.Password != "" {
		r.ParseMultipartForm(10 << 20) // < 10 MB files
		file, handler, err := r.FormFile("photo_file")
		if err != nil {
			log.Println("Error retrieving file.")
			header.Set("status", fmt.Sprintf("cannot retrieve: %s ||\n", err))
		} else {
			defer file.Close()
			log.Printf("Photo %s uploaded for user %s. The file size is %+v. MIME header is %+v.", handler.Filename, info.User.Name, handler.Size, handler.Header)
			targetDir := paths.FileBasePath
			tempFile, err := ioutil.TempFile(targetDir, "upload-*.jpeg")
			if err != nil {
				log.Printf("Error generating temporary file. %s.", err.Error())
			}
			defer tempFile.Close()
			fileBytes, err := ioutil.ReadAll(file)
			if err != nil {
				log.Println("Error reading file.")
			}
			tempFile.Write(encrypt(fileBytes, info.TempUser.Password))
			dirs := strings.Split(tempFile.Name(), "/")
			info.TempUser.PhotoUrl = fmt.Sprintf("%s/%s", paths.FileBaseRelativePath, dirs[len(dirs)-1])
			DecryptPhoto(info.TempUser.PhotoUrl, info.TempUser.Password, info.TempUser.Name, &info.Photo)
			setSession(info.User, info.TempUser, info.InfoErr, info.Photo, w)
			log.Println("Successfully uploaded file")
			header.Set("tempPhoto", info.Photo)
			header.Set("status", "success")
			header.Set("photo", info.TempUser.PhotoUrl)
		}
		http.Redirect(w, r, "/edit", 302)
	} else {
		http.Redirect(w, r, "/", 302)
	}
}

// discard handler

func DiscardHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	info := GetPageInfo(r)
	if info.User.Password != "" {
		if info.User.PhotoUrl != "" && info.User.PhotoUrl != paths.PlaceholderPath {
			os.Remove(info.Photo)
			DecryptPhoto(info.User.PhotoUrl, info.TempUser.Password, info.User.Name, &info.Photo)
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

func removeFile(src string, validator string) {
	if validator != "" && validator != paths.PlaceholderPath {
		os.Remove(src)
	}
}

func RemoveHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	info := GetPageInfo(r)
	if info.User.Password != "" {
		removeFile(info.Photo, info.Photo)
		setSession(info.User, info.TempUser, info.InfoErr, paths.PlaceholderPath, w)
		ExecuteQuery(db, "UPDATE users SET photo = NULL WHERE username = ?", info.User.Name)
		http.Redirect(w, r, "/edit", 302)
		log.Printf("Removed profile photo for user %s", info.User.Name)
	} else {
		http.Redirect(w, r, "/", 302)
	}
}

// delete handler

func DeleteHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	header.Set("status", "")
	info := GetPageInfo(r)
	if info.User.Password != "" {
		name := info.User.Name
		query := "DELETE FROM users WHERE username = ?"
		err := ExecuteQuery(db, query, name)
		if err != nil {
			log.Printf("Query %s cannot be executed due to error: %s", query, err.Error())
			header.Set("status", "cannot delete")
		} else {
			log.Printf("User %s deleted.", name)
			header.Set("status", fmt.Sprintf("delete %s", name))
		}
		removeFile(info.Photo, info.Photo)
		removeFile(info.User.PhotoUrl, info.Photo)
		clearSession(w)
	}
	http.Redirect(w, r, "/", 302)
}
