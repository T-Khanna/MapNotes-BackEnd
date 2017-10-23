package models

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Note struct {
	Title      string
	Comment    string
	Start_time string
	End_time   string
	Longitude  float64
	Latitude   float64
	Id         int
}

// Possibly will be a similar struct for any future structs we perform CRUD on.
// Refactor when the time comes.
type NoteOperations struct {
	Create func(note *Note) (int64, error)
	Delete func(id int64) error
}

// Exported API. Use as models.Notes.Create(..)
// FIXME: To compile, this is named plurarly whilst the actual Note struct must
//        be named singularly. This is inconsistent and an easy to get wrong
//        oddity of the code.
var Notes = NoteOperations{
	Create: createNote,
	Delete: deleteNote,
}

func createNote(note *Note) (int64, error) {
	// Prepare sql that inserts the note and returns the new id.
	stmt, err := db.Prepare("INSERT INTO notes(title, comments, startTime, endTime, longitude, latitude) VALUES($1, $2, $3, $4, $5, $6) RETURNING id")

	if err != nil {
		log.Println(err)
		return -1, err
	}

	// Execute the INSERT statement, marshalling the returned id into an int64.
	var id int64
	err = stmt.QueryRow(note.Title, note.Comment, note.Start_time, note.End_time,
		note.Longitude, note.Latitude).Scan(&id)

	if err != nil {
		log.Println(err)
		return -1, err
	}

	return id, nil
}

//Duplication here with DeleteUser
func deleteNote(id int64) error {
	stmt, prepErr := db.Prepare("DELETE FROM Notes WHERE id = $1")

	if prepErr != nil {
		log.Println(prepErr)
		return prepErr
	}

	_, execErr := stmt.Exec(id)

	if execErr != nil {
		log.Println(execErr)
		return execErr
	}

	return nil
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
