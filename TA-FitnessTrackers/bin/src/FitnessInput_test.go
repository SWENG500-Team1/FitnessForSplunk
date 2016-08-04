package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
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
	accessKey, err := splunk.NewSessionKey(accountName, password, splunk.LocalSplunkMgmntURL)
	if err != nil {
		t.Logf("Unable to get session key: %v\n", err)
	}

	config := &splunk.ModInputConfig{}
	config.SessionKey = accessKey.SessionKey
	stanza := splunk.ModInputStanza{}
	stanza.Params = append([]splunk.ModInputParam{},
		splunk.ModInputParam{Name: strategyParamName,
			Value: strategyGoogle})
	config.Stanzas = append(config.Stanzas, stanza)

	input := FitnessInput{ModInputConfig: config}
	clientId, clientSecret := input.getAppCredentials()
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
	input := &FitnessInput{}
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
      <use_external_validation>true</use_external_validation>
      <streaming_mode>simple</streaming_mode>
      <arg name="force_cert_validation">
         <title>ForceCertValidation</title>
         <description>If true the input requires certificate validation when making REST calls to Splunk</description>
         <data_type>boolean</data_type>
      </arg>
      <arg name="` + strategyParamName + `">
         <title>FitnessService</title>
         <description>Enter the name of the Fitness Service to be polled.  Options are: &#39;GoogleFitness&#39;, &#39;FitBit&#39;, &#39;Microsoft&#39;</description>
         <data_type>string</data_type>
      </arg>
   </scheme>`
	buf := new(bytes.Buffer)
	input := &FitnessInput{reader: strings.NewReader(scheme), writer: bufio.NewWriter(buf)}

	//Scheme should write to the buffer
	input.ReturnScheme()

	if scheme != buf.String() {
		t.Logf("Returned scheme does not match expected scheme.\nExpected:%v\nReceived:%v\n", scheme, buf.String())
		t.Fail()
	}
}

func TestSchemeValidation(t *testing.T) {
	improperValue := "Not a strategy."
	correctSchemes := []string{`<input>
			<server_host>myHost</server_host>
			<server_uri>https://127.0.0.1:8089</server_uri>
			<session_key>123102983109283019283</session_key>
			<checkpoint_dir>/opt/splunk/var/lib/splunk/modinputs</checkpoint_dir>
			<configuration>
				<stanza name="TA-GoogleFitness://test1">
						<param name="` + strategyParamName + `">` + strategyFitbit + `</param>
						<param name="other_param">other_value</param>
				</stanza>
			</configuration>
		</input>`,
		`<input>
				<server_host>myHost</server_host>
				<server_uri>https://127.0.0.1:8089</server_uri>
				<session_key>123102983109283019283</session_key>
				<checkpoint_dir>/opt/splunk/var/lib/splunk/modinputs</checkpoint_dir>
				<configuration>
					<stanza name="TA-GoogleFitness://test1">
							<param name="` + strategyParamName + `">` + strategyGoogle + `</param>
							<param name="other_param">other_value</param>
					</stanza>
				</configuration>
			</input>`,
		`<input>
					<server_host>myHost</server_host>
					<server_uri>https://127.0.0.1:8089</server_uri>
					<session_key>123102983109283019283</session_key>
					<checkpoint_dir>/opt/splunk/var/lib/splunk/modinputs</checkpoint_dir>
					<configuration>
						<stanza name="TA-GoogleFitness://test1">
								<param name="` + strategyParamName + `">` + strategyMicrosoft + `</param>
								<param name="other_param">other_value</param>
						</stanza>
					</configuration>
				</input>`}

	badScheme := `<input>
						<server_host>myHost</server_host>
						<server_uri>https://127.0.0.1:8089</server_uri>
						<session_key>123102983109283019283</session_key>
						<checkpoint_dir>/opt/splunk/var/lib/splunk/modinputs</checkpoint_dir>
						<configuration>
							<stanza name="TA-GoogleFitness://test1">
									<param name="` + strategyParamName + `">` + improperValue + `</param>
							</stanza>
						</configuration>
					</input>`

	for _, scheme := range correctSchemes {
		reader := strings.NewReader(scheme)
		input := &FitnessInput{reader: reader}
		result, msg := input.ValidateScheme()
		if result != true {
			t.Logf("Failed to validate scheme: %v\n%v\n", msg, scheme)
			t.Fail()
		}
	}

	reader := strings.NewReader(badScheme)
	input := &FitnessInput{reader: reader}
	result, _ := input.ValidateScheme()
	if result == true {
		t.Logf("Invalid scheme passed validation: \n%v\n", badScheme)
		t.Fail()
	}

}

//TestGetReader builds a ModInputConfig setting the STRATEGY_PARAM_NAME value
// then uses reflection to validate that the correct type is being generated.
func TestGetReader(t *testing.T) {
	reader, _ := readerFactory(strategyGoogle, time.Now(), time.Now())
	if reflect.TypeOf(reader) != reflect.TypeOf(&GoogleFitnessReader{}) {
		t.Log("Failed to return GoogleFitnessReader")
		t.Fail()
	}
}

func TestGetReaderFromXML(t *testing.T) {
	config := `<input>
			<server_host>myHost</server_host>
			<server_uri>https://127.0.0.1:8089</server_uri>
			<session_key>123102983109283019283</session_key>
			<checkpoint_dir>/opt/splunk/var/lib/splunk/modinputs</checkpoint_dir>
			<configuration>
				<stanza name="TA-GoogleFitness://test1">
						<param name="` + strategyParamName + `">` + strategyGoogle + `</param>
						<param name="other_param">other_value</param>
				</stanza>
			</configuration>
		</input>`

	parsed, _ := splunk.ReadModInputConfig(strings.NewReader(config))
	input := &FitnessInput{ModInputConfig: parsed}
	reader, err := readerFactory(input.getStrategy(), time.Now(), time.Now())
	if err != nil {
		t.Logf("Error getting FitnessReader: %v", err)
		t.Fail()
	}
	if reflect.TypeOf(reader) != reflect.TypeOf(&GoogleFitnessReader{}) {
		t.Log("Failed to return GoogleFitnessReader")
		t.Fail()
	}
}

func TestGetCredentials(t *testing.T) {
	accessKey, err := splunk.NewSessionKey(accountName, password, splunk.LocalSplunkMgmntURL)
	if err != nil {
		t.Logf("Unable to get session key: %v\n", err)
	}

	credentials, _ := getUsers(splunk.LocalSplunkMgmntURL, accessKey.SessionKey, strategyGoogle)
	if len(credentials) == 0 {
		t.Logf("No credentials recieved from Splunk for: %v", strategyGoogle)
		t.Fail()
	}

	t.Logf("Access Token: %v\nRefreshToken: %v\nType:%v\nExpires:%v",
		credentials[0].AccessToken,
		credentials[0].RefreshToken,
		credentials[0].TokenType,
		credentials[0].Expiry)
}
