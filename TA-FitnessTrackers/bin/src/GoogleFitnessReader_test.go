package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"os"
	"testing"
	"time"

	splunk "github.com/AndyNortrup/GoSplunk"

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

func TestGoogleStrategy(t *testing.T) {
	sessionToken, _ := splunk.NewSessionKey(accountName, password, splunk.LocalSplunkMgmntURL)
	config := &splunk.ModInputConfig{
		SessionKey: sessionToken.SessionKey,
		Stanzas:    []splunk.ModInputStanza{},
	}
	config.CheckpointDir = "/Users/andy"
	stanza := &splunk.ModInputStanza{Params: []splunk.ModInputParam{}}
	stanza.StanzaName = "TA-FitnessTrackers://google"
	param := &splunk.ModInputParam{Name: strategyParamName, Value: strategyGoogle}
	stanza.Params = append(stanza.Params, *param)
	config.Stanzas = append(config.Stanzas, *stanza)

	b, _ := xml.Marshal(config)
	input := &FitnessInput{reader: bytes.NewReader(b), writer: os.Stdout}
	input.StreamEvents()
}

//static start times and end times based on the author's dataSet
var startTime = time.Now().Add(-12 * time.Hour)
var endTime = startTime.Add(5 * time.Hour)

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
	client, _ := getClient(tok, testClientId, testClientSecret, strategyGoogle)
	latestTime := reader.getData(client, devNull, User{Name: "AndyNortrup"})

	if !latestTime.After(startTime) {
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
	client, _ := getClient(tok, testClientId, testClientSecret, strategyGoogle)
	err := reader.getSessions(
		client,
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
func TestExchangeAccessCode(t *testing.T) {
	conf := oauth2.Config{ClientID: testClientId, ClientSecret: testClientSecret}
	conf.Endpoint = google.Endpoint
	conf.Scopes = []string{"https://www.googleapis.com/auth/fitness.activity.read",
		"https://www.googleapis.com/auth/fitness.body.read"}
	conf.RedirectURL = "https://www.fitnessforsplunk.ninja:8000/en-US/app/fitness_for_splunk/Adding_Google_Account"

	tok := getTokenFromAccessCode("4/vBcdgk1KRn8r--b-_TPVqhV_LK4UU0jlC8dy_etZnO8#", conf)
	tokStr, err := json.Marshal(tok)
	if err != nil {
		t.Fail()
		t.Log(err)
	}
	t.Logf("%s\n", tokStr)
}
