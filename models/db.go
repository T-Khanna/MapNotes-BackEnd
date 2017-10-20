package models

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var db *sql.DB

type Note struct {
	Title      string
	Comment    string
	Start_time string
	End_time   string
	Longitude  float64
	Latitude   float64
	Id         int
}

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

func InsertNote(note Note) (id int64) {
	stmt, err := db.Prepare("INSERT INTO notes(title, comments, startTime, endTime, longitude, latitude) VALUES($1, $2, $3, $4, $5, $6)")

	if err != nil {
		log.Println(err)
		return -1
	}
	_, err = stmt.Exec(note.Title, note.Comment, note.Start_time, note.End_time,
		note.Longitude, note.Latitude)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT max(id) FROM notes")

	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		err = rows.Scan(&id)
	}

	return
}

func DeleteNote(id int64) {
	stmt, err := db.Prepare("DELETE FROM Notes WHERE id = $1")

	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(title)
	if err != nil {
		log.Fatal(err)
	}
}

func GetTimePeriodNotes(time string) []Note {
	rows, err := db.Query("SELECT * FROM notes WHERE (starttime <= $1 AND endtime >= $1) ", time)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	return ConvertResultToNotes(rows)
}

func GetAllNotes() []Note {
	rows, err := db.Query("SELECT title, comments, startTime, endTime, longitude, latitude, id FROM notes")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	return ConvertResultToNotes(rows)
}

func ConvertResultToNotes(rows *sql.Rows) []Note {

	list := make([]Note, 0)
	for rows.Next() {
		var n Note
		err := rows.Scan(&n.Title, &n.Comment, &n.Start_time, &n.End_time,
			&n.Longitude, &n.Latitude, &n.Id)
		if err != nil {
			log.Fatal(err)
		} else {
			list = append(list, n)
		}
	}
	return list
}
