package auth

import (
  "log"
  "net/http"
  "encoding/json"
  "io/ioutil"

  "github.com/futurenda/google-auth-id-token-verifier"
  "gitlab.doc.ic.ac.uk/g1736215/MapNotes/models"
)

type AuthUser struct {
  Iss string `json:"iss"`
  Sub string `json:"sub"`
  Azp string `json:"azp"`
  Aud string `json:"aud"`
  Iat string `json:"iat"`
  Exp string `json:"exp"`
  Email string `json:"email"`
  Email_verified string `json:"email_verified"`
  Name string `json:"name"`
  Picture string `json:"picture"`
  Given_name string `json:"given_name"`
  Family_name string `json:"family_name"`
  Locale string `json:"locale"`
}


/*
Function takes a token and verifies the integrity of the token given
After verification, it returns a user struct with the relevant information
from the token
*/
func AuthToken(token string) (isAuthenticated bool, user models.User) {

 verify := googleAuthIDTokenVerifier.Verifier{}
 aud := "371478445903-l0qtjdbu45ci2bobb5lhm41svvcbjc0u.apps.googleusercontent.com"

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


/*
  Checks that function AuthToken returns same value as Google's tokeninfo endpoint
  Result of AuthToken_test:
  1. true, nil : Both AuthToken and tokeninfo endpoint return invalid token
  2. true, AuthUser : Both AuthToken and tokeninfo endpoint have vaildated the
                      token and given the same email address
  3. false, nil: a) tokeninfo endpoint gives an invalid token but AuthToken gives
                    a valid token or
                 b) Both give valid token but email address gives are different
*/
func AuthToken_test(token string) (bool, *AuthUser) {
  // Test token
  tokeninfo := "https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=" + token
  res, err := http.Get(tokeninfo)

  // Actual
  isAuthenticated, user := AuthToken(token)

  if(err != nil) {
    log.Println(err.Error())
  }

  // Google tokeninfo says: Invalid token received
  if(res.StatusCode != http.StatusOK) {
    //check that AuthToken returns false
    return isAuthenticated == false, nil
  }

  // Get AuthUser Struct from Google tokeninfo endpoint
  body, err := ioutil.ReadAll(res.Body)
  if(err != nil) {
    log.Println(err.Error)
  }
  auth_user := getAuthUser([]byte(body))

  // Check email is the same
  email_check := auth_user.Email == user.Email

  if(email_check == false) {
    return email_check, nil
  }
  return email_check, &auth_user
}

/*
  Get AuthUser from JSON repsonse
*/
func getAuthUser(body []byte) (AuthUser) {
    var auth_user AuthUser
    err := json.Unmarshal(body, &auth_user)
    if(err != nil){
       log.Println(err.Error())
    }
    return auth_user
}
