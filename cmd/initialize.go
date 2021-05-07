package main

import (
	"database/sql"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func initialize(db *sql.DB) {
	executeQuery(db, "DROP TABLE IF EXISTS users")
	/*executeQuery(db, `CREATE TABLE users (
		username    VARCHAR(20) PRIMARY KEY,
		password    VARCHAR(20) NOT NULL,
		photo       VARCHAR(50),
		nickname    VARCHAR(30) COLLATE Latin1_General_100_CI_AI_SC_UTF8
	)`)*/
	executeQuery(db, `CREATE TABLE users (
		username    VARCHAR(20) PRIMARY KEY,
		password    VARCHAR(100) NOT NULL,
		photo       VARCHAR(50),
		nickname    VARCHAR(30)
	)`)
	/*executeQuery(db, `CREATE TABLE users (
		username	VARCHAR(20) PRIMARY KEY,
		password	VARCHAR(20) NOT NULL
	)`)*/
	for i := 0; i < 30; i++ {
		pass, err := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("pass%d%d", i*2, i*2)), 3)
		if err != nil {
			log.Printf("Error: password %s cannot be hashed.", pass)
		}
		executeQuery(db, "INSERT INTO users VALUES(?, ?, ?, ?)",
			fmt.Sprintf("user%d", i),
			pass,
			fmt.Sprintf("photo%d", i),
			fmt.Sprintf("nick%d", i),
		)
	}
}
