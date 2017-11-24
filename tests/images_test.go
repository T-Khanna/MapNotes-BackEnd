package tests

import (
	"github.com/stretchr/testify/assert"
	"gitlab.doc.ic.ac.uk/g1736215/MapNotes/models"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	//"log"
	"testing"
)

func TestGetImagesByNoteId(t *testing.T) {

	db, mock := initMockDB(t)
	defer db.Close()

	rows := generateTestRowsImages()
	mock.ExpectQuery("SELECT (.)+ FROM images WHERE note_id = (.)+").
		WithArgs(7).
		WillReturnRows(rows)

	images, _ := models.Images.GetByNote(7)
	//log.Println(images)

	assert.True(t, len(images) == 2, "Did not return correct number of images")

	assert.Equal(t, int64(7), images[0].NoteId)
	assert.Equal(t, int64(7), images[1].NoteId)
	assert.Equal(t, int64(1), images[0].Id)
	assert.Equal(t, int64(2), images[1].Id)
	assert.Equal(t, "url1.jpg", images[0].URL)
	assert.Equal(t, "url2.jpg", images[1].URL)

}

func TestInsertImage(t *testing.T) {

	db, mock := initMockDB(t)
	defer db.Close()

	testImage := models.Image{Id: 1, URL: "Innit.jpg", NoteId: 5}

	mock.ExpectPrepare("INSERT INTO images\\((.)+\\) VALUES \\(\\$1, \\$2\\)").
		ExpectExec().
		WithArgs(testImage.URL, testImage.NoteId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	models.Images.Create(testImage)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("There were unfufilled expectations: %s", err)
	}

}

//Test Helper Methods
//--------------------------------------------------------------

func generateTestRowsImages() (rows *sqlmock.Rows) {

	noteid := 7
	imageid1 := 1
	url1 := "url1.jpg"

	//noteid2 := 5
	imageid2 := 2
	url2 := "url2.jpg"

	rows = sqlmock.NewRows([]string{"url", "images.id", "note_id"}).
		AddRow(url1, imageid1, noteid).
		AddRow(url2, imageid2, noteid)
	return rows
}
