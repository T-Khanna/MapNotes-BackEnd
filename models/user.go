package models

import (
	_ "github.com/lib/pq"
	"log"
)

type User struct {
	Email string
}

type UserOperations struct {
	Create func(*User) bool
	Delete func(string) bool
}

var Users = UserOperations{
	Create: createUser,
	Delete: deleteUser,
}

func createUser(user *User) (wasCreated bool) {

	stmt, err := db.Prepare("INSERT INTO users(email, username) VALUES($1, $2)")

	if err != nil {
		log.Println(err)
		return false
	}

	_, err = stmt.Exec(user.Email, user.Username)

	if err != nil {
		log.Println(err)
		return false
	}

	return true

}

//Not a vital function, but here if a user did wish to delete their account
func deleteUser(email string) bool {

	stmt, err := db.Prepare("DELETE FROM users WHERE email = $1")

	if err != nil {
		log.Println(err)
		return false
	}

	_, err = stmt.Exec(email)

	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
