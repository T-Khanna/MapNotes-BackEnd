package models

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
}

func SetDB(otherDB *sql.DB) {
	db = otherDB
}

func InsertNote(title string, comments string, timestamp string) {
	stmt, err := db.Prepare("INSERT INTO notes(title, comments, time) VALUES($1, $2, $3)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(title, comments, timestamp)
	if err != nil {
		log.Fatal(err)
	}
}
