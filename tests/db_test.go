package tests

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/vjeantet/jodaTime"
	"gitlab.doc.ic.ac.uk/g1736215/MapNotes/models"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	//"log"
	"math"
	"reflect"
	"testing"
	"time"
)

const TOLERANCE = 0.0000001

func TestInsertNote(t *testing.T) {
	db, mock := initMockDB(t)
	defer db.Close()

	var title string = "Test title"
	var comment string = "test comment"
	var timestamp = jodaTime.Format("YYYY-MM-dd HH:mm", time.Now())
	var longitude float64 = 1.0
	var latitude float64 = 2.0
	var id int64 = -1
	//var email string = "test@mapnotes.co.uk"

	//Input string will be converted into a regex
	//Hence I need to double backslash all the special characters
	//One to escape it in a string context, another to escape in a regex context

	mock.ExpectPrepare("INSERT INTO notes\\((.)+\\) VALUES\\((.)+\\)").
		ExpectQuery().
		WithArgs(title, comment, timestamp, timestamp, longitude, latitude)

	models.Notes.Create(&models.Note{Title: &title, Comment: &comment,
		StartTime: &timestamp, EndTime: &timestamp, Longitude: &longitude, Latitude: &latitude, Id: &id})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfufilled expectations: %s", err)
	}
}

func TestGetAllNotes(t *testing.T) {
	db, mock := initMockDB(t)
	defer db.Close()

	rows2 := generateTestRowsNoteUser()
	rows3 := generateTestRowsNoteTag()

	mock.ExpectQuery("SELECT (.)+ FROM notes (.)+ JOIN").
		WillReturnRows(rows2)

	mock.ExpectQuery("SELECT (.)+ FROM notes (.)+ LEFT JOIN").
		WillReturnRows(rows3)

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

	/*
		assert.Equal(t, returnedRows[0].Title, note.Title)
		assert.Equal(t, returnedRows[0].Comment, note.Comment)
		assert.Equal(t, returnedRows[0].StartTime, note.StartTime)
		assert.Equal(t, returnedRows[0].EndTime, note.EndTime)
		assert.Equal(t, returnedRows[0].Longitude, note.Longitude)
		assert.Equal(t, returnedRows[0].Latitude, note.Latitude)
		assert.Equal(t, returnedRows[0].Id, note.Id)
		assert.Equal(t, (*returnedRows[0].Users)[0], (*note.Users)[0])
		assert.Equal(t, (*returnedRows[0].Tags)[0], (*note.Tags)[0])
	*/
}

func TestDeleteNote(t *testing.T) {
	testDelete("Notes", t, models.Notes.Delete)
}

func TestGetTimePeriodNotes(t *testing.T) {

	db, mock := initMockDB(t)
	defer db.Close()

	rows2 := generateTestRowsNoteUser()
	rows3 := generateTestRowsNoteTag()

	mock.ExpectQuery("SELECT (.)+ FROM notes (.)+ JOIN (.)+ WHERE \\(starttime <= (.)+ AND endtime >= (.)+\\)").
		WillReturnRows(rows2)

	mock.ExpectQuery("SELECT (.)+ FROM notes (.)+ LEFT JOIN (.)+ WHERE \\(starttime <= (.)+ AND endtime >= (.)+\\)").
		WillReturnRows(rows3)

	returnedRows, _ := models.Notes.GetActiveAtTime("2017-01-01 00:00")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfufilled expectations: %s", err)
	}
	if len(returnedRows) != 2 {
		t.Errorf("Function did not return correct number of rows. Returned %d rows", len(returnedRows))
		return
	}

	/*
		assert.Equal(t, returnedRows[0].Title, note.Title)
		assert.Equal(t, returnedRows[0].Comment, note.Comment)
		assert.Equal(t, returnedRows[0].StartTime, note.StartTime)
		assert.Equal(t, returnedRows[0].EndTime, note.EndTime)
		assert.Equal(t, returnedRows[0].Longitude, note.Longitude)
		assert.Equal(t, returnedRows[0].Latitude, note.Latitude)
		assert.Equal(t, returnedRows[0].Id, note.Id)
		assert.Equal(t, (*returnedRows[0].Users)[0], (*note.Users)[0])
		assert.Equal(t, (*returnedRows[0].Tags)[0], (*note.Tags)[0])
	*/

}

