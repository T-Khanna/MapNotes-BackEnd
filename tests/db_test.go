package tests

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/vjeantet/jodaTime"
	"gitlab.doc.ic.ac.uk/g1736215/MapNotes/models"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
	"time"
)

func TestInsertNote(t *testing.T) {
	db, mock := initMockDB(t)
	defer db.Close()

	var title string = "Test title"
	var comment string = "test comment"
	var timestamp = jodaTime.Format("YYYY-MM-dd HH:mm", time.Now())
	var longitude float64 = 1.0
	var latitude float64 = 2.0
	var id int = -1
	var email string = "test@mapnotes.co.uk"

	//Input string will be converted into a regex
	//Hence I need to double backslash all the special characters
	//One to escape it in a string context, another to escape in a regex context

	mock.ExpectPrepare("INSERT INTO notes\\(title, comments, startTime, endTime, longitude, latitude, user_email\\) VALUES\\(\\$1, \\$2, \\$3, \\$4, \\$5, \\$6, \\$7\\)").
		ExpectQuery().
		WithArgs(title, comment, timestamp, timestamp, longitude, latitude, email)


	models.Notes.Create(&models.Note{Title: &title, Comment: &comment,
	StartTime: &timestamp, EndTime: &timestamp, Longitude: &longitude, Latitude: &latitude, Id: &id, User_email: &email})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfufilled expectations: %s", err)
	}
}

func TestGetAllNotes(t *testing.T) {
	db, mock := initMockDB(t)
	defer db.Close()

	rows, note := generateTestRows()
	mock.ExpectQuery("SELECT (.+) FROM notes").
		WillReturnRows(rows)

		//May need to check the err returned in the line below
	returnedRows, _ := models.Notes.GetAll()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfufilled expectations: %s", err)
	}
	if len(returnedRows) != 2 {
		t.Errorf("Function did not return correct number of rows")
	}

	assert.Equal(t, returnedRows[0].Title, note.Title)
	assert.Equal(t, returnedRows[0].Comment, note.Comment)
	assert.Equal(t, returnedRows[0].StartTime, note.StartTime)
	assert.Equal(t, returnedRows[0].EndTime, note.EndTime)
	assert.Equal(t, returnedRows[0].Longitude, note.Longitude)
	assert.Equal(t, returnedRows[0].Latitude, note.Latitude)
	assert.Equal(t, returnedRows[0].Id, note.Id)
	assert.Equal(t, returnedRows[0].User_email, note.User_email)
}

func TestDeleteNote(t *testing.T) {
	testDelete("Notes", t, models.Notes.Delete)
}

func TestGetTimePeriodNotes(t *testing.T) {

	db, mock := initMockDB(t)
	defer db.Close()

	rows, note := generateTestRows()

	mock.ExpectQuery("SELECT comments, title, id, startTime, endTime, longitude, latitude, user_email FROM notes WHERE \\(starttime <= \\$1 AND endtime >= \\$1\\)").
		WithArgs("2017-01-01 00:00").
		WillReturnRows(rows)
	returnedRows, _:= models.Notes.GetActiveAtTime("2017-01-01 00:00")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfufilled expectations: %s", err)
	}
	if len(returnedRows) != 2 {
		t.Errorf("Function did not return correct number of rows")
	}

	assert.Equal(t, returnedRows[0].Title, note.Title)
	assert.Equal(t, returnedRows[0].Comment, note.Comment)
	assert.Equal(t, returnedRows[0].StartTime, note.StartTime)
	assert.Equal(t, returnedRows[0].EndTime, note.EndTime)
	assert.Equal(t, returnedRows[0].Longitude, note.Longitude)
	assert.Equal(t, returnedRows[0].Latitude, note.Latitude)
	assert.Equal(t, returnedRows[0].Id, note.Id)
	assert.Equal(t, returnedRows[0].User_email, note.User_email)


}

func TestInsertUser(t *testing.T) {
	db, mock := initMockDB(t)
	defer db.Close()

	var username string = "Harry"
	var email string = "user@email.com"

	mock.ExpectPrepare("INSERT INTO users\\(email, username\\) VALUES\\(\\$1, \\$2\\)").
		ExpectExec().
		WithArgs(email, username).
		WillReturnResult(sqlmock.NewResult(1, 1))

	models.Users.Create(&models.User{Username: username, Email: email})
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfufilled expectations: %s", err)
	}
}

type DeleteFunc func(int64) error
/*
func TestDeleteUser(t *testing.T) {
	testDelete("users", t, models.Notes.Delete)
}

*/

func testDelete(tableName string, t *testing.T, deleter DeleteFunc) {
	db, mock := initMockDB(t)
	defer db.Close()

	var id int64 = 2
	mock.ExpectPrepare("DELETE FROM " + tableName + " WHERE id = \\$1").
		ExpectExec().
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	deleter(id)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfufilled expectations: %s", err)
	}

}

func initMockDB(t *testing.T) (db *sql.DB, mock sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	models.SetDB(db)
	return
}

func generateTestRows() (rows *sqlmock.Rows, note models.Note) {
	title := "testing title"
	comment := "testing comments"
	startTime := "2017-01-01 00:00"
	endTime := "2017-05-05 00:00"
	longitude := 1.0
	latitude := 2.0
	id := 1
  email := "test@mapnotes.co.uk"

	note = models.Note{
		Title:      &title,
		Comment:    &comment,
		StartTime:  &startTime,
		EndTime:    &endTime,
		Longitude:  &longitude,
		Latitude:   &latitude,
		Id:         &id,
    User_email: &email,
	}

	rows = sqlmock.NewRows([]string{"comments", "title", "id", "startTime", "endTime",
		"longitude", "latitude", "user_email"}).
		AddRow(comment, title, id, startTime, endTime, longitude, latitude, email).
		AddRow("Harry's world", "Hi Harry", 1, "2017-01-01 00:00", "2017-05-05 00:00", 1.0, 2.0, "hello@mail.com")
	return
}
