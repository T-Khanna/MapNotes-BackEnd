package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/vjeantet/jodaTime"
	"gitlab.doc.ic.ac.uk/g1736215/MapNotes/models"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
	"time"
)

func TestInsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	models.SetDB(db)
	defer db.Close()

	var title string = "Test title"
	var comment string = "test comment"
	var timestamp = jodaTime.Format("YYYY-MM-dd HH:mm", time.Now())
	var longitude float64 = 1.0
	var latitude float64 = 2.0
	//Input string will be converted into a regex
	//Hence I need to double backslash all the special characters
	//One to escape it in a string context, another to escape in a regex context
	mock.ExpectPrepare("INSERT INTO notes\\(title, comments, startTime, endTime, longitude, latitude\\) VALUES\\(\\$1, \\$2, \\$3, \\$4, \\$5, \\$6\\)").
		ExpectExec().
		WithArgs(title, comment, timestamp, timestamp, longitude, latitude).
		WillReturnResult(sqlmock.NewResult(1, 1))
	models.InsertNote(models.Note{Title: title, Comment: comment,
		Start_time: timestamp, End_time: timestamp, Longitude: longitude, Latitude: latitude})
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfufilled expectations: %s", err)
	}
}

func TestGetAllNotes(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	models.SetDB(db)
	defer db.Close()

	var title string = "testing title"
	var comment string = "testing comments"
	var startTime string = "2017-01-01 00:00"
	var endTime string = "2017-05-05 00:00"
	var longitude float64 = 1.0
	var latitude float64 = 2.0
	var id int = 1

	rows := sqlmock.NewRows([]string{"title", "comments", "startTime", "endTime",
		"longitude", "latitude", "id"}).
		AddRow(title, comment, startTime, endTime, longitude, latitude, id).
		AddRow("Harry's world", "Hi Harry", "2017-01-01 00:00", "2017-05-05 00:00", 1.0, 2.0, 1)

	mock.ExpectQuery("SELECT (.+) FROM notes").
		WillReturnRows(rows)
	returnedRows := models.GetAllNotes()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfufilled expectations: %s", err)
	}
	if len(returnedRows) != 2 {
		t.Errorf("Function did not return correct number of rows")
	}

	assert.Equal(t, returnedRows[0].Title, title)
	assert.Equal(t, returnedRows[0].Comment, comment)
	assert.Equal(t, returnedRows[0].Start_time, startTime)
	assert.Equal(t, returnedRows[0].End_time, endTime)
	assert.Equal(t, returnedRows[0].Longitude, longitude)
	assert.Equal(t, returnedRows[0].Latitude, latitude)
	assert.Equal(t, returnedRows[0].Id, id)
}

func TestDelete(t *testing.T){

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		models.SetDB(db)
		defer db.Close()

		var title string = "Test title"

		mock.ExpectPrepare("DELETE FROM Notes WHERE title = \\$1").
		ExpectExec().
		WithArgs(title).
		WillReturnResult(sqlmock.NewResult(1,1))
		models.DeleteNote(title)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfufilled expectations: %s", err)
		}

}
