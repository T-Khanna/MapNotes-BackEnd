package models

type MappingOperations struct {
	Insert func(int64, int64) error
}

var NotesUsers = MappingOperations{
	Insert: insertNoteIdAndUserIDMapping,
}

func insertNoteIdAndUserIDMapping(note_id int64, user_id int64) error {
	stmt, prepErr := db.Prepare(`INSERT INTO notesusers(note_id, user_id)
													     VALUES($1, $2)`)
	if prepErr != nil {
		return prepErr
	}

	_, execErr := stmt.Exec(note_id, user_id)

	if execErr != nil {
		return execErr
	}

	return nil
}
