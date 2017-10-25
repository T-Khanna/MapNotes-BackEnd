package models

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// TODO: Change to StartTime and EndTime, and add json tags in camel case.
type Note struct {
	Title     *string  `json:"title,omitempty"`
	Comment   *string  `json:"comment,omitempty"`
	StartTime *string  `json:"start_time,omitempty"`
	EndTime   *string  `json:"end_time,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Id        *int     `json:"id,omitempty"`
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
		return -1, err
	}

	// Execute the INSERT statement, marshalling the returned id into an int64.
	var id int64
	err = stmt.QueryRow(note.Title, note.Comment, note.StartTime, note.EndTime,
		note.Longitude, note.Latitude).Scan(&id)

	if err != nil {
		return -1, err
	}

	return id, nil
}

//Duplication here with DeleteUser
func deleteNote(id int64) error {
	stmt, prepErr := db.Prepare("DELETE FROM Notes WHERE id = $1")

	if prepErr != nil {
		return prepErr
	}

	_, execErr := stmt.Exec(id)

	if execErr != nil {
		return execErr
	}

	return nil
}

func getNotesActiveAtTime(time string) ([]Note, error) {
	rows, err := db.Query("SELECT comments, title, id, startTime, endTime, longitude, latitude FROM notes WHERE (starttime <= $1 AND endtime >= $1) ", time)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	notes, convErr := convertResultToNotes(rows)

	if convErr != nil {
		return nil, convErr
	}

	return notes, nil
}

func getAllNotes() ([]Note, error) {
	rows, err := db.Query("SELECT comments, title, id, startTime, endTime, longitude, latitude FROM notes")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	notes, convErr := convertResultToNotes(rows)

	if convErr != nil {
		return nil, convErr
	}

	return notes, nil
}

func convertResultToNotes(rows *sql.Rows) ([]Note, error) {
	list := make([]Note, 0)
	for rows.Next() {
		var n Note

		err := rows.Scan(&n.Comment, &n.Title, &n.Id, &n.StartTime, &n.EndTime,
			&n.Longitude, &n.Latitude)
		if err != nil {
			return nil, err
		} else {
			list = append(list, n)
		}
	}
	return list, nil
}
