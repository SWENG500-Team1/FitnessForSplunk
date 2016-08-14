package main

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"
	"time"

	splunk "github.com/AndyNortrup/GoSplunk"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/fitbit"
)

// TestFitbitGetData is an integration test with the Fitbit servers.  It does
// require some setup to be run correctly.
//
// 1. Set the date range with complete days. (i.e. not today)
// 2. Set fitbitAccessToken, fitbitAccessToken, fitbitExpires constants in a file
//    that is not in your git repo.
// 3. Make sure you know the exact value of steps the first date and set the value
//   of expectedStesp to that value.
func TestFitbitGetData(t *testing.T) {
	const expectedSteps int = 14102
	endTime := time.Date(2016, time.August, 2, 02, 00, 0, 0, time.Local)
	startTime := time.Date(2016, time.August, 1, 23, 0, 0, 0, time.Local)

	reader, err := NewFitbitReader(startTime, endTime)
	if err != nil {
		t.Fail()
		t.Logf("Failed to create Fitbit Reader: %v", err)
	}

	sessionKey, _ := splunk.NewSessionKey(accountName, password, splunk.LocalSplunkMgmntURL)
	users, err := getUsers(splunk.LocalSplunkMgmntURL,
		sessionKey.SessionKey,
		strategyFitbit)

	users[0].Token.AccessToken = ""

	client, newToken := getClient(&users[0].Token, fitbitClientId, fitbitClientSecret, strategyFitbit)
	// buf := bytes.NewBuffer([]byte{})
	// writer := bufio.NewWriter(buf)

	users[0].Token = *newToken
	updateKVStoreToken(users[0], strategyFitbit, sessionKey.SessionKey)

	//Go get the data from fitbit
	date := reader.getData(client, bufio.NewWriter(os.Stdout), User{Name: "Andy"})

	if date.Day() != endTime.Day() && date.Hour() != endTime.Hour() {
		t.Fail()
		t.Logf("Wrong date returned")
	}

	// log.Printf("%s", buf)
	//Turn the data returned to the writer back into a data structure
	// var us []FitbitOutput
	// //to turn it into a JSON array we need to add commas between events and add
	// // brackets around the whole result
	// b := []byte("[" + strings.Replace(string(buf.Bytes()), "}}", "}},", 2) + "]")
	// err = json.Unmarshal(b, &us)
	// if err != nil {
	// 	t.Logf("Input: %s", b)
	// 	t.Fatal(err)
	// }
	//
	// if len(us) != 1 {
	// 	t.Logf("Failed to retrieve data from fitbit.")
	// 	t.Fail()
	// }
}

// TestGetTimeZone is an integration test with the fitbit server.  This tests the
// FitbitReader.getTimeZone method to ensure that it gets the right time zone
// back for the test user.
// Precondition: Set fitbitAccessToken, fitbitAccessToken, fitbitExpires
// constants in a file that is not in your git repo.
func TestGetTimeZone(t *testing.T) {
	tok := newTokenNoExpiry(fitbitRefreshToken, fitbitAccessToken, testTokenType)

	client, _ := getClient(tok, fitbitClientId, fitbitClientSecret, strategyFitbit)
	reader := &FitbitReader{}
	tz, _ := reader.getTimeZone(client)

	if tz != "-07:00" {
		t.Fail()
		t.Log("Wrong time zone returned.")
	}
}

func TestCreateFitbitAuthCodeURL(t *testing.T) {
	conf := oauth2.Config{ClientID: fitbitClientId, ClientSecret: fitbitClientSecret}
	conf.Endpoint = fitbit.Endpoint
	conf.Scopes = []string{"activity"}
	conf.RedirectURL = "https://localhost:8000/en-US/splunkd/services/fitness_for_splunk/fitbit_callback"
	//print a url to go get an access code
	t.Logf("URL: %v\n", conf.AuthCodeURL("state",
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("expires_in", "31536000")))
}

func disabledTestExchangeFitBitToken(t *testing.T) {
	conf := oauth2.Config{ClientID: fitbitClientId, ClientSecret: fitbitClientSecret}
	conf.Endpoint = fitbit.Endpoint
	conf.Scopes = []string{"activity"}
	conf.RedirectURL = "https://localhost:8000/en-US/splunkd/services/fitness_for_splunk/fitbit_callback"

	tok := getTokenFromAccessCode("9b38bc72beb73167ee2ccecf9af7f5cdd6869bf9", conf)
	tokStr, _ := json.Marshal(tok)
	t.Logf("%s\n", tokStr)
}
