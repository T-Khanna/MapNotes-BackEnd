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

type ReturnJSON struct {
	Notes []models.Note
}

func returnAllNotes(w http.ResponseWriter, r *http.Request) {
	notes := models.GetAllNotes()
	json.NewEncoder(w).Encode(ReturnJSON{Notes: notes})
}

func noteHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Fprintf(w, "Hello world!123\n")
	case "POST":
		fmt.Println("Received POST to ", r.URL.Path)
		fmt.Println("Inserting note into database")

		decoder := json.NewDecoder(r.Body)
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
	http.HandleFunc("/allnotes", returnAllNotes)

	//Starts the server at designated port
	http.ListenAndServe(":"+string(port), nil)
}

func main() {
	handleRequests()
}
