package middlewares

import (
	"context"
	"gitlab.doc.ic.ac.uk/g1736215/MapNotes/auth"
	"gitlab.doc.ic.ac.uk/g1736215/MapNotes/models"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type UserContextKey struct{}

/*
 Middleware that authenticates token before calling subsequent HTTP requests.
*/
func Authenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("login_token")
		log.Println(token)
		isAuthenticated, user := auth.AuthToken(token)
		isAuthenticated = true

		if !isAuthenticated {
			http.Error(w, "Token unauthenticated", http.StatusUnauthorized)
			return
		}
		//user = randomUser()
		ctx := context.WithValue(r.Context(), UserContextKey{}, user)
		rWithUser := r.WithContext(ctx)
		h.ServeHTTP(w, rWithUser)
	})
}

func randomUser() models.User {
	users := make([]models.User, 0)
	users = append(users, models.User{Name: "Beans man", Email: "beans@sdsds", Picture: "beans.jpg"})
	users = append(users, models.User{Name: "Harry", Email: "harry@hadsdysworld.com", Picture: "harry.jpg"})
	users = append(users, models.User{Name: "Bill Nye", Email: "thesrgrgegceguy@science.org", Picture: "bill.jpg"})
	rand.Seed(time.Now().UTC().UnixNano())
	index := rand.Intn(len(users))
	return users[index]
}
