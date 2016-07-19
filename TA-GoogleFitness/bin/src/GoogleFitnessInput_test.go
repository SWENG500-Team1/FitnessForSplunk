package main

import (
	"testing"

	"github.com/AndyNortrup/GoSplunk"
)

const accountName string = "testing_user"
const password string = "TestAccount"

// test_getAppCredentials is an integration test that requires the following:
// 1. Splunk server running locally.
// 2. User on that Splunk server with admin access so that we can access
//    storage/passwords.
// 3. A google clientId and clientSecret password loaded in the
//    APP_NAME local/passwords.conf file.
func TestGetAppCredentials(t *testing.T) {
	sessionKey, err := splunk.NewSessionKey(accountName, password)
	if err != nil {
		t.Fatalf("Unable to get session key: %v\n", err)
	}

	input := &GoogleFitnessInput{}
	clientId, clientSecret := input.getAppCredentials(sessionKey.SessionKey)
	t.Logf("ClientId Expected: %v\tReceived: %v\n", testClientId, clientId)
	if clientId != testClientId {
		t.Fail()
	}

	t.Logf("ClientSecret Expected: %v\tReceived: %v\n", testClientSecret, clientSecret)
	if clientSecret != testClientSecret {
		t.Fail()
	}
}
