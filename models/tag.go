package models

import (
	_ "github.com/lib/pq"
	"log"
)

func linkTag(tag string, note_id int64) error {

	// Check first that tag doesn't exist in tag table
	var tag_id int64

	err := db.QueryRow("SELECT id FROM tags WHERE tag = $1", tag).Scan(&tag_id)

	// If s not null, then tag already exists in tag table
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
