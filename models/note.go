package models

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"log"
  "sort"

	_ "github.com/lib/pq"
)

// TODO: Change to StartTime and EndTime, and add json tags in camel case.
type Note struct {
	Title     *string   `json:"title,omitempty"`
	Comment   *string   `json:"comment,omitempty"`
	StartTime *string   `json:"start_time,omitempty"`
	EndTime   *string   `json:"end_time,omitempty"`
	Longitude *float64  `json:"longitude,omitempty"`
	Latitude  *float64  `json:"latitude,omitempty"`
	Id        *int      `json:"id,omitempty"`
	Users     *[]User   `json:"users,omitempty"`
	Tags      *[]string `json:"tags,omitempty"`
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

func mergeNotes(oldIds []int64, newNote Note) {

	deleteNotes(oldIds)
	createNote(&newNote)

	//so we get an array of ids of notes to delete, which will use cascades
	//the cascades will handle all of the tags and users deletion
	//the agg function will have already figured out who the users, tags and other attributes are for this new note
	//then we just have to create a new note.

}

func deleteNotes(deleteids []int64) (err error) {

	//delete all of the notes using deleteids

	for i := 0; i < len(deleteids); i++ {

		err = deleteNote(deleteids[i])

		if err != nil {

			return

		}

	}

	return

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

	// Get tags from note and insert each tag in database
	tags := note.Tags

	for _, t := range *tags {

		err := linkTag(t, id)

		if err != nil {
			return -1, err
		}
	}

	users := note.Users

	for _, u := range *users {

		_, uid := GetUserId(u)

		err := NotesUsers.Insert(id, uid)

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

	//also need code to delete notestags entries and notesuser entries
	//or set up cascade deletes in database?

	if prepErr != nil {
		return prepErr
	}

	_, execErr := stmt.Exec(id)

	if execErr != nil {
		return execErr
	}

	return nil
}

func filterNotes(whereClause string) ([]Note, error) {
	notesWithUsersQuery := fmt.Sprintf(
		`SELECT comments, title, n.id, startTime, endTime, longitude, latitude, u.id, u.name, u.email
    FROM notes as n
    JOIN notesusers as nu ON n.id = nu.note_id
    JOIN users as u ON nu.user_id = u.id
    %s`, whereClause,
	)

	notesWithTagsQuery := fmt.Sprintf(
		`SELECT n.id, t.tag
    FROM notes as n
    LEFT JOIN notestags as nt ON n.id = nt.note_id
    LEFT JOIN tags as t ON nt.tag_id = t.id
    %s`, whereClause,
	)

	// Get rows of (...note, ...user)
	notesWithUsersRows, uErr := db.Query(notesWithUsersQuery)
	if uErr != nil {
		return nil, uErr
	}
	defer notesWithUsersRows.Close()

	// Get rows of (note.id, tag)
	notesWithTagsRows, tErr := db.Query(notesWithTagsQuery)
	if tErr != nil {
		return nil, tErr
	}
	defer notesWithTagsRows.Close()

	return rowsToNotes(notesWithUsersRows, notesWithTagsRows)
}

type reverseChronologicalOrder []Note

func (a reverseChronologicalOrder) Len() int {
  return len(a)
}

func (a reverseChronologicalOrder) Swap(i, j int) {
  a[i], a[j] = a[j], a[i]
}

func (a reverseChronologicalOrder) Less(i, j int) bool {
    if *a[i].StartTime > *a[j].StartTime {
       return true
    }
    if *a[i].StartTime < *a[j].StartTime {
       return false
    }
    return *a[i].EndTime > *a[j].EndTime
}

/**
 * Takes rows of (...note, ...user) and (note.id, tag) and constructs a slice
 * of note objects with the tag and user arrays filled in.
 */
func rowsToNotes(notesWithUsersRows *sql.Rows, notesWithTagsRows *sql.Rows) ([]Note, error) {
	/*
	   Loop over notes with users, populating each note's Users field. As this is
	   done, insert the notes into a hash map by their id. When iterating over
	   notes with tags, insert the tags into the note struct taken from the map
	   with the correct id.
	*/
	notesById := make(map[int]Note)

	// Iterate over the notesWithUsersRows, populating notesById and each note's
	// users field.
  var emptyNote Note
	var note Note
	var user User

	for notesWithUsersRows.Next() {
		err := notesWithUsersRows.Scan(
			&note.Comment,
			&note.Title,
			&note.Id,
			&note.StartTime,
			&note.EndTime,
			&note.Longitude,
			&note.Latitude,
			&user.Id,
			&user.Name,
			&user.Email,
		)
		if err != nil {
			return nil, err
		}

		// If not already hit this note, add it to the map and initialise its users.
		// Else, get the note from the map and add this user to its users.
		notesWithUsers := notesById[*note.Id]
		if notesWithUsers == emptyNote {
			note.Users = &[]User{user}
			notesById[*note.Id] = note
		} else {
			noteUsers := notesWithUsers.Users
			*noteUsers = append(*noteUsers, user)
		}
	}

	// Iterate over notes with tags rows, adding tags to the notes from the map,
	// whose Users field has been created.
	var tag *string
	for notesWithTagsRows.Next() {
		err := notesWithTagsRows.Scan(&note.Id, &tag)
		if err != nil {
			return nil, err
		}

		// If tags not yet constructed, construct it, else append.
		noteWithUsers := notesById[*note.Id]
		if noteWithUsers == emptyNote {
			// FIXME: wtf why is this happening?
			log.Printf("Error in models.rowsToResults(): Found note with tag but not user: %+v", note)
		} else if tag != nil && noteWithUsers.Tags == nil {
			noteWithUsers.Tags = &[]string{*tag}
		} else if tag != nil {
			*noteWithUsers.Tags = append(*noteWithUsers.Tags, *tag)
		}
	}

	// Convert map to slice
	// FIXME: Is there a way to do this without required this extra iteration
	//        over all notes?
	var notes []Note = make([]Note, 0)
	for _, note := range notesById {
		notes = append(notes, note)
	}

  // Sorting notes
  sort.Sort(reverseChronologicalOrder(notes))

	return notes, nil
}

func getNotesActiveByUser(userEmail string) ([]Note, error) {
	s := fmt.Sprintf("WHERE email = '%s'", userEmail)
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
