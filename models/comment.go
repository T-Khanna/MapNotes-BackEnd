package models

import (
	_ "github.com/lib/pq"
	"log"
	"fmt"
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
	Delete    func(string) error
	Merge     func([]int64, int64) error
}

var Comments = CommentOperations{
	GetByNote: getCommentsByNoteId,
	Create:    createComment,
	Delete:    deleteComment,
	Merge:     mergeComments,
}

func mergeComments(oldnoteids []int64, newnoteid int64) (err error) {

	var idString string = ConvertIntArrayToString(oldnoteids)
	q1 :=  fmt.Sprintf("UPDATE comments SET note_id = %d WHERE note_id in %s", newnoteid ,idString)
	_, err = db.Exec(q1)

	return

}

func getCommentsByNoteId(note_id int64) ([]Comment, error) {
	comments := make([]Comment, 0)

	log.Println("Attempting to retrieve comments with note id ", note_id)
	rows, err := db.Query(`SELECT comment, comments.id, note_id, users.id, users.name, users.email, users.picture
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
	user_id := comment.User.Id
	if user_id == -1 {
		userErr, id := GetUserId(comment.User)
		if userErr != nil {
			return userErr
		}
		user_id = id
	}
	stmt, err := db.Prepare("INSERT INTO comments(comment, note_id, user_id) VALUES ($1, $2, $3)")

	if err != nil {
		log.Println(err)
		return err
	}

	_, err = stmt.Exec(comment.Comment, comment.NoteId, user_id)

	if err != nil {
		log.Println(err)
		return err
	}

	return err
}

// Delete comment with a specific comment string
func deleteComment(comment string) error {
	stmt, prepErr := db.Prepare("DELETE FROM comments WHERE comment = $1")

	if prepErr != nil {
		return prepErr
	}

	_, execErr := stmt.Exec(comment)

	if execErr != nil {
		return execErr
	}

	return nil
}

// Delete all comments for a specific note_id
func deleteAllCommentinNote(id int64) error {
	stmt, prepErr := db.Prepare("DELETE FROM comments WHERE note_id = $1")

	if prepErr != nil {
		return prepErr
	}

	_, execErr := stmt.Exec(id)

	if execErr != nil {
		return execErr
	}

	return nil
}
