package main

import (
	"encoding/json"
	"errors"

	"github.com/AndyNortrup/GoSplunk"
	"golang.org/x/oauth2"
)

type User struct {
	name   string
	userID string
	scope  []string
	oauth2.Token
}

// getTokens gets a list of tokens that are in the storage/passwords endpoint
// for the given strategy
func getUsers(serverURI, sessionKey, strategy string) ([]User, error) {

	var collection string
	switch {
	case strategy == STRATEGY_GOOGLE:
		collection = "google_tokens"
	case strategy == STRATEGY_FITBIT:
		collection = "fitbit_tokens"
	case strategy == STRATEGY_MICROSOFT:
		collection = "microsoft_tokens"
	}

	tokenReader, err := splunk.KVStoreGetCollection(splunk.LocalSplunkMgmntURL,
		collection, fitnessForSplunkAppName, "nobody", sessionKey)
	if err != nil {
		return []User{}, errors.New("Unable to get user tokens from Splunk:" + err.Error())
	}
	defer tokenReader.Close()

	var result []User

	//Temporary struct so we can get string values out then make a JSON token
	// by properly converting the date stamp
	type tempUser struct {
		Name  string            `json:"name"`
		Id    string            `json:"id"`
		Token map[string]string `json:"token"`
	}

	type userCollection struct {
		users []tempUser
	}

	temp := &userCollection{}
	decode := json.NewDecoder(tokenReader)
	err = decode.Decode(temp)
	if err != nil {
		return []User{}, err
	}
	for _, tempUser := range temp.users {
		token := newToken(tempUser.Token["refresh_token"],
			tempUser.Token["access_token"],
			tempUser.Token["expires_at"],
			tempUser.Token["token_type"])
		result = append(result,
			User{name: tempUser.Name,
				userID: tempUser.Id,
				scope:  []string{tempUser.Token["scope"]},
				Token:  *token})
	}

	return result, nil
}
