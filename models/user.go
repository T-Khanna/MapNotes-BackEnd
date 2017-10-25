package models

import (
	_ "github.com/lib/pq"
	"log"
)

type User struct {
	Email string
}

type UserOperations struct {
	Create func(*User) (int64)
	Delete func(string)
}

var Users = UserOperations{
	Create: createUser,
	Delete: deleteUser,
}

func createUser(user *User) (id int64) {

	stmt, err := db.Prepare("INSERT INTO users(email) VALUES($1)")

	if err != nil {
		log.Println(err)
		return -1
	}

	_, err = stmt.Exec(user.Email)

	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT max(id) FROM users")

	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		err = rows.Scan(&id)
	}

	return

}

//Not a vital function, but here if a user did wish to delete their account
func deleteUser(email string) {

	stmt, err := db.Prepare("DELETE FROM users WHERE email = $1")

	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(email)

	if err != nil {
		log.Fatal(err)
	}
}
