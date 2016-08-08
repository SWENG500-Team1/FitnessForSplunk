package main

import (
	"encoding/json"
	"errors"

	"github.com/AndyNortrup/GoSplunk"
	"golang.org/x/oauth2"
)

type User struct {
	Name         string `json:"name"`
	UserID       string `json:"id"`
	Scope        []string
	oauth2.Token `json:"token"`
}

// getTokens gets a list of tokens that are in the storage/passwords endpoint
// for the given strategy
func getUsers(serverURI, sessionKey, strategy string) ([]User, error) {

	var collection string
	switch {
	case strategy == strategyGoogle:
		collection = "google_tokens"
	case strategy == strategyFitbit:
		collection = "fitbit_tokens"
	case strategy == strategyMicrosoft:
		collection = "microsoft_tokens"
	}

	tokenReader, err := splunk.KVStoreGetCollection(splunk.LocalSplunkMgmntURL,
		collection, fitnessForSplunkAppName, "nobody", sessionKey)
	if err != nil {
		return []User{}, errors.New("Unable to get user tokens from Splunk:" + err.Error())
	}
	defer tokenReader.Close()
	var user []User
	decoder := json.NewDecoder(tokenReader)
	err = decoder.Decode(&user)
	if err != nil {
		return user, err
	}
	return user, err
}

//Temporary struct so we can get string values out then make a JSON token
// by properly converting the date stamp
type KVStoreUser struct {
	Name  string            `json:"name"`
	Id    string            `json:"id"`
	Token map[string]string `json:"token"`
}

type UserCollection struct {
	users []KVStoreUser
}

func getTokenTimeFormat(strategy string) string {
	switch strategy {
	case strategyGoogle:
		return googleOauthTimeFormat
	case strategyFitbit:
		return fitbitOauthTimeFormat
	}
	return ""
}