/*
func TestMergeNotes(t *testing.T) {

	models.InitDB()

	newtitle := "newnote"
	title := "testing title"
	comment := "testing comments"
	startTime := "2017-01-01 00:00"
	endTime := "2017-05-05 00:00"
	longitude := 1.0
	latitude := 2.0
	var id int64 = 1
	//email := "test@mapnotes.co.uk"
	users := []models.User{{Name: "Harry", Email: "beans@classic.com"}}
	tags := []string{"Harry"}

	note := generateTestNote(newtitle, comment, startTime, endTime, longitude, latitude, id, users, tags)
	note1 := generateTestNote(title, comment, startTime, endTime, longitude, latitude, id, users, tags)
	note2 := generateTestNote(title, comment, startTime, endTime, longitude, latitude, id, users, tags)
	note3 := generateTestNote(title, comment, startTime, endTime, longitude, latitude, id, users, tags)

	id1, _ := models.Notes.Create(&note1)
	id2, _ := models.Notes.Create(&note2)
	id3, _ := models.Notes.Create(&note3)

	ids := []int64{id1, id2, id3}

	models.Notes.Merge(ids, note)

}
*/

func TestInsertUser(t *testing.T) {
	db, mock := initMockDB(t)
	defer db.Close()

	var email string = "user@email.com"
	var name string = "Harry"
	var picture string = "bds"

	mock.ExpectPrepare("INSERT INTO users\\((.)+\\) VALUES\\(\\$1, \\$2, \\$3\\) RETURNING id").
		ExpectQuery().
		WithArgs(email, name, picture).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow("3"))

	models.Users.Create(&models.User{Email: email, Name: name, Picture: picture})
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("There were unfufilled expectations: %s", err)
	}

}

func TestGetNotesWithinRange(t *testing.T) {
	db, mock := initMockDB(t)
	defer db.Close()

	rows2 := generateTestRowsNoteUser()
	rows3 := generateTestRowsNoteTag()

	comment := "testing comments"
	title := "testing title"
	startTime := "2017-01-01 00:00"
	endTime := "2017-05-05 00:00"
	longitude := 1.0
	latitude := 2.0001
	var id int64 = 7

	users := []models.User{{Name: "Harry", Email: "beans.yeah@youwhat.not"}}
	tags := []string{"Harry"}

	note := generateTestNote(title, comment, startTime, endTime, longitude, latitude, id, users, tags)

	mock.ExpectQuery("SELECT (.)+ FROM notes (.)+ JOIN").
		WillReturnRows(rows2)

	mock.ExpectQuery("SELECT (.)+ FROM notes (.)+ LEFT JOIN").
		WillReturnRows(rows3)

	returnedNotes, err := models.GetNotesWithinRange(50, note)
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfufilled expectations: %s", err)
		return
	}
	if len(returnedNotes) < 1 {
		t.Errorf("Function did not return correct number of rows. Returned %d rows", len(returnedNotes))
		return
	}

	assert.Equal(t, returnedNotes[0].Title, note.Title)
	assert.Equal(t, returnedNotes[0].Comment, note.Comment)
	assert.Equal(t, returnedNotes[0].StartTime, note.StartTime)
	assert.Equal(t, returnedNotes[0].EndTime, note.EndTime)
	assert.Equal(t, returnedNotes[0].Longitude, note.Longitude)
	assert.Equal(t, returnedNotes[0].Latitude, note.Latitude)
	assert.Equal(t, returnedNotes[0].Id, note.Id)
	assert.Equal(t, (*returnedNotes[0].Users)[0].Email, (*note.Users)[0].Email)
	assert.Equal(t, (*returnedNotes[0].Tags)[0], (*note.Tags)[0])

}

