package models

import (
	_ "github.com/lib/pq"
	"log"
)

type Comment struct {
	Id      int64
	User    User
	Comment string
	NoteId  int64
}

type CommentOperations struct {
	GetByNote func(int64) ([]Comment, error)
	Create    func(Comment) error
	//Delete func(string) error
}

var Comments = CommentOperations{
	GetByNote: getCommentsByNoteId,
	Create:    createComment,
}

func getCommentsByNoteId(note_id int64) ([]Comment, error) {
	comments := make([]Comment, 0)
	rows, err := db.Query(`SELECT (comment, comments.id, note_id, users.id, users.name, users.email, users.picture) 
                         FROM comments JOIN users on comments.user_id = users.id WHERE note_id = $1`, note_id)
	if err != nil {
		log.Println(err)
		return comments, err
	}
	defer rows.Close()
	for rows.Next() {
		var comment Comment
		var user User
		err = rows.Scan(&comment.Comment, &comment.Id, &comment.NoteId, &user.Id, &user.Name, &user.Email, &user.Picture)
		if err != nil {
			log.Println(err)
			return comments, err
		}
		comment.User = user
		comments = append(comments, comment)
	}
	return comments, err
}

func createComment(comment Comment) error {
	stmt, err := db.Prepare("INSERT INTO comments(comment, note_id, user_id) VALUES ($1, $2, $3)")

	if err != nil {
		log.Println(err)
		return err
	}

	_, err = stmt.Exec(comment.Comment, comment.NoteId, comment.User.Id)

	if err != nil {
		log.Println(err)
		return err
	}

	return err
}
