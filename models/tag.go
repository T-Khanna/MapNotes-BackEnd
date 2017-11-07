package models

import (
	_ "github.com/lib/pq"
)

type Tag struct {
  Title string `json:"title,omitempty"`
	Id    int    `json:"id,omitempty"`
}

func createTag(tag *Tag) (int64, error) {

	// Check first that tag doesn't exist in tag table
	var s string
	err := db.QueryRow("SELECT title FROM tags WHERE title = ?", tag.Title).Scan(&s)

  // If s not null, then tag already exists in tag table
	// so return the id of tag
	if s != "" {
		return int64(tag.Id), nil
	}

	// if s null, then tag doesn't exist so we can insert tag in tag table
	stmt_tag, err := db.Prepare("INSERT INTO tags(title) VALUES ($1) RETURNING id")

	if err != nil {
		return -1, err
	}

	// Execute the INSERT statement, marshalling the returned id into an int64.
	var id int64
	err = stmt_tag.QueryRow(tag.Title).Scan(&id)

	if err != nil {
		return -1, err
	}

	return id, nil;

}
