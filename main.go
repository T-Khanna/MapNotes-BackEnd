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
	router.GET("/api/notes/time/:time", handlers.NotesGetByTime)
	router.GET("/api/notes/user/:user_email", handlers.NotesGetByUser)
	router.GET("/api/all/notes", handlers.NotesGetAll)
	router.PUT("/api/notes", handlers.NotesUpdate)
	router.POST("/api/notes", handlers.NotesCreate)
	router.DELETE("/api/notes/:id", handlers.NotesDelete)
	router.GET("/api/comments/:note_id", handlers.CommentsGetByNote)
	router.POST("/api/comments", handlers.CommentsCreate)
	router.GET("/api/images/:note_id", handlers.ImagesGetByNote)
	router.POST("/api/images", handlers.ImagesCreate)

	// Users
	router.GET("/api/users", handlers.UserGet)
	router.POST("/api/users", handlers.UserCreate)
}

func mkHandlerWithMiddleware(router http.Handler) http.Handler {
	return alice.New(
		middlewares.Logger,
		middlewares.Timeout,
		middlewares.Authenticate,
	).Then(router)
}
