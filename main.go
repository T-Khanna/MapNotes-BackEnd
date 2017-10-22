package main

import (
	"encoding/json"
	"fmt"
	"gitlab.doc.ic.ac.uk/g1736215/MapNotes/models"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Page struct {
	Title string
	Body  []byte
}

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
func groupNotesHandler(w http.ResponseWriter, r *http.Request) {

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
func noteHandler(w http.ResponseWriter, r *http.Request) {

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

		//For now we are just passing in a dummy string
		//Eventually we will pass in the id
		//This will be changed once the DeleteNote query and its test have been changed
		models.DeleteNote(deleteID.Id)

	default:
		http.Error(w, "Invalid request method.", 405)
	}
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

//Handler for the '/' path
func handler(w http.ResponseWriter, r *http.Request) {
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

func handleRequests() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	//Connect to database
	models.InitDB()
	fmt.Println("Connected to database.")

	fmt.Println("Starting server on port", port)

	//Tells the server how to handle paths that equal the first arg
	http.HandleFunc("/", handler)
	http.HandleFunc("/note", noteHandler)
	http.HandleFunc("/allnotes", groupNotesHandler)

	//Starts the server at designated port
	http.ListenAndServe(":"+string(port), nil)
}

func main() {
	handleRequests()
}
