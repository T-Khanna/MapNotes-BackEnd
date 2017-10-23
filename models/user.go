package models

import (
	_ "github.com/lib/pq"
	"log"
)

func insertUser(user User) (id int64) {

	stmt, err := db.Prepare("INSERT INTO users(username, password) VALUES($1, $2)")

	if err != nil {
		log.Println(err)
		return -1
	}

	_, err = stmt.Exec(user.Username, user.Password)

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
func DeleteUser(id int64) {

	stmt, err := db.Prepare("DELETE FROM users WHERE id = $1")

	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(id)

	if err != nil {
		log.Fatal(err)
	}
}
