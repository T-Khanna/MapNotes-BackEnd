package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"gitlab.doc.ic.ac.uk/g1736215/MapNotes/middlewares"
	"gitlab.doc.ic.ac.uk/g1736215/MapNotes/models"
	validation "gitlab.doc.ic.ac.uk/g1736215/MapNotes/validation"
)

func decodeNoteStruct(r *http.Request) (error, *models.Note) {
	var note models.Note
	decodeErr := json.NewDecoder(r.Body).Decode(&note)

	if decodeErr != nil {
		return decodeErr, nil
	}
	user := r.Context().Value(middlewares.UserContextKey{}).(models.User)

	if note.Users != nil {
		*note.Users = append(*note.Users, user)
	} else {
		note.Users = &[]models.User{user}
	}

	return nil, &note
}

/*
 Route: GET /api/notes/user/:user_email
 Gets the Note with the specified user email
*/
func NotesGetByUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user_email := ps.ByName("user_email")

	if user_email == "" {
		msg := fmt.Sprintf("Error: Could not parse user_email param: %s", user_email)
		logAndRespondWithError(w, msg, msg)
		return
	}

	notes, err := models.Notes.GetByUser(user_email)

	if err != nil {
		logAndRespondWithError(
			w,
			err.Error(),
			fmt.Sprintf("Error: Database failed to retrieve notes active at user_email %s", user_email),
		)
		return
	}

	respondWithJson(w, struct{ Notes []models.Note }{notes}, http.StatusOK)
}

/*
 Route: GET /api/notes/time/:time
 Gets the Note with the specified time.
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

	respondWithJson(w, struct{ Notes []models.Note }{notes}, http.StatusOK)
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

	log.Println(r.Context().Value(middlewares.UserContextKey{}))

	respondWithJson(w, struct{ Notes []models.Note }{notes}, http.StatusOK)
}

/*
 Route: POST /api/notes
 Creates a new Note with attributes from the request body given in JSON format.
*/
func NotesCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Decode body into Note struct
	decodeErr, note := decodeNoteStruct(r)

	if decodeErr != nil {
		logAndRespondWithError(
			w,
			"Error: Could not decode JSON body into Note struct.",
			decodeErr.Error(),
		)
		return
	}

	if !(validation.ValidateNote(note)) {

		logAndRespondWithError(
			w,
			"Error: Note could not be validated.",
			"Error: Note could not be validated.",
		)
		return

	}

	// Create new Note
	newId, createErr := models.Notes.Create(note)

	if createErr != nil {
		logAndRespondWithError(
			w,
			"Error: Could not insert Note into database.",
			createErr.Error(),
		)
		return
	}
	merge := false
	if models.TimeForAggregation() {
		var RANGE float64 = 50
		//Put the newId in the note struct
		note.Id = &newId
		notes, err := models.GetNotesWithinRange(RANGE, *note)
		if err != nil {
			logAndRespondWithError(
				w,
				"Error: Failed to perform a filter of notes with a certain range",
				err.Error(),
			)
			return
		}
		notes, err = models.GetAllNotesAroundSameTime(notes)
		if err != nil {
			logAndRespondWithError(
				w,
				"Error: Failed to perform a filter of notes by a certain timeframe",
				err.Error(),
			)
			return
		}
		notes = models.GetNotesWithSimilarText(notes)
		notes, err = models.GetNotesWithSimilarTags(notes)
		if err != nil {
			logAndRespondWithError(
				w,
				"Error: Failed to perform an aggregation of notes by their tags",
				err.Error(),
			)
			return
		}
		if len(notes) == 0 {
			//No notes were similar so no need to continue
			log.Println("Did not find any similar notes")
			return
		}
		note_ids, note := models.ConstructAggregatedNote(notes)
		mergeId, mergeErr := models.Notes.Merge(note_ids, note)
		if mergeErr != nil {
			logAndRespondWithError(
				w,
				"Error: Could not insert merged Note into database.",
				mergeErr.Error(),
			)
			return
		}
		newId = mergeId
		merge = true
	}
	//TODO Tell the front-end to perform some sort of refresh to get rid of all
	//the deleted notes
	// Return { id: newId, Merge: merge} as JSON.
	respondWithJson(w, struct {
		Id    int64
		Merge bool
	}{newId, merge}, http.StatusCreated)

}

/*
 Route: PUT /api/notes
 Creates a new Note with attributes from the request body given in JSON format.
*/
func NotesUpdate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decodeErr, note := decodeNoteStruct(r)

	if decodeErr != nil {
		logAndRespondWithError(
			w,
			"Error: Could not decode JSON body into Note struct.",
			decodeErr.Error(),
		)
		return
	}

	if !validation.ValidatePartialNote(note) {
		msg := "Error: Invalid Note."
		logAndRespondWithError(w, msg, fmt.Sprintf("%s\n    Note = %+v", msg, *note))
		return
	}

	// Create new Note
	updateErr := models.Notes.Update(note)

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
