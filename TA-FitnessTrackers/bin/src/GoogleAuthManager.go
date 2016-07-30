package main

import (
	"log"
	"net/http"
	"time"

	"google.golang.org/api/fitness/v1"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func getTokenFromAccessCode(accessCode string, conf oauth2.Config) *oauth2.Token {
	tok, err := conf.Exchange(oauth2.NoContext, accessCode)
	if err != nil {
		log.Fatalf("Error fettching token: %v\n", err)
	}
	return tok
}

func newToken(
	refreshToken string,
	accessToken string,
	expiryStr string,
	tokenType string) *oauth2.Token {

	expires, err := time.Parse(oauth_time_format, expiryStr)
	if err != nil {
		log.Fatalf("Error fetting token from refresh token: %v\n", err)
	}

	tok := new(oauth2.Token)
	tok.RefreshToken = refreshToken
	tok.AccessToken = accessToken
	tok.Expiry = expires
	tok.TokenType = tokenType

	return tok
}

func getClient(tok *oauth2.Token, clientID string, clientSecret string) *http.Client {
	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes: []string{fitness.FitnessActivityReadScope,
			fitness.FitnessBodyReadScope},
		Endpoint: google.Endpoint,
	}

	client := conf.Client(oauth2.NoContext, tok)
	return client
}
