package main

import (
	"fmt"
	"gitlab.doc.ic.ac.uk/g1736215/MapNotes/handlers"
	"gitlab.doc.ic.ac.uk/g1736215/MapNotes/models"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	//Connect to database
	models.InitDB()
	fmt.Println("Connected to database.")

	fmt.Println("Starting server on port", port)

	//Tells the server how to handle paths that equal the first arg
	http.HandleFunc("/", handlers.Handler)
	http.HandleFunc("/note", handlers.NoteHandler)
	http.HandleFunc("/allnotes", handlers.GroupNotesHandler)
	http.HandleFunc("/user", handlers.UserHandler)

	//Starts the server at designated port
	http.ListenAndServe(":"+string(port), nil)
}
