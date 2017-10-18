package main

import (
	"fmt"
	"gitlab.doc.ic.ac.uk/g1736215/MapNotes/models"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"encoding/json"
	"github.com/gorilla/mux"
)

type Page struct {
	Title string
	Body  []byte
}


func returnAllNotes(w http.ResponseWriter, r *http.Request) {
		notes := models.GetAllNotes()
		json.NewEncoder(w).Encode(notes)
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
	// mux.Router matches incoming requests against a list of registered routes
	// and calls a handler for the route that matches the URL or other conditions.
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", handler)
	myRouter.HandleFunc("/notes", returnAllNotes)

	//Starts the server at designated port
	http.ListenAndServe(":"+string(port), myRouter)
}

func main() {
	handleRequests()
}
