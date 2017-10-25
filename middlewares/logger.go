package middlewares

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

/*
 Middleware that logs all requests in the Apache Common Log Format.
*/
func Logger(h http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, h)
}
