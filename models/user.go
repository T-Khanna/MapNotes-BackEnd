package models

import (
	_ "github.com/lib/pq"
	"log"
	"sync"

)

type User struct {
	Id    int64
	Name  string
	Email string
}

type SynchronisedMap struct {
	sync.RWMutex
	usermap map[string]int64
}

var user_map_sync = SynchronisedMap{usermap: make(map[string]int64)}

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
	user_map_sync.RUnlock()

	return

}

func insertUserMap(email string, id int64) {

	user_map_sync.Lock()
	user_map_sync.usermap[email] = id
	user_map_sync.Unlock()

}

func GetUserId(u User) (err error, id int64) {

	email := u.Email
	name := u.Name

	keyExists, id := checkUserMap(email)

	if !(keyExists) {

		user, userErr := getUserByEmail(email)

		if userErr != nil {
			return userErr, -1
		}

		if user != nil && user.Id != -1 {

			insertUserMap(email, user.Id)
			return nil, user.Id

		}

		var newuser User = User{Name: name, Email: email}

		err, id = createUser(&newuser)

		if err == nil {

			insertUserMap(email, id)

		}

	}

	return

}

func createUser(user *User) (err error, id int64) {

	stmt, err := db.Prepare("INSERT INTO users(email, name) VALUES($1, $2) RETURNING id")

	if err != nil {
		log.Println(err)
		return
	}

	err = stmt.QueryRow(user.Email, user.Name).Scan(&id)

	if err != nil {
		log.Println(err)
		return
	}

	return

}

//Not a vital function, but here if a user did wish to delete their account
func deleteUser(email string) (err error) {

	//TODO: make this function deletr from the hashmap
	//TODO: set up cascade deletes

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

func getUserByEmail(email string) (user *User, err error) {
	rows, err := db.Query("SELECT id, name FROM Users WHERE email = $1", email)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	user = &User{Id: -1}
	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Name)
	}
	return user, err
}
