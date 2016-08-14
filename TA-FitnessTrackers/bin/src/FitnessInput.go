package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/AndyNortrup/GoSplunk"
)

const APP_NAME string = "TA-FitnessTrackers"
const ENFORCE_CERT_VALIDATION string = "force_cert_validation"

type FitnessInput struct {
	*splunk.ModInputConfig
	reader io.Reader //Location to read configurations from
	writer io.Writer //Location to write configurations to
}

//Write the scheme to input.writer
func (input *FitnessInput) ReturnScheme() {
	arguments := append([]splunk.Argument{}, splunk.Argument{
		Name:        ENFORCE_CERT_VALIDATION,
		Title:       "ForceCertValidation",
		Description: "If true the input requires certificate validation when making REST calls to Splunk",
		DataType:    "boolean",
	},
		splunk.Argument{
			Name:        strategyParamName,
			Title:       "FitnessService",
			Description: "Enter the name of the Fitness Service to be polled.  Options are: 'GoogleFitness', 'FitBit', 'Microsoft'",
			DataType:    "string",
		})

	scheme := &splunk.Scheme{
		Title:                 "Fitness Trackers",
		Description:           "Retrieves fitness data from Google Fitness and fitbit.",
		UseExternalValidation: true,
		StreamingMode:         "simple",
		Args:                  arguments,
	}

	enc := xml.NewEncoder(input.writer)
	enc.Indent("   ", "   ")
	if err := enc.Encode(scheme); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}

func (input *FitnessInput) ValidateScheme() (bool, string) {
	config, err := splunk.ReadModInputConfig(input.reader)
	if err != nil {
		return false, "Unable to parse configuration." + err.Error()
	}

	for _, stanza := range config.Stanzas {
		for _, param := range stanza.Params {
			//Check that the parameter STRAGEGY_PARAM_NAME is one of our defined
			// strategies for getting data
			if param.Name == strategyParamName &&
				!(param.Value == string(strategyGoogle) ||
					param.Value == strategyFitbit ||
					param.Value == strategyMicrosoft) {
				return false, "Improper service '" + param.Value + "' name indicated."
			}
		}
	}
	return true, ""
}

func (input *FitnessInput) StreamEvents() {

	config, err := splunk.ReadModInputConfig(input.reader)
	if err != nil {
		log.Printf("Unable to read Modular Input config from reader.")
	}
	input.ModInputConfig = config

	tokens, err := getUsers(splunk.LocalSplunkMgmntURL, input.SessionKey, input.getStrategy())
	if err != nil {
		log.Printf("Unable to get user tokens: %v", err)
	}

	clientId, clientSecret := input.getAppCredentials()

	for _, token := range tokens {

		//Quick and dirty fix because it's late.
		//KV store escapes the backslash in the refresh tokens from google.
		token.RefreshToken = strings.Replace(token.RefreshToken, "\\", "", -1)
		//Create HTTP client
		client, newToken := getClient(&token.Token,
			clientId,
			clientSecret,
			input.getStrategy())

		//Fitbit is stupid and makes you cache a new refresh token every time.
		//Probably more secure, but much more of a pain to handle.
		if input.getStrategy() == strategyFitbit {
			token.RefreshToken = newToken.RefreshToken
			err := updateKVStoreToken(token, input.getStrategy(), input.SessionKey)
			if err != nil {
				log.Printf("FITBIT: Failed to update KV Store with new key for user: %v Err: %v",
					token.UserID, err)
			}
		}

		//Get start and end points from checkpoint
		startTime, endTime := input.getTimes(input.getStrategy(), token.Name, token.UserID)
		//Create a Fitness Reader to go get the data
		fitnessReader, err := readerFactory(input.getStrategy(), startTime, endTime)
		if err != nil {
			log.Fatal(err)
		}

		input.writeCheckPoint(input.getStrategy(),
			token.Name,
			token.UserID,
			fitnessReader.getData(client, bufio.NewWriter(os.Stdout), token))
	}
}

//get the value of the strategy parameter from the configuration.
func (input *FitnessInput) getStrategy() string {
	var strategy string

	for _, stanza := range input.Stanzas {
		for _, param := range stanza.Params {
			if param.Name == strategyParamName {
				strategy = param.Value
			}
		}
	}
	if strategy == "" {
		log.Fatalf("No strategy passed to Fitness Input")
	}
	return strategy
}

// getAppCredentials makes a call to the storage/passwords enpoint and retrieves
// an appId and clinetSecret for the application.  The appId is stored in the
// password field of the endpoint data and the appId is in the username.
func (input *FitnessInput) getAppCredentials() (string, string) {
	passwords, err := splunk.GetEntities(splunk.LocalSplunkMgmntURL,
		[]string{"storage", "passwords"},
		APP_NAME,
		"nobody",
		input.SessionKey)

	if err != nil || len(passwords.Entries) == 0 {
		log.Fatalf("Unable to retrieve password entries for: %v"+
			"Error: %v\n", err, input.Stanzas[0].StanzaName)
	}

	for _, entry := range passwords.Entries {
		var clientId, clientSecret string
		strategyKey := false
		for _, key := range entry.Contents.Keys {
			if key.Name == "clear_password" {
				clientSecret = key.Value
			}
			if key.Name == "username" {
				clientId = key.Value
			}
			if key.Name == "realm" && key.Value == input.getStrategy() {
				strategyKey = true
			}
		}

		if strategyKey {
			return clientId, clientSecret
		}
	}
	log.Fatalf("No application credentials found for service \"%v\"", input.getStrategy())
	return "", ""
}

//getTimes returns a startTime and an endTime value.  endTime is retrived from
// a checkpoint file, if not it returns the current time.
// The end time is always the current time.
func (input *FitnessInput) getTimes(service, username, userid string) (time.Time, time.Time) {
	startTime, err := input.readCheckPoint(service, username, userid)
	if err != nil {
		startTime = time.Now().AddDate(0, 0, -5)
	}
	endTime := time.Now()
	return startTime, endTime
}

func (input *FitnessInput) writeCheckPoint(service, username, userid string, t time.Time) {

	//Encode the time we've been given into bytes
	g, err := t.GobEncode()
	if err != nil {
		log.Fatalf("Unable to encode checkpoint time: %v\n", err)
	}

	//Write the checkpoint
	err = ioutil.WriteFile(input.getCheckPointPath(service, username, userid), g, 0644)
	if err != nil {
		log.Fatalf("Error writing checkpoint file: %v\n", err)
	}
	log.Printf("Wrote checkpoint for %v - %v: %v", service, username, t)

}

func (input *FitnessInput) readCheckPoint(service, username, userid string) (time.Time, error) {
	b, err := ioutil.ReadFile(input.getCheckPointPath(service, username, userid))
	if err != nil {
		log.Printf("Unable to read checkpoint file:%v\n", err)
		return time.Now(), err
	}
	var t time.Time
	err = t.GobDecode(b)
	if err != nil {
		log.Printf("Unable to decode checkpoint file: %v\n", err)
		return time.Now().AddDate(0, 0, -5), err
	}

	log.Printf("Read checkpoint for %v - %v: %v", service, username, t)

	return t, nil
}

// Takes the checkpoint dir from and config stanza name from the input and
// creates a checkpoint dir.  Should be unique for each input
func (input *FitnessInput) getCheckPointPath(service, username, userid string) string {
	//Create a hash of the stanza name as a filename
	fileName := service + "_" + username + "_" + userid
	path := path.Join(input.CheckpointDir, fileName)
	return path
}
