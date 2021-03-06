// Package handlers contains all handlers and helper functions for session and photo encryption handling.
package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"git.garena.com/jiayu.li/entry-task/cmd/logging"
	"git.garena.com/jiayu.li/entry-task/cmd/paths"
	"git.garena.com/jiayu.li/entry-task/cmd/protocol"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"os"

	"google.golang.org/protobuf/proto"
)

func LoginHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	header.Set("user", "")
	header.Set("status", "")

	name := r.FormValue("name")
	pass := r.FormValue("password")

	redirectTarget := "/"
	u := protocol.User{}
	tu := createUser(name, pass)
	ie := &InfoErr{}
	photo := ""

	if protocol.IsExistingUsername(db, name, &u) {
		logging.Log(logging.INFO, fmt.Sprintf("User %s found.", name))
		if protocol.IsCorrectPassword(pass, u.Password) {
			logging.Log(logging.INFO, fmt.Sprintf("Login to %s successful!", name))
			u.Name = name
			err := DecryptPhoto(u.PhotoUrl, pass, name, &photo)
			if err != nil {
				logging.Log(logging.ERROR, fmt.Sprintf(err.Error()))
			}
			tu.PhotoUrl = u.PhotoUrl
			tu.Nickname = u.Nickname
			redirectTarget = "/view"
			user, err := proto.Marshal(&u)
			if err != nil {
				logging.Log(logging.DEBUG, fmt.Sprintf("User %v cannot be parsed as a user: %s", &u, err.Error()))
			}

			header.Set("user", string(user))
			header.Set("status", "successful login")
		} else {
			logging.Log(logging.INFO, fmt.Sprintf("Login to %s unsuccessful due to wrong password!", name))
			tu.Password = ""
			ie.PasswordErr = "Incorrect password."
			header.Set("status", "incorrect password")
		}
	} else {
		logging.Log(logging.INFO, fmt.Sprintf("User %s does not exists. Redirect to sign up page.", name))
		redirectTarget = "/signup"
		header.Set("status", "user not exist")
	}
	setSession(&u, &tu, ie, photo, w)
	http.Redirect(w, r, redirectTarget, 302)
}

// DecryptPhoto decrypts photo at url with pass as key of the user of given username and copy it to a local
// location with path stored at photo.
func DecryptPhoto(url string, pass string, name string, photo *string) error {
	if url == "" || url == paths.PlaceholderPath {
		*photo = paths.PlaceholderPath
	} else {
		encrypted, err := ioutil.ReadFile(url)
		if err != nil {
			return errors.New(fmt.Sprintf("The encrypted file %s is invalid: %s.", url, err.Error()))
		}
		decrypted := decrypt(encrypted, pass)
		*photo = fmt.Sprintf("%s/user%s.jpeg", paths.TempPath, name)
		err = ioutil.WriteFile(*photo, decrypted, 0600)
		if err != nil {
			return errors.New(fmt.Sprintf("Cannot write file: %s.", err.Error()))
		}
	}
	return nil
}

func LogoutHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	info := GetPageInfo(r)
	clearSession(w)
	logging.Log(logging.INFO, fmt.Sprintf("User %s logged out.", info.User.Name))
	if info.Photo != "" && info.Photo != paths.PlaceholderPath {
		err := os.Remove(info.Photo)
		if err != nil {
			logging.Log(logging.DEBUG, fmt.Sprintf("Removing file %s unsuccessful: %s", info.Photo, err.Error()))
		}
	}
	http.Redirect(w, r, "/", 302)
}

func SignupHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	userInfoHandler(db, w, r, "/signup", "/edit", "INSERT INTO users VALUES (?, ?, NULL, NULL) ON DUPLICATE KEY UPDATE username = ?")
}

func ResetHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	info := GetPageInfo(r)
	if info.User.Password != "" {
		userInfoHandler(db, w, r, "/reset", "/view",
			"UPDATE users SET username = ?, password = ? WHERE username = ?")
	} else {
		logging.Log(logging.INFO, fmt.Sprintf("Access denied. Redirect to homepage."))
		http.Redirect(w, r, "/", 302)
	}
}