func TestGetSimilarTitles(t *testing.T) {
	var c string = "Goodbye World"
	var s = "test_start"
	var e = "test_end"
	var long float64 = 1.0
	var lat float64 = 2.0
	var id int64 = -1
	var user models.User
	user.Name = "beans"
	users := []models.User{user}
	tags := []string{}

	title1 := "Hello"
	title2 := "hello"
	title3 := "goodbye"

	note1 := models.Note{Title: &title1, Comment: &c, StartTime: &s, EndTime: &e,
		Longitude: &long, Latitude: &lat, Id: &id, Users: &users, Tags: &tags}

	note2 := models.Note{Title: &title2, Comment: &c, StartTime: &s, EndTime: &e,
		Longitude: &long, Latitude: &lat, Id: &id, Users: &users, Tags: &tags}

	note3 := models.Note{Title: &title3, Comment: &c, StartTime: &s, EndTime: &e,
		Longitude: &long, Latitude: &lat, Id: &id, Users: &users, Tags: &tags}

	var notes []models.Note = make([]models.Note, 3)
	notes[0] = note1
	notes[1] = note2
	notes[2] = note3

	filtered := models.GetNotesWithSimilarText(notes)

	// Filtered should contain only note1 and note2 since both are the same
	// title, just one has a title format while the other is lowercase
	assert.Equal(t, len(filtered), 2)
	assert.Equal(t, filtered[0], note1)
	assert.Equal(t, filtered[1], note2)
}

func TestGetSimilarTags1(t *testing.T) {

	var title string = "Test title"
	var c string = "test comment"
	var s = "test_start"
	var e = "test_end"
	var long float64 = 1.0
	var lat float64 = 2.0
	var id int64 = -1
	var user models.User
	user.Name = "u11"
	users := []models.User{user}
	tags1 := []string{"tag1", "tag2", "tag3"}

	note1 := models.Note{Title: &title, Comment: &c, StartTime: &s, EndTime: &e,
		Longitude: &long, Latitude: &lat, Id: &id, Users: &users, Tags: &tags1}

	tags2 := []string{"tag1", "tag4"}

	note2 := models.Note{Title: &title, Comment: &c, StartTime: &s, EndTime: &e,
		Longitude: &long, Latitude: &lat, Id: &id, Users: &users, Tags: &tags2}

	tags3 := []string{"tag5", "tag6"}

	note3 := models.Note{Title: &title, Comment: &c, StartTime: &s, EndTime: &e,
		Longitude: &long, Latitude: &lat, Id: &id, Users: &users, Tags: &tags3}

	var notes []models.Note = make([]models.Note, 3)
	notes[0] = note1
	notes[1] = note2
	notes[2] = note3

	filtered, err := models.GetNotesWithSimilarTags(notes)

	if err != nil {
		t.Errorf("Function: GetNotesWithSimilarTags threw an error: %s", err)
	}

	//filtered should contain only note1 and note2 since both are similar
	assert.Equal(t, len(filtered), 2)
	assert.Equal(t, filtered[0], note1)
	assert.Equal(t, filtered[1], note2)
}

