package handlers

import (
	"encoding/json"
	"fmt"
	//"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"gitlab.doc.ic.ac.uk/g1736215/MapNotes/middlewares"
	"gitlab.doc.ic.ac.uk/g1736215/MapNotes/models"
)

func decodeCommentStruct(r *http.Request) (error, *models.Comment) {
	var comment models.Comment
	decodeErr := json.NewDecoder(r.Body).Decode(&comment)

	if decodeErr != nil {
		return decodeErr, nil
	}
	user := r.Context().Value(middlewares.UserContextKey{}).(models.User)
	comment.User = user
	comment.User.Id = -1

	return nil, &comment
}

/*
 Route: GET /api/comments/:note_id
 Gets the Note with the specified user email
*/
func CommentsGetByNote(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	comments, err := models.Comments.GetByNote(id)

	if err != nil {
		logAndRespondWithError(
			w,
			err.Error(),
			fmt.Sprintf("Error: Database failed to retrieve comments for note id %s", note_id),
		)
		return
	}

	respondWithJson(w, struct{ Comments []models.Comment }{comments}, http.StatusOK)
}

func CommentsCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decodeErr, comment := decodeCommentStruct(r)

	if decodeErr != nil {
		logAndRespondWithError(
			w,
			"Error: Could not decode JSON body into comment struct.",
			decodeErr.Error(),
		)
		return
	}

	createErr := models.Comments.Create(*comment)
	if createErr != nil {
		logAndRespondWithError(
			w,
			"Error: Could not insert comment into database.",
			createErr.Error(),
		)
		return
	}

	respondWithJson(w, "success", http.StatusOK)

}
