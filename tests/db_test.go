package main

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
	//Input string will be converted into a regex
	//Hence I need to double backslash all the special characters
	//One to escape it in a string context, another to escape in a regex context

	mock.ExpectPrepare("INSERT INTO notes\\(title, comments, startTime, endTime, longitude, latitude\\) VALUES\\(\\$1, \\$2, \\$3, \\$4, \\$5, \\$6\\)").
		ExpectExec().
		WithArgs(title, comment, timestamp, timestamp, longitude, latitude).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rows := sqlmock.NewRows([]string{"max"}).AddRow(1)
	mock.ExpectQuery("SELECT max\\(id\\) FROM notes").WillReturnRows(rows)
	models.InsertNote(models.Note{Title: title, Comment: comment,
		Start_time: timestamp, End_time: timestamp, Longitude: longitude, Latitude: latitude, Id: id})
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
	returnedRows := models.GetAllNotes()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfufilled expectations: %s", err)
	}
	if len(returnedRows) != 2 {
		t.Errorf("Function did not return correct number of rows")
	}

	assert.Equal(t, returnedRows[0].Title, note.Title)
	assert.Equal(t, returnedRows[0].Comment, note.Comment)
	assert.Equal(t, returnedRows[0].Start_time, note.Start_time)
	assert.Equal(t, returnedRows[0].End_time, note.End_time)
	assert.Equal(t, returnedRows[0].Longitude, note.Longitude)
	assert.Equal(t, returnedRows[0].Latitude, note.Latitude)
	assert.Equal(t, returnedRows[0].Id, note.Id)
}

func TestDeleteNote(t *testing.T) {
	testDelete("Notes", t, models.DeleteNote)
}

func TestGetTimePeriodNotes(t *testing.T) {

	db, mock := initMockDB(t)
	defer db.Close()

	rows, note := generateTestRows()

	mock.ExpectQuery("SELECT \\* FROM notes WHERE \\(starttime <= \\$1 AND endtime >= \\$1\\)").
		WithArgs("2017-01-01 00:00").
		WillReturnRows(rows)
	returnedRows := models.GetTimePeriodNotes("2017-01-01 00:00")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfufilled expectations: %s", err)
	}
	if len(returnedRows) != 2 {
		t.Errorf("Function did not return correct number of rows")
	}

	assert.Equal(t, returnedRows[0].Title, note.Title)
	assert.Equal(t, returnedRows[0].Comment, note.Comment)
	assert.Equal(t, returnedRows[0].Start_time, note.Start_time)
	assert.Equal(t, returnedRows[0].End_time, note.End_time)
	assert.Equal(t, returnedRows[0].Longitude, note.Longitude)
	assert.Equal(t, returnedRows[0].Latitude, note.Latitude)
	assert.Equal(t, returnedRows[0].Id, note.Id)

}

func TestInsertUser(t *testing.T) {
	db, mock := initMockDB(t)
	defer db.Close()

	var username string = "Harry"
	var password string = "1234"

	mock.ExpectPrepare("INSERT INTO users\\(username, password\\) VALUES\\(\\$1, \\$2\\)").
		ExpectExec().
		WithArgs(username, password).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rows := sqlmock.NewRows([]string{"max"}).AddRow(1)
	mock.ExpectQuery("SELECT max\\(id\\) FROM users").WillReturnRows(rows)

	models.InsertUser(models.User{Userid: -1, Username: username, Password: password})
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfufilled expectations: %s", err)
	}
}

type DeleteFunc func(int64)

func TestDeleteUser(t *testing.T) {
	testDelete("users", t, models.DeleteUser)
}

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
	var title string = "testing title"
	var comment string = "testing comments"
	var startTime string = "2017-01-01 00:00"
	var endTime string = "2017-05-05 00:00"
	var longitude float64 = 1.0
	var latitude float64 = 2.0
	var id int = 1

	note = models.Note{title, comment, startTime, endTime, longitude, latitude, id}

	rows = sqlmock.NewRows([]string{"title", "comments", "startTime", "endTime",
		"longitude", "latitude", "id"}).
		AddRow(title, comment, startTime, endTime, longitude, latitude, id).
		AddRow("Harry's world", "Hi Harry", "2017-01-01 00:00", "2017-05-05 00:00", 1.0, 2.0, 1)
	return
}
