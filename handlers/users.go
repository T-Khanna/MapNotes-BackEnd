package handlers

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func UserGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "GET sent to /user")
}

func UserCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "POST sent to /user")
}
