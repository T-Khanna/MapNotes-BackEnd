package models

import (
	_ "github.com/lib/pq"
	"log"
	"sync"

)

type User struct {
	Email string
}

type SynchronisedMap struct {

	sync.RWMutex
	usermap map[string]int64

}

var user_map_sync =  SynchronisedMap{usermap: make(map[string]int64)}


type UserOperations struct {
	Create func(*User) (error, int64)
	Delete func(string) error
}

var Users = UserOperations{
	Create: createUser,
	Delete: deleteUser,
}

func checkUserMap(email string) (exists bool, id int64) {

user_map_sync.RLock()
id, exists = user_map_sync.usermap[email]
user_map_sync.Unlock()

return

}

func insertUserMap(email string, id int64) {

  user_map_sync.Lock()
	user_map_sync.usermap[email] = id
	user_map_sync.Unlock()

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
