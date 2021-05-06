// Package main is the entry-point for the go-sockets client sub-project.
// The go-sockets project is available under the GPL-3.0 License in LICENSE.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

// Application constants, defining host, port, and protocol.
const (
	connHost = "localhost"
	connPort = "8080"
	connType = "tcp"
)

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello world</h1><p>content</p>")
	log.Println("Open home page")
}

func main() {
	// Start the client and connect to the server.
	fmt.Println("Connecting to", connType, "server", connHost+":"+connPort)
	conn, err := net.Dial(connType, connHost+":"+connPort)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}

	/*
		for {
			http.HandleFunc("/", homePageHandler)
			conn.Write([]byte("connected" + connPort))
			msg, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				fmt.Println("Err: ", err.Error())
			}
			fmt.Println(msg)
			log.Fatal(http.ListenAndServe(":"+connPort, nil))
		}
	*/

	// Create new reader from Stdin.
	reader := bufio.NewReader(os.Stdin)
	// go http.HandleFunc("/", homePageHandler)
	// go log.Fatal(http.ListenAndServe(":"+connPort, nil))

	// run loop forever, until exit.
	for {
		// Prompting message.
		fmt.Print("Text to send: ")

		// Read in input until newline, Enter key.
		input, _ := reader.ReadString('\n')

		// Send to socket connection.
		conn.Write([]byte(input))

		// Listen for relay.
		message, _ := bufio.NewReader(conn).ReadString('\n')

		// Print server relay.
		log.Print("Server relay: " + message)

		// added
		http.HandleFunc("/", homePageHandler)
		// go log.Fatal(http.ListenAndServe(":"+connPort, nil))
	}
}

// Reference:
// socket programming: https://github.com/Alice-Williams-Tech/go-sockets/tree/v0.1.0
