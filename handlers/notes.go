package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"gitlab.doc.ic.ac.uk/g1736215/MapNotes/models"
)

/*
 Route: GET /api/notes/:time
 Gets the Note with the specified id.
*/
func NotesGetByTime(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	time := ps.ByName("time")

	if time == "" {
		msg := fmt.Sprintf("Error: Could not parse time param: %s", time)
		logAndRespondWithError(w, msg, msg)
		return
	}

	notes, err := models.Notes.GetActiveAtTime(time)

	if err != nil {
		logAndRespondWithError(
			w,
			err.Error(),
			fmt.Sprintf("Error: Database failed to retrieve notes active at time %s", time),
		)
		return
	}

	respondWithJson(w, struct{Notes []models.Note}{notes}, http.StatusOK)
}

/*
 Route: GET /api/all/notes
*/
func NotesGetAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	notes, err := models.Notes.GetAll()

	if err != nil {
		logAndRespondWithError(
			w,
			err.Error(),
			"Error: Database failed to retrieve all notes.",
		)
		return
	}

	respondWithJson(w, struct{Notes []models.Note}{notes}, http.StatusOK)
}

/*
 Route: POST /api/notes
 Creates a new Note with attributes from the request body given in JSON format.
*/
func NotesCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Decode body into Note struct
	note := models.Note{}
	decodeErr := json.NewDecoder(r.Body).Decode(&note)

	if decodeErr != nil {
		logAndRespondWithError(
			w,
			"Error: Could not decode JSON body into Note struct.",
			decodeErr.Error(),
		)
		return
	}

	// Create new Note
	newId, createErr := models.Notes.Create(&note)

	if createErr != nil {
		logAndRespondWithError(
			w,
			"Error: Could not insert Note into database.",
			createErr.Error(),
		)
		return
	}

	// Return { id: newId } as JSON.
	respondWithJson(w, struct{ Id int64 }{newId}, http.StatusCreated)
}

/*
 Route: PUT /api/notes
 Creates a new Note with attributes from the request body given in JSON format.
*/
func NotesUpdate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Decode body into Note struct
	note := models.Note{}
	decodeErr := json.NewDecoder(r.Body).Decode(&note)

	if decodeErr != nil {
		logAndRespondWithError(
			w,
			"Error: Could not decode JSON body into Note struct.",
			decodeErr.Error(),
		)
		return
	}

	// Create new Note
	updateErr := models.Notes.Update(&note)

	if updateErr != nil {
		logAndRespondWithError(
			w,
			"Error: Could not update note.",
			updateErr.Error(),
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

/*
 Route: DELETE /api/notes/:id
*/
func NotesDelete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idStr := ps.ByName("id")
	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		logAndRespondWithError(
			w,
			fmt.Sprintf("Error: Could not parse id param: %s", idStr),
			err.Error(),
		)
		return
	}

	err = models.Notes.Delete(id)

	// TODO: Test insertion of bad ID
	if err != nil {
		logAndRespondWithError(
			w,
			fmt.Sprintf("Error: Could not delete note with id: %d", id),
			err.Error(),
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

/*
 Serialises the `object` and writes it to `w`, the HTTP response, with status
 code `statusCode`.
*/
func respondWithJson(w http.ResponseWriter, object interface{}, statusCode int) {
	// Set content type and status code.
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)

	// Serialise object and write to ResponseWriter.
	err := json.NewEncoder(w).Encode(object)

	// If error, return error instead.
	if err != nil {
		logAndRespondWithError(
			w,
			"Error: Failed to encode object as json",
			err.Error(),
		)
	}
}

/*
Logs the `logMsg` on the server and writes an error to the response using
`http.Error()`, with `responseMsg` as the message and `http.StatusBadRequest` as
the status code.
*/
// TODO: Should we response with an { error: responseMsg } JSON? What does http.Error actually respond with?
func logAndRespondWithError(w http.ResponseWriter, responseMsg string, logMsg string) {
	log.Println(logMsg)
	http.Error(w, responseMsg, http.StatusBadRequest)
}
