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

func InsertNote(title string, comments string, timestamp string) {
	rows, err := db.Query("SELECT id FROM Notes ORDER BY id DESC LIMIT 1")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var id int = 0
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
	}
	stmt, err := db.Prepare("INSERT INTO notes(title, comments, time, id) VALUES($1, $2, $3, $4)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(title, comments, timestamp, id+1)
	if err != nil {
		log.Fatal(err)
	}
}
