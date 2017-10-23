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
