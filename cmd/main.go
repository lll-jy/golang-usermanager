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

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// server main method

var router = mux.NewRouter()
var db *sql.DB

func setHandleFunc(router *mux.Router) {
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
	router.HandleFunc("/delete", deleteHandler).Methods("POST")
	router.HandleFunc("/upload", uploadHandler).Methods("POST")
	router.HandleFunc("/discard", discardHandler).Methods("POST")
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

func executeQuery(db *sql.DB, query string, args ...interface{}) {
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Printf("Query %s cannot be executed due to error: %s", query, err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec(args...)
	if err != nil {
		log.Printf("Query %s cannot be executed due to error: %s", query, err.Error())
	}
}

func main() {
	setDb()
	defer db.Close()
	// initialize(db)

	setHandleFunc(router)

	// https://www.sohamkamani.com/golang/how-to-build-a-web-application/
	staticFileDir := http.Dir("./")
	staticFileHandler := http.StripPrefix("/", http.FileServer(staticFileDir))
	router.PathPrefix("/").Handler(staticFileHandler).Methods("GET")

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
