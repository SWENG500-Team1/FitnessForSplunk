package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

/*
Create a file for testing with your secrets:
const testClientId string
const testClinetSecret string
const testRefreshToken string
const testAccessToken string
const testExpires string
const testTokenType string
*/

//static start times and end times based on the author's dataSet
var startTime = time.Unix(0, 1468305494396000000)
var endTime = startTime.Add(5 * time.Millisecond)

func TestLatestTime(t *testing.T) {

	//TODO: Replace hard coded values with pull from arguments
	// reader := NewFitnessReader(testClientId, testClientSecret)

	//TODO: Replace hard coded values with pull from storage/passwords
	/*TODO: Determine if the value from storage/passwords has a refresh token.
	  Yes: Refresh the existing token.
	  No: Get a refresh token and store new token
	*/
	tok := newTokenWithExpiry(testRefreshToken,
		testAccessToken,
		testExpires,
		testTokenType,
		getTokenTimeFormat(strategyGoogle))

	reader := &GoogleFitnessReader{startTime: startTime, endTime: endTime}
	devNull := bufio.NewWriter(ioutil.Discard)
	latestTime := reader.getData(
		getClient(tok, testClientId, testClientSecret, strategyGoogle),
		devNull,
		User{Name: "AndyNortrup"})

	if latestTime.Nanosecond() == startTime.Nanosecond() {
		t.Logf("Unable to pull new data from data sources.  "+
			"\tStart Time: %v\tLatest Time: %v\n",
			startTime,
			latestTime)
		t.Fail()
	}
}

func TestGetSessions(t *testing.T) {
	startTime := time.Date(2016, 07, 16, 04, 0, 0, 0, time.Local)
	tok := newTokenWithExpiry(testRefreshToken,
		testAccessToken,
		testExpires,
		testTokenType,
		getTokenTimeFormat(strategyGoogle))
	reader := &GoogleFitnessReader{startTime: startTime, endTime: startTime.Add(12 * time.Hour)}
	reader.username = "andrew.nortrup"
	err := reader.getSessions(
		getClient(tok,
			testClientId,
			testClientSecret,
			strategyGoogle),
		bufio.NewWriter(os.Stdout))

	if err != nil {
		t.Logf("Unable to get sessions: %v\n", err)
		t.Fail()
	}
}

//TestCreateGoogleAuthCodeURL creates a authorization URL.
func TestCreateGoogleAuthCodeURL(t *testing.T) {
	conf := oauth2.Config{ClientID: testClientId, ClientSecret: testClientSecret}
	conf.Endpoint = google.Endpoint
	conf.Scopes = []string{"https://www.googleapis.com/auth/fitness.activity.read",
		"https://www.googleapis.com/auth/fitness.body.read"}
	conf.RedirectURL = "https://www.fitnessforsplunk.ninja:8000/en-US/app/fitness_for_splunk/Adding_Google_Account"
	//print a url to go get an access code
	t.Logf("URL: %v\n", conf.AuthCodeURL("state",
		oauth2.AccessTypeOffline))
}

//This is a utility method for generating tokens for testing purposes.
func disabledTestExchangeAccessCode(t *testing.T) {
	conf := oauth2.Config{ClientID: testClientId, ClientSecret: testClientSecret}
	conf.Endpoint = google.Endpoint
	conf.Scopes = []string{"https://www.googleapis.com/auth/fitness.activity.read",
		"https://www.googleapis.com/auth/fitness.body.read"}
	conf.RedirectURL = "https://www.fitnessforsplunk.ninja:8000/en-US/app/fitness_for_splunk/Adding_Google_Account"

	tok := getTokenFromAccessCode("4/jPQ5pJSfhM2nN0ZwwNkUzT6Fmr-EKmejunjswa8LJuM#", conf)
	tokStr, err := json.Marshal(tok)
	if err != nil {
		t.Fail()
		t.Log(err)
	}
	t.Logf("%s\n", tokStr)
}
