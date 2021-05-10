package main

// https://gist.github.com/mschoebel/9398202
// https://golang.org/doc/articles/wiki/
// https://github.com/go-sql-driver/mysql/wiki/Examples

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"git.garena.com/jiayu.li/entry-task/cmd/handlers"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// server main method

var router = mux.NewRouter()
var db *sql.DB

func setHandleFunc(router *mux.Router) {
	router.HandleFunc("/", makeHandler(handlers.IndexPageHandler))
	router.HandleFunc("/view", makeHandler(handlers.ViewPageHandler))
	router.HandleFunc("/signup", makeHandler(handlers.SignupPageHandler)).Methods("GET")
	router.HandleFunc("/edit", makeHandler(handlers.EditPageHandler)).Methods("GET")
	router.HandleFunc("/reset", makeHandler(handlers.ResetPageHandler)).Methods("GET")

	router.HandleFunc("/login", makeHandler(handlers.LoginHandler)).Methods("POST")
	router.HandleFunc("/logout", makeHandler(handlers.LogoutHandler)).Methods("POST")
	router.HandleFunc("/signup", makeHandler(handlers.SignupHandler)).Methods("POST")
	router.HandleFunc("/edit", makeHandler(handlers.EditHandler)).Methods("POST")
	router.HandleFunc("/reset", makeHandler(handlers.ResetHandler)).Methods("POST")
	router.HandleFunc("/delete", makeHandler(handlers.DeleteHandler)).Methods("POST")
	router.HandleFunc("/upload", makeHandler(handlers.UploadHandler)).Methods("POST")
	router.HandleFunc("/discard", makeHandler(handlers.DiscardHandler)).Methods("POST")
	router.HandleFunc("/remove", makeHandler(handlers.RemoveHandler)).Methods("POST")
}

func tryConnection(db *sql.DB) {
	retryCount := 30
	for {
		err := db.Ping()
		if err != nil {
			if retryCount == 0 {
				log.Fatalf("Not able to establish connection to database")
			}

			log.Printf(fmt.Sprintf("Could not connect to database. Wait 2 seconds. %d retries left...", retryCount))
			retryCount--
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}
}

func setDb() {
	var err error
	db, err = sql.Open("mysql", "root:@/entryTask")
	if err != nil {
		log.Printf("Error connecting to database. %s", err.Error())
	}
}

func makeHandler(fn func(*sql.DB, http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(db, w, r)
	}
}

func main() {
	setDb()
	defer db.Close()

	setHandleFunc(router)

	// https://www.sohamkamani.com/golang/how-to-build-a-web-application/
	staticFileDir := http.Dir("./")
	staticFileHandler := http.StripPrefix("/", http.FileServer(staticFileDir))
	router.PathPrefix("/").Handler(staticFileHandler).Methods("GET")

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
