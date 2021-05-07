package main

// https://gist.github.com/mschoebel/9398202
// https://golang.org/doc/articles/wiki/

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// server main method

var router = mux.NewRouter()

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

func setDb() {
	db, err := sql.Open("mysql", "root:@/mysql")
	if err != nil {
		log.Printf("Error connecting to database. %s", err.Error())
	}
	defer db.Close()

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

func main() {
	setDb()
	setHandleFunc(router)

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
