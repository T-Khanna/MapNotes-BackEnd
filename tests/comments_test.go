package tests

import (
	"github.com/stretchr/testify/assert"
	"gitlab.doc.ic.ac.uk/g1736215/MapNotes/models"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	//"log"
	//"math"
	//"reflect"
	"testing"
	//"time"
)

func TestGetCommentsByNoteId(t *testing.T) {

	db, mock := initMockDB(t)
	defer db.Close()

	rows := generateTestRowsComments()
	mock.ExpectQuery("SELECT (.)+ FROM comments JOIN users on (.)+").
		WithArgs(7).
		WillReturnRows(rows)

	comments, _ := models.Comments.GetByNote(7)
	//log.Println(comments)

	assert.True(t, len(comments) == 2, "Did not return correct number of comments")

	assert.Equal(t, "Quantum Beans", comments[0].Comment)
	assert.Equal(t, "Harry", comments[0].User.Name)
	assert.Equal(t, "wwww.harrysworld.com/harry.jpg", comments[0].User.Picture)
	assert.Equal(t, int64(7), comments[0].NoteId)
	assert.Equal(t, int64(7), comments[1].NoteId)
	assert.Equal(t, "wwww.billnyethescienceguy.com/bill.jpg", comments[1].User.Picture)
	assert.Equal(t, int64(1), comments[0].User.Id)

}

func TestInsertComment(t *testing.T) {

	db, mock := initMockDB(t)
	defer db.Close()

	commentUser := generateTestUser("Harry", "Harry@HarrysWorld.com", 1, "blah")
	testComment := models.Comment{Id: 1, Comment: "Innit", NoteId: 5, User: commentUser}

	mock.ExpectPrepare("INSERT INTO comments\\((.)+\\) VALUES \\(\\$1, \\$2\\, \\$3\\)").
		ExpectExec().
		WithArgs(testComment.Comment, testComment.NoteId, testComment.User.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	models.Comments.Create(testComment)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("There were unfufilled expectations: %s", err)
	}

}

func TestDeleteComment(t *testing.T) {
	db, mock := initMockDB(t)
	defer db.Close()

	var comment string = "Quantum Beans"
	mock.ExpectPrepare("DELETE FROM comments WHERE comment = \\$1").
		ExpectExec().
		WithArgs(comment).
		WillReturnResult(sqlmock.NewResult(1, 1))
	models.Comments.Delete(comment)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfufilled expectations: %s", err)
	}

}

//Test Helper Methods
//--------------------------------------------------------------

func generateTestRowsComments() (rows *sqlmock.Rows) {

	noteid := 7
	commentid1 := 1
	user1 := generateTestUser("Harry", "Harry@HarrysWorld.com", 1, "wwww.harrysworld.com/harry.jpg")
	comment1 := "Quantum Beans"

	user2 := generateTestUser("Bill Nye", "thescienceguy@science.com", 2, "wwww.billnyethescienceguy.com/bill.jpg")
	commentid2 := 2
	comment2 := "Science rules!"

	rows = sqlmock.NewRows([]string{"comment", "comments.id", "note_id", "users.id", "users.name", "users.email", "users.picture"}).
		AddRow(comment1, commentid1, noteid, user1.Id, user1.Name, user1.Email, user1.Picture).
		AddRow(comment2, commentid2, noteid, user2.Id, user2.Name, user2.Email, user2.Picture)
	return rows
}
