package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	"gitlab.doc.ic.ac.uk/g1736215/MapNotes/handlers"
	"gitlab.doc.ic.ac.uk/g1736215/MapNotes/middlewares"
	"gitlab.doc.ic.ac.uk/g1736215/MapNotes/models"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	initDatabase()
	fmt.Println("Connected to database.")

	router := initRouter()

	fmt.Println("Starting server on port", port)
	err := http.ListenAndServe(":"+string(port), router)

	if err != nil {
		log.Fatal(err.Error())
	}
}

func initRouter() http.Handler {
	router := httprouter.New()

	setupRoutes(router)
	routerWithMiddleware := mkHandlerWithMiddleware(router)

	return routerWithMiddleware
}

func initDatabase() {
	models.InitDB()
}

func setupRoutes(router *httprouter.Router) {
	// Notes
	router.GET("/note/:time", handlers.NotesGetByTime)
	router.GET("/allnotes", handlers.NotesGetAll)
	router.POST("/note", handlers.NotesCreate)
	router.DELETE("/note", handlers.NotesDelete)

	// Users
	router.GET("/user", handlers.UserHandler)
}

func mkHandlerWithMiddleware(router http.Handler) http.Handler {
	return alice.New(
		middlewares.Logger,
		middlewares.Timeout,
	).Then(router)
}
