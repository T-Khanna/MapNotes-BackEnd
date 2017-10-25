package models

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"fmt"
)

// TODO: Change to StartTime and EndTime, and add json tags in camel case.
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
	GetAll          func() ([]Note, error)
	GetActiveAtTime func(string) ([]Note, error)
	Create          func(*Note) (int64, error)
	Delete          func(int64) error
}

// Exported API. Use as models.Notes.Create(..)
// FIXME: To compile, this is named plurarly whilst the actual Note struct must
//        be named singularly. This is inconsistent and an easy to get wrong
//        oddity of the code.
var Notes = NoteOperations{
	GetAll:          getAllNotes,
	GetActiveAtTime: getNotesActiveAtTime,
	Create:          createNote,
	Delete:          deleteNote,
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

func getNotesActiveAtTime(time string) ([]Note, error) {
	rows, err := db.Query("SELECT comments, title, id, startTime, endTime, longitude, latitude FROM notes WHERE (starttime <= $1 AND endtime >= $1) ", time)

	fmt.Println(time)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer rows.Close()
	return convertResultToNotes(rows), nil
}

func getAllNotes() ([]Note, error) {
	rows, err := db.Query("SELECT comments, title, id, startTime, endTime, longitude, latitude FROM notes")
	defer rows.Close()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return convertResultToNotes(rows), nil
}

func convertResultToNotes(rows *sql.Rows) []Note {
	list := make([]Note, 0)
	for rows.Next() {
		var n Note

		err := rows.Scan(&n.Comment, &n.Title,  &n.Id, &n.Start_time, &n.End_time,
			&n.Longitude, &n.Latitude)
		if err != nil {
			log.Println(err)
		} else {
			list = append(list, n)
		}
	}
	return list
}
