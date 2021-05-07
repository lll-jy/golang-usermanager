package main_test

import (
	"database/sql"
	"fmt"
	"log"
)

func executeQuery(db *sql.DB, query string, args ...interface{}) {
	stmt, err := db.Prepare(query)
	if err != nil {
		panic(fmt.Sprintf("query %s err: %s", query, err.Error()))
	}
	defer stmt.Close()
	_, err = stmt.Exec(args...)
	if err != nil {
		panic(err.Error())
	}
	log.Println("sucessful!!")
}

func initialize(db *sql.DB) {
	executeQuery(db, "DROP TABLE IF EXISTS users")
	/*executeQuery(db, `CREATE TABLE users (
		username    VARCHAR(20) PRIMARY KEY,
		password    VARCHAR(20) NOT NULL,
		photo       VARCHAR(50),
		nickanme    VARCHAR(30) COLLATE Latin1_General_100_CI_AI_SC_UTF8
	)`)*/
	executeQuery(db, `CREATE TABLE users (
		username	VARCHAR(20) PRIMARY KEY,
		password	VARCHAR(20) NOT NULL
	)`)
	for i := 0; i < 30; i++ {
		executeQuery(db, "INSERT INTO users VALUES(?, ?)", fmt.Sprintf("user%d", i), fmt.Sprintf("pass%d%d", i*2, i*2))
	}
}
