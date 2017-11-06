package middlewares

import (
	"context"
	"gitlab.doc.ic.ac.uk/g1736215/MapNotes/auth"
	"log"
	"net/http"
)

type UserContextKey struct{}

/*
 Middleware that authenticates token before calling subsequent HTTP requests.
*/
func Authenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
    token := r.Header.Get("login_token")
    log.Println(token)
		isAuthenticated, user := auth.AuthToken(token)
		//isAuthenticated = true
    if (!isAuthenticated) {
			http.Error(w, "Token unauthenticated", http.StatusUnauthorized)
			return
    }
    ctx := context.WithValue(r.Context(), UserContextKey{}, user.Email)
    rWithUser := r.WithContext(ctx)
		h.ServeHTTP(w, rWithUser)
	})
}
