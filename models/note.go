package models

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// TODO: Change to StartTime and EndTime, and add json tags in camel case.
type Note struct {
	Title      *string   `json:"title,omitempty"`
	Comment    *string   `json:"comment,omitempty"`
	StartTime  *string   `json:"start_time,omitempty"`
	EndTime    *string   `json:"end_time,omitempty"`
	Longitude  *float64  `json:"longitude,omitempty"`
	Latitude   *float64  `json:"latitude,omitempty"`
	Id         *int      `json:"id,omitempty"`
	User_email *string   `json:"user_email,omitempty"`
	Tags       *[]string `json:"tags,omitempty"`
}

// Possibly will be a similar struct for any future structs we perform CRUD on.
// Refactor when the time comes.
type NoteOperations struct {
	GetAll          func() ([]Note, error)
	GetActiveAtTime func(string) ([]Note, error)
  GetByUser       func(string) ([]Note, error)
	Create          func(*Note) (int64, error)
	Update          func(*Note) error
	Delete          func(int64) error
}

// Exported API. Use as models.Notes.Create(..)
// FIXME: To compile, this is named plurarly whilst the actual Note struct must
//        be named singularly. This is inconsistent and an easy to get wrong
//        oddity of the code.
var Notes = NoteOperations{
	GetAll:          getAllNotes,
	GetActiveAtTime: getNotesActiveAtTime,
  GetByUser:       getNotesActiveByUser,
	Create:          createNote,
	Update:          updateNote,
	Delete:          deleteNote,
}

func createNote(note *Note) (int64, error) {
	// Prepare sql that inserts the note and returns the new id.
	stmt, err := db.Prepare("INSERT INTO notes(title, comments, startTime, endTime, longitude, latitude, user_email) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id")

	if err != nil {
		return -1, err
	}

	// Execute the INSERT statement, marshalling the returned id into an int64.
	var id int64
	err = stmt.QueryRow(note.Title, note.Comment, note.StartTime, note.EndTime,
		note.Longitude, note.Latitude, note.User_email).Scan(&id)

	if err != nil {
		return -1, err
	}

	// Get tags from note and insert each tag in database
	tags := note.Tags

	for _, t := range *tags {

		err := linkTag(t, id)

		if err != nil {
			return -1, err
		}
	}

	return id, nil
}

func updateNote(note *Note) error {
	/*
	  To implement partial updates:

	  The fields of the Note struct must be pointers, so that they we can
	  distinguish when they've been ommitted from the JSON by checking if the
	  pointer is nil.

	  To dynamically construct the query based on what columns are included, uses a
	  bunch of if statements that check if a column is present and, if so, appends
	  "column name = $n" to the byte buffer.

	  Uses a byte buffer to avoid re-concatenating strings over and over.
	*/

	if note.Id == nil {
		return errors.New("Error: Attempting to update Note but ID not provided")
	}

	// This will be the parameter number of the column-to-update's value in the
	// query that is constructed.. If a column needs to be updated and it's the
	// 'numCols'th column to be added to the query, then it will become parameter
	// '$numCols' in the query.
	numCols := 1

	// Contains the values of the columns to be added. Each time a non-nil field
	// is found in note, that field will be appended to values and numCols
	// incremented. Thus, values[i] will match $i in the query.
	values := []interface{}{}

	// Initialise buffer in which to build the query string.
	var buffer bytes.Buffer
	buffer.WriteString("UPDATE notes SET ")

	if note.Title != nil {
		buffer.WriteString(fmt.Sprintf("title = $%d, ", numCols))
		numCols++
		values = append(values, *note.Title)
	}

	if note.Comment != nil {
		buffer.WriteString(fmt.Sprintf("comments = $%d, ", numCols))
		numCols++
		values = append(values, *note.Comment)
	}

	if note.StartTime != nil {
		buffer.WriteString(fmt.Sprintf("startTime = $%d, ", numCols))
		numCols++
		values = append(values, *note.StartTime)
	}

	if note.EndTime != nil {
		buffer.WriteString(fmt.Sprintf("endTime = $%d, ", numCols))
		numCols++
		values = append(values, *note.EndTime)
	}

	if note.Longitude != nil {
		buffer.WriteString(fmt.Sprintf("Longitude = $%d, ", numCols))
		numCols++
		values = append(values, *note.Longitude)
	}

	if note.Latitude != nil {
		buffer.WriteString(fmt.Sprintf("Latitude = $%d, ", numCols))
		numCols++
		values = append(values, *note.Latitude)
	}

	if note.Tags != nil {
		buffer.WriteString(fmt.Sprintf("Tags = $%d, ", numCols))
		numCols++
		values = append(values, *note.Tags)
	}

	// FIXME: For some reason, bytes.TrimSuffix does not exist, so the trailing
	// comma cannot be removed. Instead, add a superflous 'id = id'.
	buffer.WriteString(fmt.Sprintf("id = %d", *note.Id))

	buffer.WriteString(fmt.Sprintf(" WHERE id = %d;", *note.Id))

	query := buffer.String()

	_, err := db.Exec(query, values...)

	return err
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

func filterNotes(filter string) ([]Note, error) {
	query := fmt.Sprintf(`SELECT comments, title, n.id, startTime, endTime, longitude, latitude, user_email, tag
                        FROM notes as n
                        LEFT JOIN notestags as nt
                        ON n.id = nt.note_id
                        LEFT JOIN tags as t
                        ON t.id = nt.tag_id %s`, filter)

  rows, err := db.Query(query)

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

func getNotesActiveByUser(userEmail string) ([]Note, error) {
  s := fmt.Sprintf("WHERE user_email = '%s'", userEmail)
  log.Println(s)
  return filterNotes(s)
}

func getNotesActiveAtTime(time string) ([]Note, error) {
  s := fmt.Sprintf("WHERE (starttime <= '%[1]s' AND endtime >= '%[1]s')", time)
  log.Println(s)
  return filterNotes(s)
}

func getAllNotes() ([]Note, error) {
  return filterNotes("")
}

func printNote(n Note) {
	log.Println("Printing note...")
	log.Println(*n.Title)
	log.Println(*n.Comment)
	log.Println(*n.StartTime)
	log.Println(*n.EndTime)
	log.Println(*n.Longitude)
	log.Println(*n.Latitude)
	log.Println(*n.Id)
	log.Println(*n.User_email)
}

func convertResultToNotes(rows *sql.Rows) ([]Note, error) {
	log.Println("Entered conversion")
	list := make([]Note, 0)
	var fstNote *Note = nil
	for rows.Next() {
		log.Println("Entered loop")
		var n Note
		var t *string
		err := rows.Scan(&n.Comment, &n.Title, &n.Id, &n.StartTime, &n.EndTime,
			&n.Longitude, &n.Latitude, &n.User_email, &t)
		if err != nil {
			return nil, err
		}
		tagarr := make([]string, 0)
		n.Tags = &tagarr
		if t != nil {
			*n.Tags = append(*n.Tags, *t)
		}
		if fstNote == nil {
			fstNote = &n
		} else if *(*fstNote).Id == *n.Id {
			*fstNote.Tags = append(*fstNote.Tags, *t)
		} else {
			log.Println("appending to list")
			list = append(list, *fstNote)
			fstNote = &n
		}
	}
  if fstNote != nil {
	  list = append(list, *fstNote)
  }
	log.Println(len(list))
	return list, nil
}