func TestGetSimilarTags2(t *testing.T) {

	var title string = "Test title"
	var c string = "test comment"
	var s = "test_start"
	var e = "test_end"
	var long float64 = 1.0
	var lat float64 = 2.0
	var id int64 = -1
	var user models.User
	user.Name = "u11"
	users := []models.User{user}

	tags1 := []string{"tag1", "tag2", "tag3"}
	note1 := models.Note{Title: &title, Comment: &c, StartTime: &s, EndTime: &e,
		Longitude: &long, Latitude: &lat, Id: &id, Users: &users, Tags: &tags1}

	tags2 := []string{"tag1", "tag4"}
	note2 := models.Note{Title: &title, Comment: &c, StartTime: &s, EndTime: &e,
		Longitude: &long, Latitude: &lat, Id: &id, Users: &users, Tags: &tags2}

	tags3 := []string{"tag5", "tag6"}
	note3 := models.Note{Title: &title, Comment: &c, StartTime: &s, EndTime: &e,
		Longitude: &long, Latitude: &lat, Id: &id, Users: &users, Tags: &tags3}

	tags4 := []string{"tag6", "tag7", "tag8"}
	note4 := models.Note{Title: &title, Comment: &c, StartTime: &s, EndTime: &e,
		Longitude: &long, Latitude: &lat, Id: &id, Users: &users, Tags: &tags4}

	tags5 := []string{"tag1", "tag6"}
	note5 := models.Note{Title: &title, Comment: &c, StartTime: &s, EndTime: &e,
		Longitude: &long, Latitude: &lat, Id: &id, Users: &users, Tags: &tags5}

	var notes []models.Note = make([]models.Note, 5)
	notes[0] = note1
	notes[1] = note2
	notes[2] = note3
	notes[3] = note4
	notes[4] = note5

	filtered, err := models.GetNotesWithSimilarTags(notes)

	if err != nil {
		t.Errorf("Function: GetNotesWithSimilarTags threw an error: %s", err)
	}

	//filtered should contain only note3, note4 and note2 because tag6 is the tag
	// that occurs the most across all the notes.
	assert.Equal(t, len(filtered), 3)
	assert.Equal(t, filtered[0], note3)
	assert.Equal(t, filtered[1], note4)
	assert.Equal(t, filtered[2], note5)
}

func TestGetSimilarTags3(t *testing.T) {

	var title string = "Test title"
	var c string = "test comment"
	var s = "test_start"
	var e = "test_end"
	var long float64 = 1.0
	var lat float64 = 2.0
	var id int64 = -1
	var user models.User
	user.Name = "u11"
	users := []models.User{user}

	tags1 := []string{"tag1"}
	note1 := models.Note{Title: &title, Comment: &c, StartTime: &s, EndTime: &e,
		Longitude: &long, Latitude: &lat, Id: &id, Users: &users, Tags: &tags1}

	tags2 := []string{"tag2"}
	note2 := models.Note{Title: &title, Comment: &c, StartTime: &s, EndTime: &e,
		Longitude: &long, Latitude: &lat, Id: &id, Users: &users, Tags: &tags2}

	tags3 := []string{"tag3"}
	note3 := models.Note{Title: &title, Comment: &c, StartTime: &s, EndTime: &e,
		Longitude: &long, Latitude: &lat, Id: &id, Users: &users, Tags: &tags3}

	var notes []models.Note = make([]models.Note, 3)
	notes[0] = note1
	notes[1] = note2
	notes[2] = note3

	filtered, err := models.GetNotesWithSimilarTags(notes)

	if err != nil {
		t.Errorf("Function: GetNotesWithSimilarTags threw an error: %s", err)
	}

	//filtered should not contain any notes,because no notes have the same tags
	assert.Equal(t, len(filtered), 0)
}

type DeleteFunc func(int64) error

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

func TestConstructAggregatedNote(t *testing.T) {
	n1Users := []models.User{generateTestUser("Harry", "harry@harrysworld.com", 1, "blah")}
	n1Tags := []string{"Harry", "beans"}
	n1 := generateTestNote("First note", "", "2017-11-18 10:10", "2017-11-18 10:50", 1.0, 2.01, 1, n1Users, n1Tags)

	n2Users := []models.User{generateTestUser("Beans", "beans@classic.com", 2, "blah")}
	n2Tags := []string{"heinz", "beans"}
	n2 := generateTestNote("Second note", "Beans baby", "2017-11-18 10:00", "2017-11-18 11:00", 1.01, 2.02, 2, n2Users, n2Tags)

	n3Users := []models.User{generateTestUser("Bill", "thescienceguy@gmail.com", 3, "blah")}
	n3Tags := []string{}
	n3 := generateTestNote("Third note", "Science rules!", "2017-11-18 09:50", "2017-11-18 11:05", 1.02, 2.03, 3, n3Users, n3Tags)

	users := append(n1Users, n2Users...)
	users = append(users, n3Users...)
	tags := []string{"Harry", "beans", "heinz"}

	notes := []models.Note{n1, n2, n3}
	ids, result := models.ConstructAggregatedNote(notes)

	assert.True(t, containsId(ids, 1), "Missing note id 1")
	assert.True(t, containsId(ids, 2), "Missing note id 2")
	assert.True(t, containsId(ids, 3), "Missing note id 3")

	assert.Equal(t, *n1.Title, *result.Title)
	assert.Equal(t, *n2.Comment, *result.Comment)
	assert.Equal(t, "2017-11-18 10:00", *result.StartTime)
	assert.Equal(t, "2017-11-18 11:00", *result.EndTime)
	assert.True(t, math.Abs(1.01-*result.Longitude) < TOLERANCE)
	assert.True(t, math.Abs(2.02-*result.Latitude) < TOLERANCE)
	assert.True(t, reflect.DeepEqual(users, *result.Users))
	assert.True(t, reflect.DeepEqual(tags, *result.Tags))
}

