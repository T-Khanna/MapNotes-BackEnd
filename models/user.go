package models

import (
	_ "github.com/lib/pq"
	"log"

)

type User struct {
	Email string
}

var user_map map[string]int64


type UserOperations struct {
	Create func(*User) (error, int64)
	Delete func(string) error
}

var Users = UserOperations{
	Create: createUser,
	Delete: deleteUser,
}

func checkUserMap(email string) (exists bool, id int64) {

id, exists = user_map[email]

return

}

func insertUserMap(email string, id int64) {

	user_map[email] = id

}

func getUserId(email string) (err error, id int64) {

  keyExists, id := checkUserMap(email)

	if !(keyExists) {

		var newuser User = User{Email: email}

		err, id = createUser(&newuser)

		if err == nil {

	  insertUserMap(email, id)

	  }

	}

	return

}

func createUser(user *User) (err error, id int64) {

	stmt, err := db.Prepare("INSERT INTO users(email) VALUES($1) RETURNING id")

	if err != nil {
		log.Println(err)
		return
	}

	err = stmt.QueryRow(user.Email).Scan(&id)

	if err != nil {
		log.Println(err)
		return
	}

	return

}

//Not a vital function, but here if a user did wish to delete their account
func deleteUser(email string) (err error) {

	stmt, err := db.Prepare("DELETE FROM users WHERE email = $1")

	if err != nil {
		log.Println(err)
		return
	}

	_, err = stmt.Exec(email)

	if err != nil {
		log.Println(err)
		return
	}

	return
}
