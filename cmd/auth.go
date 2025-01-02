package cmd

import (
	"encoding/json"
	"errors"
	"net/http"
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

// try to open the auth file with the required permission, check for the error.
// if the file does not exist, return a customer error asking user to login
func openAuthFile(perm int) (*os.File, error) {
	file, err := os.OpenFile(".auth_token.json", perm, 0600)
	if err != nil {
		switch {
		case err.Error() == `open .auth_token.json: no such file or directory`:
			return nil, errors.New("authorization required, please login using [login] command")
		default:
			return nil, err
		}
	}
	return file, nil
}

// add the token to the json file
func addAuthTokenToFile(body []byte) error {
	//open file with create and write only permissions, this will ensure that if
	//the file does not exist it gets created
	file, err := openAuthFile(os.O_CREATE | os.O_WRONLY)
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
	file, err := openAuthFile(os.O_RDONLY)
	if err != nil {
		return authToken{}, err
	}
	defer file.Close()

	//create a json Decoder on the file
	dec := json.NewDecoder(file)

	//since the token is stored as { "authentication_token": {"token":"token"...
	//create a map to decode the data from the file into. "authentication_token" will
	//be created as a map key and token and expiry will be docoded to authToken struct
	tokenMap := make(map[string]authToken)
	err = dec.Decode(&tokenMap)
	if err != nil {
		return authToken{}, err
	}

	return tokenMap["authentication_token"], nil
}

func setAuthHeader() (http.Header, error) {
	authToken, err := extractAuthToken()
	if err != nil {
		return nil, err
	}
	if authToken.Token == "" || authToken.Expiry.Before(time.Now()) {
		return nil, errors.New("autorization token expired or not found; Please login again using [login] command")
	}

	//create an empty header map and add the authorizion token that will be passed onto the api calls
	authorizationHeader := http.Header{}
	authorizationHeader.Add("Authorization", "Bearer "+authToken.Token)

	return authorizationHeader, nil
}

// remove the auth file to allow the user to logout as needed
func removeAuth() error {
	//extract the authoriztion token and check if it is empty or expired
	authToken, err := extractAuthToken()
	if err != nil {
		return err
	}
	if authToken.Token == "" || authToken.Expiry.Before(time.Now()) {
		return errors.New("autorization token expired or not found; Please login again using [login] command")
	}

	file, err := openAuthFile(os.O_WRONLY)
	if err != nil {
		return err
	}
	defer file.Close()

	return os.Remove(".auth_token.json")
}