//------------------------------------------------------------------------------------------------------------------------------
//Test helper methods

func initMockDB(t *testing.T) (db *sql.DB, mock sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	models.SetDB(db)
	return
}

func containsId(ids []int64, id int64) bool {
	for _, el := range ids {
		if el == id {
			return true
		}
	}
	return false
}

func generateTestUser(name string, email string, id int64, picture string) models.User {
	return models.User{Name: name, Email: email, Id: id, Picture: picture}
}

func generateTestNote(title string, comment string, startTime string,
	endTime string, longitude float64, latitude float64, id int64, users []models.User, tags []string) (note models.Note) {

	note = models.Note{
		Title:     &title,
		Comment:   &comment,
		StartTime: &startTime,
		EndTime:   &endTime,
		Longitude: &longitude,
		Latitude:  &latitude,
		Id:        &id,
		Users:     &users,
		Tags:      &tags,
	}

	return

}

func generateTestRows() (rows *sqlmock.Rows, note models.Note) {

	title := "testing title"
	comment := "testing comments"
	startTime := "2017-01-01 00:00"
	endTime := "2017-05-05 00:00"
	longitude := 1.0
	latitude := 2.0
	var id int64 = 1
	//email := "test@mapnotes.co.uk"
	users := []models.User{{Name: "Harry", Email: "beans@classic.com"}}
	tags := []string{"Harry"}

	note = generateTestNote(title, comment, startTime, endTime, longitude, latitude, id, users, tags)

	rows = sqlmock.NewRows([]string{"comments", "title", "n.id", "startTime", "endTime",
		"longitude", "latitude", "users", "tag"}).
		AddRow(comment, title, id, startTime, endTime, longitude, latitude, "Harry", "Harry").
		AddRow("Harry's world", "Hi Harry", 2, "2017-01-01 00:00", "2017-05-05 00:00", 3.0, 2.0, "Beans", "Beans")
	return
}

func generateTestRowsNoteUser() (rows *sqlmock.Rows) {

	comment := "testing comments"
	title := "testing title"
	noteid := 7
	startTime := "2017-01-01 00:00"
	endTime := "2017-05-05 00:00"
	longitude := 1.0
	latitude := 2.0
	userid := 1
	name := "Harry"
	email := "beans.yeah@youwhat.not"

	rows = sqlmock.NewRows([]string{"comments", "title", "id", "startTime", "endTime",
		"longitude", "latitude", "id", "name", "email"}).
		AddRow(comment, title, noteid, startTime, endTime, longitude, latitude, userid, name, email).
		AddRow("Harry's world", "Hi Harry", 2, "2017-01-01 00:00", "2017-05-05 00:00", 3.0, 2.0, 1.0, "Beans", "Beans@beans")
	return
}

func generateTestRowsNoteTag() (rows *sqlmock.Rows) {

	tag := "Harry"
	noteid := 7

	rows = sqlmock.NewRows([]string{"id", "tag"}).
		AddRow(noteid, tag)
		//AddRow("Harry's world", "Hi Harry", 2, "2017-01-01 00:00", "2017-05-05 00:00", 1.0, 2.0, "Beans", "Beans")
	return
}
