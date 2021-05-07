package main

// https://gist.github.com/mschoebel/9398202
// https://golang.org/doc/articles/wiki/

import (
	"fmt"
	"net/http"
	"reflect"

	"git.garena.com/jiayu.li/entry-task/cmd/protocol"
	"github.com/gorilla/mux"
	"google.golang.org/protobuf/proto"
)

// server main method

var router = mux.NewRouter()

func main() {

	u := &protocol.User{
		Name:     "name",
		Password: "password",
	}
	out, err := proto.Marshal(u)
	if err != nil {
		fmt.Println("error: ", err)
	} else {
		fmt.Println(out)
		fmt.Println(reflect.TypeOf(out))
		fmt.Println(string(out))
		fmt.Println([]uint8(string(out)))
	}
	u2 := &protocol.User{}
	if err := proto.Unmarshal(out, u2); err != nil {
		fmt.Println("2error: ", err)
	} else {
		fmt.Println(u2)
	}

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

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
