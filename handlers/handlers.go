package handlers

import (
	"encoding/json"
	"fmt"
	"gitlab.doc.ic.ac.uk/g1736215/MapNotes/models"
	"log"
	"net/http"
	"strconv"
)

type TimeKey struct {
	Time string
}

type IdKey struct {
	Id int64
}

type ReturnJSON struct {
	Notes []models.Note
}

//Handler for the '/allnotes' Path
//Used to get groups of notes
func GroupNotesHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "GET":
		decoder := json.NewDecoder(r.Body)
		var timejson TimeKey
		err := decoder.Decode(&timejson)
		var notes []models.Note

		if err != nil {

			fmt.Fprintf(w, "Incorrect format for getting notes in a time period")

		} else {

			notes = models.GetTimePeriodNotes(timejson.Time)
			json.NewEncoder(w).Encode(ReturnJSON{Notes: notes})

		}

	default:
		http.Error(w, "Invalid request method.", 405)

	}

}

//handler for the '/note' Path
//Used for inserting notes and deleting notes
func NoteHandler(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)

	switch r.Method {
	case "GET":
		fmt.Fprintf(w, "Hello world!123\n")
	case "POST":
		fmt.Println("Received POST to ", r.URL.Path)
		fmt.Println("Inserting note into database")

		var note models.Note
		err := decoder.Decode(&note)

		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, "Error could not decode POST request. Possibly incorrect JSON string\n")
			return
		}

		var id int64 = models.InsertNote(note)
		if id == -1 {
			fmt.Fprintf(w, "Database encountered an error. Failed to insert note.")
			return
		}
		fmt.Fprintf(w, strconv.FormatInt(id, 10)+"\n")

	case "DELETE":
		var deleteID IdKey
		err := decoder.Decode(&deleteID)

		if err != nil {
			fmt.Fprintf(w, "Incorrect format for deleting a note")
			return
		}

		models.DeleteNote(deleteID.Id)

	default:
		http.Error(w, "Invalid request method.", 405)
	}
}

//Handler for the '/' path
func Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	//middle argument is HTML string
	case "GET":
		fmt.Fprintf(w, "Hello world!123\n")
	case "POST":
		fmt.Fprintf(w, "You sent a post request to \"%s\"\n", r.URL.Path)
	default:
		http.Error(w, "Invalid request method.", 405)
	}
}

func UserHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "GET":
		fmt.Fprintf(w, "GET sent to /user")
	case "POST":
		fmt.Fprintf(w, "POST sent to /user")
	default:
		http.Error(w, "Invalid request method.", 405)

	}

}
