package main

import (
  "net/http"
  "log"
  "strconv"

  "google.golang.org/api/oauth2/v2"
  "gitlab.doc.ic.ac.uk/g1736215/MapNotes/models"

)
var client = &http.Client{}


/*
Function takes a token and verifies the integrity of the token given
After verification, it returns a user struct with the relevant information
from the token
*/
func authToken(token string) (models.User) {

	service, err := oauth2.New(client)
	tokenInfoCall := service.Tokeninfo()

  // Sets parameter of IdToken to token in tokenInfoCall
	tokenInfoCall.IdToken(token)

	tokenInfo, err := tokenInfoCall.Do()
  if err != nil {
			 log.Fatal(err.Error())
	 }

	//Get User Id from tokenInfo
	id_string := tokenInfo.UserId

	//Converts string to int
	id, err := strconv.Atoi(id_string)

	// Get User email from tokenInfo
	email := tokenInfo.Email

	//Fill User struct
	user := models.User{Userid: id, Username: email, Password: ""}

	return user

}
