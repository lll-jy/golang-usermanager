package handlers

import (
	"database/sql"
	"fmt"
	"git.garena.com/jiayu.li/entry-task/cmd/logging"
	"golang.org/x/crypto/bcrypt"
)

// Initialize setup the initial database data with 200 users with the following configuration:
// username = user<i>, password = pass<i*2><i*2>, photo = nil, nickname = nick<i>.
func Initialize(db *sql.DB) {
	err := ExecuteQuery(db, "DROP TABLE IF EXISTS users")
	if err != nil {
		logging.Log(logging.ERROR, fmt.Sprintf("Cannot drop table: %s", err.Error()))
	}
	err = ExecuteQuery(db, `CREATE TABLE users (
		username    VARCHAR(20) PRIMARY KEY,
		password    VARCHAR(100) NOT NULL,
		photo       VARCHAR(50),
		nickname    VARCHAR(30) COLLATE utf8mb4_unicode_ci
	)`)
	if err != nil {
		logging.Log(logging.ERROR, fmt.Sprintf("Cannot create table: %s", err.Error()))
	}
	for i := 0; i < 200; i++ {
		pass, err := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("pass%d%d", i*2, i*2)), bcrypt.MinCost)
		if err != nil {
			logging.Log(logging.ERROR, fmt.Sprintf("Error: password %s cannot be hashed: %s", pass, err.Error()))
		}
		err = ExecuteQuery(db, "INSERT INTO users VALUES(?, ?, ?, ?)",
			fmt.Sprintf("user%d", i),
			pass, nil,
			fmt.Sprintf("nick%d", i),
		)
		if err != nil {
			logging.Log(logging.ERROR, fmt.Sprintf("Cannot insert user user%d: %s", i, err.Error()))
		}
	}
}

// ExecuteQuery executes the given query (with ? at some places to replace with arguments) to run on db.
func ExecuteQuery(db *sql.DB, query string, args ...interface{}) error {
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(args...)
	stmt.Close()
	if err != nil {
		return err
	}
	return nil
}
