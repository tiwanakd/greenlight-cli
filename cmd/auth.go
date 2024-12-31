package cmd

import (
	"encoding/json"
	"os"
	"time"
)

// Crate a struch to handle the authorization tokens.
// Authoriazation token will be stored in a local hidden json file.
// This will allow users login session to persist as long as the token is valid
type authToken struct {
	Token  string    `json:"token"`
	Expiry time.Time `json:"expiry"`
}

// add the token to the json file
func addAuthTokenToFile(body []byte) error {
	file, err := os.OpenFile(".auth_token.json", os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(body)
	if err != nil {
		return err
	}

	return nil
}

// get the token from the file to pass as the Authoriztion header as needed
func extractAuthToken() (authToken, error) {
	body, err := os.ReadFile(".auth_token.json")
	if err != nil {
		return authToken{}, err
	}

	//since the token is stored as { "authentication_token": {"token":"token"...
	//create a map to decode the data from the file into. "authentication_token" will
	//be created as a map key and token and expiry will be docoded to authToken struct
	tokenMap := make(map[string]authToken)
	err = json.Unmarshal(body, &tokenMap)
	if err != nil {
		return authToken{}, err
	}

	return tokenMap["authentication_token"], nil
}
