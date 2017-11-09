package tests

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/vjeantet/jodaTime"
	"gitlab.doc.ic.ac.uk/g1736215/MapNotes/models"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	//"log"
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
	mock.ExpectQuery("SELECT comments, title, n.id, startTime, endTime, longitude, latitude, user_email, tag FROM notes").
		WillReturnRows(rows)

	//May need to check the err returned in the line below
	returnedRows, err := models.Notes.GetAll()
	if err != nil {
		t.Errorf(err.Error())
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfufilled expectations: %s", err)
	}
	if len(returnedRows) != 2 {
		t.Errorf("Function did not return correct number of rows. Returned %d rows", len(returnedRows))
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

	mock.ExpectQuery(`SELECT comments, title, n.id, startTime, endTime, longitude, latitude, user_email, tag
										FROM notes as n
										LEFT JOIN notestags as nt
										ON n.id = nt.note_id
										LEFT JOIN tags as t
										ON t.id = nt.tag_id
										WHERE \(starttime <= (.)+ AND endtime >= (.)+\)`).
		WillReturnRows(rows)
	returnedRows, _ := models.Notes.GetActiveAtTime("2017-01-01 00:00")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfufilled expectations: %s", err)
	}
	if len(returnedRows) != 2 {
		t.Errorf("Function did not return correct number of rows. Returned %d rows", len(returnedRows))
		return
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

	var email string = "user@email.com"

	mock.ExpectPrepare("INSERT INTO users\\(email\\) VALUES\\(\\$1\\)").
		ExpectExec().
		WithArgs(email).
		WillReturnResult(sqlmock.NewResult(1, 1))

	models.Users.Create(&models.User{Email: email})
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
	tags := []string{"Harry", "Beans"}

	note = models.Note{
		Title:      &title,
		Comment:    &comment,
		StartTime:  &startTime,
		EndTime:    &endTime,
		Longitude:  &longitude,
		Latitude:   &latitude,
		Id:         &id,
		User_email: &email,
		Tags:       &tags,
	}

	rows = sqlmock.NewRows([]string{"comments", "title", "n.id", "startTime", "endTime",
		"longitude", "latitude", "user_email", "tag"}).
		AddRow(comment, title, id, startTime, endTime, longitude, latitude, email, "Harry").
		AddRow("Harry's world", "Hi Harry", 2, "2017-01-01 00:00", "2017-05-05 00:00", 1.0, 2.0, "hello@mail.com", "Beans")
	return
}
