package handlers

import (
	"database/sql"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func Initialize(db *sql.DB) {
	ExecuteQuery(db, "DROP TABLE IF EXISTS users")
	/*executeQuery(db, `CREATE TABLE users (
		username    VARCHAR(20) PRIMARY KEY,
		password    VARCHAR(20) NOT NULL,
		photo       VARCHAR(50),
		nickname    VARCHAR(30) COLLATE Latin1_General_100_CI_AI_SC_UTF8
	)`)*/
	ExecuteQuery(db, `CREATE TABLE users (
		username    VARCHAR(20) PRIMARY KEY,
		password    VARCHAR(100) NOT NULL,
		photo       VARCHAR(50),
		nickname    VARCHAR(30) COLLATE utf8mb4_unicode_ci
	)`)
	for i := 0; i < 200; i++ {
		pass, err := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("pass%d%d", i*2, i*2)), 3)
		if err != nil {
			log.Printf("Error: password %s cannot be hashed.", pass)
		}
		ExecuteQuery(db, "INSERT INTO users VALUES(?, ?, ?, ?)",
			fmt.Sprintf("user%d", i),
			pass, nil,
			fmt.Sprintf("nick%d", i),
		)
	}
}

func ExecuteQuery(db *sql.DB, query string, args ...interface{}) {
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
