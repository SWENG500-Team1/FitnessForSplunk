package main

import (
	"testing"
	"time"

	"github.com/AndyNortrup/GoSplunk"
)

func TestRetrieval(t *testing.T) {

	sessionKey, err := splunk.GetSessionKey("testing_user", "TestAccount")
	if err != nil {
		t.Logf("Unable to get a Session key: %v\n", err)
	}

	//TODO: Replace hard coded values with pull from arguments
	reader := NewFitnessReader(getAppCredentials(sessionKey.SessionKey))

	//TODO: Replace hard coded values with pull from storage/passwords
	/*TODO: Determine if the value from storage/passwords has a refresh token.
	  Yes: Refresh the existing token.
	  No: Get a refresh token and store new token
	*/
	tok := reader.getTokenFromRefreshToken("1/7u5ngLKEF2MiVYHvnWwYKRIb8s3s8u2e8JtHZ2yjUAQ",
		"ya29.Ci8IA_du7mknNus-G_UTfiWB3FHeqdpIqEj_bwaUSvB2lYvsZSuKB7E-2TVuDM44sw",
		"2016-06-21 07:59:23.44961918 -0700 PDT",
		"Bearer")

	startTime := time.Now().Add(-6 * time.Hour)
	endTime := time.Now()
	latestTime := outputData(reader, tok, startTime, endTime)
	if latestTime == startTime {
		t.Log("Unable to pull new data from data srouces.")
		t.Fail()
	}
}
