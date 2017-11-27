package models

import (
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type Image struct {
	URL    string
	Id     int64
	NoteId int64
}

type ImageOperations struct {
	GetByNote func(int64) ([]Image, error)
	Create    func(Image) error
	Merge     func([]int64, int64) error
	//Delete    func(string) error
}

var Images = ImageOperations{
	GetByNote: getImagesByNoteId,
	Create:    createImage,
	Merge:     mergeImages,
	//Delete:    deleteImage,
}

func mergeImages(oldnoteids []int64, newnoteid int64) (err error) {
	var idString string = ConvertIntArrayToString(oldnoteids)
	q1 := fmt.Sprintf("UPDATE images SET note_id = %d WHERE note_id in %s", newnoteid, idString)
	_, err = db.Exec(q1)
	return
}

func getImagesByNoteId(note_id int64) ([]Image, error) {
	images := make([]Image, 0)

	log.Println("Attempting to retrieve images with id ", note_id)
	rows, err := db.Query(`SELECT url, images.id, note_id
                         FROM images WHERE note_id = $1`, note_id)
	if err != nil {
		log.Println(err)
		return images, err
	}

	defer rows.Close()
	for rows.Next() {
		var image Image
		err = rows.Scan(&image.URL, &image.Id, &image.NoteId)
		if err != nil {
			log.Println(err)
			return images, err
		}
		images = append(images, image)
	}
	return images, err
}

func createImage(image Image) error {
	stmt, err := db.Prepare("INSERT INTO images(url, note_id) VALUES ($1, $2)")

	if err != nil {
		log.Println(err)
		return err
	}

	_, err = stmt.Exec(image.URL, image.NoteId)

	if err != nil {
		log.Println(err)
		return err
	}

	return err
}
