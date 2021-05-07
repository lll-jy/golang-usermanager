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

func dummySelect(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		panic(err.Error())
	}

	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	values := make([]sql.RawBytes, len(columns))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}

		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			log.Println(columns[i], ": ", value)
		}
		log.Println("--------")
	}
	if err = rows.Err(); err != nil {
		panic(err.Error())
	}
}

func setDb() {
	var err error
	db, err = sql.Open("mysql", "root:@/mysql")
	if err != nil {
		log.Printf("Error connecting to database. %s", err.Error())
	}
}

func main() {
	// setDb()
	// defer db.Close()
	// dummySelect(db)
	setHandleFunc(router)

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
