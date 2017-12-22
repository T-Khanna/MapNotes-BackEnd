package handlers

import (
	"encoding/json"
	"fmt"
	//"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	//"gitlab.doc.ic.ac.uk/g1736215/MapNotes/middlewares"
	"gitlab.doc.ic.ac.uk/g1736215/MapNotes/models"
)

func decodeImageStruct(r *http.Request) (error, *models.Image) {
	var image models.Image
	decodeErr := json.NewDecoder(r.Body).Decode(&image)

	if decodeErr != nil {
		return decodeErr, nil
	}
	//user := r.Context().Value(middlewares.UserContextKey{}).(models.User)
	return nil, &image
}

/*
 Route: GET /api/images/:note_id
 Gets the images with the specified note id
*/
func ImagesGetByNote(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	note_id := ps.ByName("note_id")
	id, err := strconv.ParseInt(note_id, 10, 64)

	if err != nil {
		logAndRespondWithError(
			w,
			fmt.Sprintf("Error: Could not parse note id param: %s", note_id),
			err.Error(),
		)
		return
	}

	images, err := models.Images.GetByNote(id)

	if err != nil {
		logAndRespondWithError(
			w,
			err.Error(),
			fmt.Sprintf("Error: Database failed to retrieve images for note id %s", note_id),
		)
		return
	}

	respondWithJson(w, struct{ Images []models.Image }{images}, http.StatusOK)
}

/*
 Route: POST /api/images
 Inserts the image
*/
func ImagesCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decodeErr, image := decodeImageStruct(r)

	if decodeErr != nil {
		logAndRespondWithError(
			w,
			"Error: Could not decode JSON body into image struct.",
			decodeErr.Error(),
		)
		return
	}

	createErr := models.Images.Create(*image)
	if createErr != nil {
		logAndRespondWithError(
			w,
			"Error: Could not insert image into database.",
			createErr.Error(),
		)
		return
	}

	respondWithJson(w, "success", http.StatusOK)

}

//Route: DELETE /api/images/:image_id
//Deletes the image
func ImagesDelete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	image_id := ps.ByName("image_id")

	id, err := strconv.ParseInt(image_id, 10, 64)
	if err != nil {
		logAndRespondWithError(
			w,
			fmt.Sprintf("Error: Could not parse image id param: %s", image_id),
			err.Error(),
		)
		return
	}

	deleteErr := models.Images.Delete(id)
	if deleteErr != nil {
		logAndRespondWithError(
			w,
			"Error: Could not delete image into database.",
			deleteErr.Error(),
		)
		return
	}

	respondWithJson(w, "success", http.StatusOK)

}
