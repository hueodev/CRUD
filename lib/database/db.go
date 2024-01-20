package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func DB() {
	// Create database if it doesn't exist
	conn, err := sql.Open("sqlite3", "./lib/database/db.sqlite")
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = conn.Exec("CREATE TABLE IF NOT EXISTS messages (id INTEGER, created TEXT, username TEXT, msg TEXT)")
	if err != nil {
		log.Fatal(err.Error())
	}
}
