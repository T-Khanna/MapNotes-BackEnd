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
	Delete    func(int64) error
	Update    func(int64, []string) error
}

var Images = ImageOperations{
	GetByNote: getImagesByNoteId,
	Create:    createImage,
	Merge:     mergeImages,
	Delete:    deleteImageByImageID,
	Update:    updateImages,
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

func deleteImageByImageID(image_id int64) error {
	stmt, err := db.Prepare("DELETE FROM images WHERE id = $1")
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = stmt.Exec(image_id)
	if err != nil {
		log.Println(err)
		return err
	}
	return err
}

func updateImages(note_id int64, image_urls []string) error {
	for i := 0; i < len(image_urls); i++ {
		image_url := image_urls[i]
		var n int64
		err := db.QueryRow("SELECT note_id FROM images WHERE url = $1", image_url).Scan(&n)
		if err == nil {
			continue
		}
		image := Image{URL: image_url, NoteId: note_id}
		err = createImage(image)
		if err != nil {
			return err
		}
	}
	return nil
}
