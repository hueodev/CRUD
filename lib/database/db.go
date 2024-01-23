package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func DB() {
	conn, err := sql.Open("sqlite3", "./lib/database/db.sqlite")
	if err != nil {
		log.Fatal(err.Error())
	}

	// Create database if it doesn't exist
	_, err = conn.Exec("CREATE TABLE IF NOT EXISTS messages (id INTEGER, created TEXT, username TEXT, msg TEXT)")
	if err != nil {
		log.Fatal(err.Error())
	}
}
