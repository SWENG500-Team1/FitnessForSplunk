package main

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/AndyNortrup/GoSplunk"
	"golang.org/x/oauth2"
)

type User struct {
	Name         string `json:"name"`
	UserID       string `json:"id"`
	Scope        []string
	oauth2.Token `json:"token"`
	TokenExpiry  string `json:"token_expiry"`
}

// getTokens gets a list of tokens that are in the storage/passwords endpoint
// for the given strategy
func getUsers(serverURI, sessionKey, strategy string) ([]User, error) {

	collection := getKVStoreCollection(strategy)

	tokenReader, err := splunk.KVStoreGetCollection(splunk.LocalSplunkMgmntURL,
		collection, fitnessForSplunkAppName, "nobody", sessionKey)
	if err != nil {
		return []User{}, errors.New("Unable to get user tokens from Splunk:" + err.Error())
	}
	defer tokenReader.Close()
	var users []User
	decoder := json.NewDecoder(tokenReader)
	err = decoder.Decode(&users)
	if err != nil {
		return users, err
	}

	for i, user := range users {
		if user.TokenExpiry != "" && strategy == strategyGoogle {
			users[i].Token.Expiry, err = time.Parse(getTokenTimeFormat(strategy), user.TokenExpiry)
			if err != nil {
				log.Printf("Failed to convert token expiry time to time.Time obj: %v", err)
			}
		}
	}
	return users, err
}

//updateKVStoreToken updates the Splunk KV Store with a new refreshToken for the user
func updateKVStoreToken(u User, strategy, sessionKey string) error {
	return splunk.KVStoreUpdate(splunk.LocalSplunkMgmntURL,
		getKVStoreCollection(strategy),
		u.UserID, u,
		fitnessForSplunkAppName,
		"nobody",
		sessionKey)
}

func getKVStoreCollection(strategy string) string {
	switch {
	case strategy == strategyGoogle:
		return "google_tokens"
	case strategy == strategyFitbit:
		return "fitbit_tokens"
	case strategy == strategyMicrosoft:
		return "microsoft_tokens"
	}
	return ""
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
