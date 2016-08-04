package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/fitbit"
)

/* TestFitbitGetData is an integration test with the Fitbit servers.  It does
require some setup to be run correctly.

1. Set the date range with complete days. (i.e. not today)
2. Set fitbitAccessToken, fitbitAccessToken, fitbitExpires constants in a file
   that is not in your git repo.
3. Make sure you know the exact value of steps the first date and set the value
  of expectedStesp to that value.
*/
func TestFitbitGetData(t *testing.T) {
	const expectedSteps int = 14102
	endTime := time.Date(2016, time.August, 3, 0, 0, 0, 0, time.Local)
	startTime := time.Date(2016, time.August, 1, 0, 0, 0, 0, time.Local)

	reader, err := NewFitbitReader(startTime, endTime)

	tok := newToken(fitbitRefreshToken,
		fitbitAccessToken,
		fitbitExpires,
		testTokenType,
		getTokenTimeFormat(strategyFitbit))

	client := getClient(tok, fitbitClientId, fitbitClientSecret, strategyFitbit)
	buf := bytes.NewBuffer([]byte{})
	writer := bufio.NewWriter(buf)

	//Go get the data from fitbit
	date := reader.getData(client, writer, User{name: "Andy"})

	if date.Day() != 3 {
		t.Logf("Wrong date returned.\nExpected:%v\nRecieved:%v", endTime, date)
		t.Fail()
	}

	//Turn the data returned to the writer back into a data structure
	var us []UserSummary
	//to turn it into a JSON array we need to add commas between events and add
	// brackets around the whole result
	b := []byte("[" + strings.Replace(string(buf.Bytes()), "}}", "}},", 2) + "]")
	err = json.Unmarshal(b, &us)
	if err != nil {
		t.Logf("Input: %s", b)
		t.Fatal(err)
	}

	if len(us) != 3 {
		t.Logf("Failed to retrieve data from fitbit.")
		t.Fail()
	} else {
		expected := expectedSteps
		if us[0].Summary.Steps != 14102 {
			t.Logf("Incorrect step count retrived for 1AUG16. Expected: %v\tRecieved:%v\n",
				expected, us[0].Summary.Steps)
			t.Fail()
		}
	}

}

func TestFitbitDates(t *testing.T) {
	now := time.Now()
	_, err := NewFitbitReader(now, now.AddDate(0, 0, -1))
	if err == nil {
		t.Log("Failed to detect invalid date range.")
		t.Fail()
	}

	fbr, err := NewFitbitReader(now.AddDate(0, 0, -2), now)
	if err != nil {
		t.Log("Failed to create Fitbit reader with valid date range.")
		t.Fail()
	}

	if fbr.endTime.Day() != now.Day()-1 {
		t.Logf("Failed to back date off of current day.")
		t.Fail()
	}

	endTime := time.Date(2016, time.August, 3, 0, 0, 0, 0, time.Local)
	startTime := time.Date(2016, time.August, 1, 0, 0, 0, 0, time.Local)
	fbr, err = NewFitbitReader(startTime, endTime)
	if err != nil {
		t.Log("Failed to create Fitbit reader with valid date range.")
		t.Fail()
	}

	if fbr.endTime.Day() != endTime.Day() {
		t.Log("Date changed improperly when creating Fitbit reader.")
		t.Fail()
	}
}

func TestCreateFitbitAuthCodeURL(t *testing.T) {
	conf := oauth2.Config{ClientID: fitbitClientId, ClientSecret: fitbitClientSecret}
	conf.Endpoint = fitbit.Endpoint
	conf.Scopes = []string{"activity"}
	conf.RedirectURL = "https://www.fitnessforsplunk.ninja:8000/en-US/splunkd/services/fitness_for_splunk/fitbit_callback"
	//print a url to go get an access code
	t.Logf("URL: %v\n", conf.AuthCodeURL("state", oauth2.AccessTypeOnline))
}

func TestExchangeFitBitToken(t *testing.T) {
	conf := oauth2.Config{ClientID: fitbitClientId, ClientSecret: fitbitClientSecret}
	conf.Endpoint = fitbit.Endpoint
	conf.Scopes = []string{"activity"}
	conf.RedirectURL = "https://www.fitnessforsplunk.ninja:8000/en-US/splunkd/services/fitness_for_splunk/fitbit_callback"

	tok := getTokenFromAccessCode("ee0cb70f9effe948e1719cd8c8b700470c27f904", conf)
	tokStr, _ := json.Marshal(tok)
	t.Logf("%s\n", tokStr)
}
