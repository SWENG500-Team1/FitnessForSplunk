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

func getClient(tok *oauth2.Token,
	clientID string,
	clientSecret string,
	strategy string) *http.Client {

	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	switch {
	case strategy == STRATEGY_GOOGLE:
		conf.Endpoint = google.Endpoint
		conf.Scopes = []string{"https://www.googleapis.com/auth/fitness.activity.read",
			"https://www.googleapis.com/auth/fitness.body.read",
			"https://www.googleapis.com/auth/userinfo.profile"}
	case strategy == STRATEGY_FITBIT:
		conf.Endpoint = fitbit.Endpoint
		conf.Scopes = []string{"activity"}
	}

	client := conf.Client(oauth2.NoContext, tok)
	return client
}
