package main

import (
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/fitbit"
	"golang.org/x/oauth2/google"
)

func getTokenFromAccessCode(accessCode string, conf oauth2.Config) *oauth2.Token {
	tok, err := conf.Exchange(oauth2.NoContext, accessCode)
	if err != nil {
		log.Fatalf("Error fettching token: %v\n", err)
	}
	return tok
}

func newTokenWithExpiry(
	refreshToken string,
	accessToken string,
	expiryStr string,
	tokenType string,
	timeFormat string) *oauth2.Token {

	expires, err := time.Parse(timeFormat, expiryStr)
	if err != nil {
		log.Fatalf("Error fetting token from refresh token: %v\n", err)
	}

	tok := new(oauth2.Token)
	tok.RefreshToken = refreshToken
	tok.AccessToken = accessToken
	tok.TokenType = tokenType
	tok.Expiry = expires

	return tok
}

func newTokenNoExpiry(refreshToken, accessToken, tokenType string) *oauth2.Token {
	tok := &oauth2.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    tokenType,
	}

	return tok
}

func getClient(tok *oauth2.Token,
	clientID string,
	clientSecret string,
	strategy string) (*http.Client, *oauth2.Token) {

	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	//Fitbit is silly and wants you to remove the old token, cus they can't ignore
	// an expired token in the presece of a refresh token
	if strategy == strategyFitbit {
		tok.AccessToken = ""
	}

	switch {
	case strategy == strategyGoogle:
		conf.Endpoint = google.Endpoint
		conf.Scopes = []string{"https://www.googleapis.com/auth/fitness.activity.read",
			"https://www.googleapis.com/auth/fitness.body.read",
			"https://www.googleapis.com/auth/userinfo.profile"}
	case strategy == strategyFitbit:
		conf.Endpoint = fitbit.Endpoint
		conf.Scopes = []string{"activity", "profile"}
	}

	client := conf.Client(oauth2.NoContext, tok)

	newToken, err := conf.TokenSource(oauth2.NoContext, tok).Token()
	if err != nil {
		log.Printf("Failed to capture updated token. %v", tok.RefreshToken)
	}

	return client, newToken
}
