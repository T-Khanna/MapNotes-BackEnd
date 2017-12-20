package models

import (
	_ "github.com/lib/pq"
	"log"
)

type Attended struct {
	User   User
	NoteId int64
}

type AttendedOperations struct {
	Create func(Attended) error
	Delete func(int64) error
}

var Attend = AttendedOperations{
	Create: createAttended,
	Delete: deleteAllAttendanceByNoteId,
}

func createAttended(attended Attended) error {
	user_id := attended.User.Id
	if user_id == -1 {
		userErr, id := GetUserId(attended.User)
		if userErr != nil {
			return userErr
		}
		user_id = id
	}
	stmt, err := db.Prepare("INSERT INTO attended(note_id, user_id) VALUES ($1, $2)")

	if err != nil {
		log.Println(err)
		return err
	}

	_, err = stmt.Exec(attended.NoteId, user_id)

	if err != nil {
		log.Println(err)
		return err
	}

	return err
}

// Delete all attendance for a specific note_id
func deleteAllAttendanceByNoteId(id int64) error {
	stmt, prepErr := db.Prepare("DELETE FROM attended WHERE note_id = $1")

	if prepErr != nil {
		return prepErr
	}

	_, execErr := stmt.Exec(id)

	if execErr != nil {
		return execErr
	}

	return nil
}
