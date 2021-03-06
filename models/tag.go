package models

import (
	_ "github.com/lib/pq"
	"log"
)

func linkTag(tag string, note_id int64) error {

	// Check first that tag doesn't exist in tag table
	var tag_id int64

	err := db.QueryRow("SELECT id FROM tags WHERE tag = $1", tag).Scan(&tag_id)

	// If tag_id not null then tag already exists in tag table
	// so return the id of tag
	if err != nil {

		//if you error, is it always because the tag didnt exist?
		tag_id, err = createTag(tag)
		if err != nil {
			return err
		}

	}

	// Prepare sql that inserts note_id and tag_id into notes_tag table

	stmt_notetag, err := db.Prepare("INSERT INTO notestags(note_id, tag_id) VALUES($1, $2)")

	if err != nil {
		return err
	}

	// Execute the INSERT statement
	_, err = stmt_notetag.Exec(note_id, tag_id)

	if err != nil {
		return err
	}

	return nil

}

func createTag(tag string) (int64, error) {

	log.Printf("Inserting tag: %s", tag)
	stmt_tag, err := db.Prepare("INSERT INTO tags(tag) VALUES($1) RETURNING id")

	if err != nil {
		return -1, err
	}

	// Execute the INSERT statement, marshalling the returned id into an int64.
	var id int64
	err = stmt_tag.QueryRow(tag).Scan(&id)

	if err != nil {
		return -1, err
	}

	return id, nil

}

func UpdateTags(note_id int64, tags []string) error {

	for i := 0; i < len(tags); i++ {
		tag := tags[i]
		var tag_id int64

		// Check first that tag doesn't exist in tag table
		err := db.QueryRow("SELECT id FROM tags WHERE tag = $1", tag).Scan(&tag_id)

		// If there is an error, then tag does not exist in tag table
		if err != nil {

			// Create the tag in Tags table and put new tag in tag_id
			tag_id, err = createTag(tag)
			if err != nil {
				return err
			}

		}

		//Next we check the tag is not already linked to that note
		var n int64
		err = db.QueryRow("SELECT note_id FROM NotesTags WHERE note_id = $1 AND tag_id = $2", note_id, tag_id).Scan(&n)
		//If no error, then it's already present in table
		if err == nil {
			continue
		}

		// Prepare sql that inserts note_id and tag_id into notes_tag table
		stmt_notetag, err := db.Prepare("INSERT INTO notestags(note_id, tag_id) VALUES($1, $2)")
		if err != nil {
			return err
		}

		// Execute the INSERT statement
		_, err = stmt_notetag.Exec(note_id, tag_id)
		if err != nil {
			return err
		}
	}
	return nil
}
