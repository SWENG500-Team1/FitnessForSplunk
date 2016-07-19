package main

import (
	"bufio"
	"io/ioutil"
	"testing"
	"time"
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

	input := &GoogleFitnessInput{}

	//TODO: Replace hard coded values with pull from arguments
	reader := NewFitnessReader(testClientId, testClientSecret)

	//TODO: Replace hard coded values with pull from storage/passwords
	/*TODO: Determine if the value from storage/passwords has a refresh token.
	  Yes: Refresh the existing token.
	  No: Get a refresh token and store new token
	*/
	tok := reader.getTokenFromRefreshToken(testRefreshToken,
		testAccessToken,
		testExpires,
		testTokenType)

	devNull := bufio.NewWriter(ioutil.Discard)

	latestTime := input.fetchData(reader, tok, startTime, endTime, devNull)
	if latestTime.Nanosecond() == startTime.Nanosecond() {
		t.Logf("Unable to pull new data from data sources.  "+
			"\tStart Time: %v\tLatest Time: %v\n",
			startTime,
			latestTime)
		t.Fail()
	}
}
