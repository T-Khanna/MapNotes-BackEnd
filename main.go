package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Page struct {
	Title string
	Body  []byte
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

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	//Connect to database
	_, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to database.")

	fmt.Println("Starting server on port", port)

	//Tells the server how to handle paths that equal the first arg
	http.HandleFunc("/", handler)

	//Starts the server at designated port
	http.ListenAndServe(":"+string(port), nil)

	fmt.Println("Started server on port", port)
}
