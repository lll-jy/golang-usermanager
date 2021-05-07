package main

// https://gist.github.com/mschoebel/9398202
// https://golang.org/doc/articles/wiki/

import (
	"net/http"

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

func main() {

	setHandleFunc(router)

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