// userInfoHandler is the shared part of signup and reset page, with rt as the route, tgt as the redirecting target,
// and query as the database SQL query to execute (insert or update).
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
	ie := &InfoErr{}
	if protocol.IsValidUsername(name) {
		if protocol.IsExistingUsername(db, name, u) {
			logging.Log(logging.INFO, fmt.Sprintf("User signup failure: duplicate user %s found.", name))
			ie.NameErr = fmt.Sprintf("The username %s already exists.", name)
			header.Set("status", "user already exists")
		} else if protocol.IsValidPassword(pass) {
			if pass == repeatPass {
				if rt == "/signup" {
					logging.Log(logging.INFO, fmt.Sprintf("New user %s signed up.", name))
				}
				hashed, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
				if err != nil {
					logging.Log(logging.ERROR, fmt.Sprintf("Password %s cannot be hashed: %s", pass,
						err.Error()))
				}
				err = ExecuteQuery(db, query, name, hashed, tu.Name)
				if err != nil {
					logging.Log(logging.DEBUG, fmt.Sprintf("Query %s cannot be executed with arguments %s, " +
						"%s, %s: %s", query, name, hashed, tu.Name, err.Error()))
				}
				if rt == "/reset" {
					u.Name = ""
					protocol.IsExistingUsername(db, name, u)
				}
				u.Name = name
				u.Password = string(hashed)
				if info.Photo != "" && info.Photo != paths.PlaceholderPath {
					err := os.Remove(u.PhotoUrl)
					if err != nil {
						logging.Log(logging.DEBUG, fmt.Sprintf("The original file path is invalid: %s", err.Error()))
					}
					fileBytes, err := ioutil.ReadFile(info.Photo)
					if err != nil {
						logging.Log(logging.ERROR, fmt.Sprintf("The temporary file cannot be read properly: %s", err.Error()))
					}
					err = ioutil.WriteFile(u.PhotoUrl, encrypt(fileBytes, pass), 0600)
					if err != nil {
						logging.Log(logging.ERROR, fmt.Sprintf("The file cannot be re-encrypted: %s", err.Error()))
					}
					logging.Log(logging.INFO, fmt.Sprintf("The original file key is updated."))
				}
				err = DecryptPhoto(u.PhotoUrl, pass, name, &info.Photo)
				if err != nil {
					logging.Log(logging.ERROR, fmt.Sprintf("Cannot decrypt and copy photo %s: %s", u.PhotoUrl, err.Error()))
				}
				tu.Nickname = u.Nickname
				redirectTarget = tgt
				user, err := proto.Marshal(u)
				if err != nil {
					logging.Log(logging.DEBUG, fmt.Sprintf("User %v cannot be parsed as a user: %s", &u, err.Error()))
				}
				header.Set("user", string(user))
				header.Set("status", "successful signup")
			} else {
				logging.Log(logging.INFO, fmt.Sprintf("User signup/reset failure: password does not match."))
				ie.PasswordRepeatErr = "The password does not match."
				header.Set("status", "mismatch password")
			}
		} else {
			logging.Log(logging.INFO, fmt.Sprintf("User signup/reset failure: password format invalid."))
			u.Name = name
			ie.PasswordErr = "The password is not valid."
			header.Set("status", "wrong password format")
		}
	} else {
		logging.Log(logging.INFO, fmt.Sprintf("User signup/reset failture: invalid username format of %s.", name))
		ie.NameErr = "The username format is not valid."
		header.Set("status", "wrong username format")
	}
	setSession(u, &tu, ie, info.Photo, w)
	http.Redirect(w, r, redirectTarget, 302)
}

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
			err := ExecuteQuery(db, "UPDATE users SET photo = ?, nickname = ? WHERE username = ?",
				photo, nickname, info.User.Name)
			if err != nil {
				logging.Log(logging.DEBUG, fmt.Sprintf("Execution of query updating user profile failed: %s", err.Error()))
			}
			logging.Log(logging.INFO, fmt.Sprintf("User information of %s updated.", info.User.Name))

			if info.User.PhotoUrl != "" && info.User.PhotoUrl != paths.PlaceholderPath {
				err := os.Remove(info.User.PhotoUrl)
				if err == nil {
					logging.Log(logging.INFO, fmt.Sprintf("Removed original photo from database."))
				} else {
					logging.Log(logging.ERROR, fmt.Sprintf("Cannot remove the file %s: %s", info.User.PhotoUrl, err.Error()))
				}
			}
			info.User.PhotoUrl = protocol.ConvertToString(photo)
			info.User.Nickname = info.TempUser.Nickname
			setSession(info.User, info.TempUser, info.InfoErr, info.Photo, w)
			http.Redirect(w, r, "/view", 302)
		}
	} else {
		logging.Log(logging.INFO, fmt.Sprintf("Access denied. Redirect to homepage."))
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
		err := r.ParseMultipartForm(10 << 20) // < 10 MB files
		if err != nil {
			logging.Log(logging.DEBUG, fmt.Sprintf("Cannot parse mutlipart form: %s", err.Error()))
		}
		file, handler, err := r.FormFile("photo_file")
		if err != nil {
			logging.Log(logging.ERROR, fmt.Sprintf("Error retrieving file: %s", err.Error()))
			header.Set("status", fmt.Sprintf("cannot retrieve: %s ||\n", err))
		} else {
			defer file.Close()
			logging.Log(logging.INFO, fmt.Sprintf("Photo %s uploaded for user %s. The file size is %+v. " +
				"MIME header is %+v.", handler.Filename, info.User.Name, handler.Size, handler.Header))
			if err != nil {
				logging.Log(logging.DEBUG, fmt.Sprintf("Error generating temporary file: %s.", err.Error()))
			}
			fileBytes, err := ioutil.ReadAll(file)
			if err != nil {
				logging.Log(logging.ERROR, fmt.Sprintf("Error reading file: %s", err.Error()))
			}
			filename := fmt.Sprintf("%s/user%s.jpeg", paths.FileBasePath, info.User.Name)
			ioutil.WriteFile(filename, encrypt(fileBytes, info.TempUser.Password), 0600)
			if err != nil {
				logging.Log(logging.ERROR, fmt.Sprintf("Cannot write bytes to file: %s", err.Error()))
			}
			info.TempUser.PhotoUrl = filename
			err = DecryptPhoto(info.TempUser.PhotoUrl, info.TempUser.Password, info.TempUser.Name, &info.Photo)
			if err != nil {
				logging.Log(logging.ERROR, fmt.Sprintf(err.Error()))
			}
			setSession(info.User, info.TempUser, info.InfoErr, info.Photo, w)
			logging.Log(logging.INFO, fmt.Sprintln("Successfully uploaded file"))
			header.Set("tempPhoto", info.Photo)
			header.Set("status", "success")
			header.Set("photo", info.TempUser.PhotoUrl)
		}
		http.Redirect(w, r, "/edit", 302)
	} else {
		http.Redirect(w, r, "/", 302)
	}
}

func DiscardHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	info := GetPageInfo(r)
	if info.User.Password != "" {
		if info.TempUser.PhotoUrl != "" && info.TempUser.PhotoUrl != paths.PlaceholderPath {
			err := os.Remove(info.Photo)
			if err != nil {
				logging.Log(logging.DEBUG, fmt.Sprintf("Cannot remove file %s: %s", info.Photo, err.Error()))
			}
			err = DecryptPhoto(info.User.PhotoUrl, info.TempUser.Password, info.User.Name, &info.Photo)
			if err != nil {
				logging.Log(logging.ERROR, fmt.Sprintf(err.Error()))
			}
			logging.Log(logging.INFO, fmt.Sprintf("Temporary file removed."))
		} else {
			logging.Log(logging.DEBUG, fmt.Sprintf("No file to remove."))
		}
		if info.User.PhotoUrl == "" || info.User.PhotoUrl == paths.PlaceholderPath {
			err := os.Remove(info.TempUser.PhotoUrl)
			if err != nil {
				logging.Log(logging.DEBUG, fmt.Sprintf("Cannot remove file %s: %s",
					info.User.PhotoUrl, err.Error()))
			}
			logging.Log(logging.INFO, fmt.Sprintf("Removed temporarily uploaded file."))
		}
		info.TempUser.PhotoUrl = info.User.PhotoUrl
		setSession(info.User, info.TempUser, info.InfoErr, info.Photo, w)
		http.Redirect(w, r, "/view", 302)
	} else {
		logging.Log(logging.INFO, fmt.Sprintf("Access denied. Redirect to homepage."))
		http.Redirect(w, r, "/", 302)
	}
}

// removeFile removes file at src if validator is not empty string or placeholder
func removeFile(src string, validator string) {
	if validator != "" && validator != paths.PlaceholderPath {
		err := os.Remove(src)
		if err != nil {
			logging.Log(logging.DEBUG, fmt.Sprintf("Cannot remove file %s: %s", src, err.Error()))
		}
	}
}

func RemoveHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	info := GetPageInfo(r)
	if info.User.Password != "" {
		removeFile(info.Photo, info.Photo)
		setSession(info.User, info.TempUser, info.InfoErr, paths.PlaceholderPath, w)
		err := ExecuteQuery(db, "UPDATE users SET photo = NULL WHERE username = ?", info.User.Name)
		if err != nil {
			logging.Log(logging.DEBUG, fmt.Sprintf("Cannot remove photo of %s: %s", info.User.Name, err.Error()))
		}
		http.Redirect(w, r, "/edit", 302)
		logging.Log(logging.INFO, fmt.Sprintf("Removed profile photo for user %s", info.User.Name))
	} else {
		logging.Log(logging.INFO, fmt.Sprintf("Access denied. Redirect to homepage."))
		http.Redirect(w, r, "/", 302)
	}
}

func DeleteHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	header.Set("status", "")
	info := GetPageInfo(r)
	if info.User.Password != "" {
		name := info.User.Name
		query := "DELETE FROM users WHERE username = ?"
		err := ExecuteQuery(db, query, name)
		if err != nil {
			logging.Log(logging.ERROR, fmt.Sprintf("Query %s cannot be executed: %s", query, err.Error()))
			header.Set("status", "cannot delete")
		} else {
			logging.Log(logging.INFO, fmt.Sprintf("User %s deleted.", name))
			header.Set("status", fmt.Sprintf("delete %s", name))
		}
		removeFile(info.Photo, info.Photo)
		removeFile(info.User.PhotoUrl, info.Photo)
		clearSession(w)
	}
	http.Redirect(w, r, "/", 302)
}
