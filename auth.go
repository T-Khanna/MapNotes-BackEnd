package main

import (
  "log"

  "github.com/futurenda/google-auth-id-token-verifier"
  "gitlab.doc.ic.ac.uk/g1736215/MapNotes/models"
)

/*
Function takes a token and verifies the integrity of the token given
After verification, it returns a user struct with the relevant information
from the token
*/
func AuthToken(token string) (isAuthenticated bool, user models.User) {

 verify := googleAuthIDTokenVerifier.Verifier{}
 aud := "xxxxxx-yyyyyyy.apps.googleusercontent.com"

 isAuthenticated = false


 err := verify.VerifyIDToken(token, []string{aud,})


 if err != nil {
   log.Println(err.Error())
   return false, models.User{Email: ""}
 }

 claimSet, err := googleAuthIDTokenVerifier.Decode(token)

 if err != nil {
   log.Println(err.Error())
   return false, models.User{Email: ""}
 }

 // Get User email from claimSet
 email := claimSet.Email

 //Fill User struct
 user = models.User{Email: email}

 //Token has been authenticated
 isAuthenticated = true

 return isAuthenticated, user
}
