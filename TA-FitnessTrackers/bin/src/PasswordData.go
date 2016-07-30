package main

import (
	"encoding/json"
	"log"

	"github.com/AndyNortrup/GoSplunk"
	"golang.org/x/oauth2"
)

type User struct {
	username string
	userID   string
	oauth2.Token
}

// getTokens gets a list of tokens that are in the storage/passwords endpoint
// for the given strategy
func getUsers(serverURI, sessionKey, strategy string) []User {

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
		collection, "fittness_for_splunk", "nobody", sessionKey)
	if err != nil {
		log.Fatalf("Unable to get user tokens from Splunk: %v\n", err)
	}
	defer tokenReader.Close()

	var result []User
	var username string

	//Temporary struct so we can get string values out then make a JSON token
	// by properly converting the date stamp
	type tempUser struct {
		Name  string    `json:"name"`
		Id    string    `json:"id"`
		Token tokenData `json:"token"`
	}

	type tokenData struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		Expires      string `json:"expires_at"`
		Scope        string `json:"scope"`
	}

	temp := []tempUser{}
	decode := json.NewDecoder(tokenReader)
	err = decode.Decode(temp)
	if err != nil {
		log.Fatalf("Failed to decode passwords from storage/passwords: %v\n JSON to Decode: %v\n", err, tokenJSON)
	}
	token := newToken(temp.RefreshToken, temp.AccessToken, temp.Expires, temp.TokenType)
	result = append(result, User{username: username, id: Id, Token: *token})
}
