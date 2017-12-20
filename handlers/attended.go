package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"gitlab.doc.ic.ac.uk/g1736215/MapNotes/middlewares"
	"gitlab.doc.ic.ac.uk/g1736215/MapNotes/models"
)

func AttendedCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("note_id")
	note_id, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		logAndRespondWithError(
			w,
			fmt.Sprintf("Error: Could not parse note_id param: %s", id),
			err.Error(),
		)
		return
	}
	user := r.Context().Value(middlewares.UserContextKey{}).(models.User)
	user.Id = -1
	attended := models.Attended{User: user, NoteId: note_id}

	createErr := models.Attend.Create(attended)
	if createErr != nil {
		logAndRespondWithError(
			w,
			"Error: Could not insert attendance into database.",
			createErr.Error(),
		)
		return
	}

	respondWithJson(w, "success", http.StatusOK)

}
