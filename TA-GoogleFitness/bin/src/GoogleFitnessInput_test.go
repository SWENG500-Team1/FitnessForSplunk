package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

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
	sessionKey, err := splunk.NewSessionKey(accountName, password, splunk.LocalSplunkMgmntURL)
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

//Write's a value to a checkpoint file an then reads it back to confirm that it
// was written correctly
func TestWriteCheckpoint(t *testing.T) {

	tempDir, err := ioutil.TempDir("", "checkpointDir")
	if err != nil {
		t.Logf("Unable to create temporary checkpointDir: %v\n", err)
	}
	//Delete the checkpoint directory when we are done with the test
	defer os.RemoveAll(tempDir)

	//create a time in the past to check against
	checkpointTime := time.Now().Add(-30 * time.Minute)

	//create an input and configuration for that input that drives the
	// checkpoint dir
	input := &GoogleFitnessInput{}
	config := &splunk.ModInputConfig{}
	config.CheckpointDir = tempDir
	stanza := &splunk.ModInputStanza{}
	stanza.StanzaName = "input://TestStanza"
	config.Stanzas = []splunk.ModInputStanza{*stanza}
	input.ModInputConfig = config

	//write the checkpoint.
	input.writeCheckPoint(checkpointTime)

	//get the checkpoint back
	startTime, _ := input.getTimes()

	//Validate that the time we sent in is the same as the time we get back.
	if startTime != checkpointTime {
		t.Log("Incorrect start time retreived after checkpoint written.")
		t.Fail()
	}
}

func TestScheme(t *testing.T) {
	scheme := `   <scheme>
      <title>Google Fitness</title>
      <description>Retrieves fitness data from Google Fitness.</description>
      <use_external_validation>false</use_external_validation>
      <streaming_mode>simple</streaming_mode>
      <arg name="force_cert_validation">
         <title>ForceCertValidation</title>
         <description>If true the input requires certificate validation when making REST calls to Splunk</description>
      </arg>
   </scheme>`
	buf := new(bytes.Buffer)
	input, err := NewGoogleFitnessInput(strings.NewReader(scheme), bufio.NewWriter(buf))
	if err != nil {
		t.Logf("Unable to create GoogleFitnessInput\n%v\n", err)
		t.Fail()
	}

	//Scheme should write to the buffer
	input.ReturnScheme()

	if scheme != buf.String() {
		t.Logf("Returned scheme does not match expected scheme.\nExpected:%v\nReceived:%v\n", scheme, buf.String())
		t.Fail()
	}
}
