package main

import (
	"fmt"
	"gitlab.doc.ic.ac.uk/g1736215/MapNotes/handlers"
	"gitlab.doc.ic.ac.uk/g1736215/MapNotes/models"
	"github.com/julienschmidt/httprouter"
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
	router := initRouter()

	//Starts the server at designated port
	err := http.ListenAndServe(":"+string(port), router)

	if err != nil {
		log.Fatal(err.Error())
	}
}

func initRouter() http.Handler  {

	router := httprouter.New()

	router.GET("/", handlers.Handler)
	router.GET("/note", handlers.NoteHandler)
	router.GET("/allnotes", handlers.GroupNotesHandler)
	router.GET("/user", handlers.UserHandler)

	return router

}
