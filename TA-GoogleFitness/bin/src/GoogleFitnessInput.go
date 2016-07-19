package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/AndyNortrup/GoSplunk"

	"golang.org/x/oauth2"
)

const APP_NAME string = "TA-GoogleFitness"

type GoogleFitnessInput struct{}

func (input *GoogleFitnessInput) ReturnScheme() {
	arguments := append([]splunk.Argument{}, splunk.Argument{
		Name:        "client_id",
		Title:       "Googl API Client ID",
		Description: "This ID identifies the application to Google.",
	})

	scheme := &splunk.Scheme{
		Title:                 "Google Fitness",
		Description:           "Retrieves fitness data from Google Fitness.",
		UseExternalValidation: false,
		StreamingMode:         "simple",
		Args:                  arguments,
	}

	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("   ", "   ")
	if err := enc.Encode(scheme); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}

func (input *GoogleFitnessInput) ValidateScheme() {
	fmt.Printf("Validate Scheme\n")
}

func (input *GoogleFitnessInput) StreamEvents() {

	//Get parameters from std in.
	config, err := splunk.ReadModInputConfig(bufio.NewReader(os.Stdin))
	if err != nil {
		log.Fatalf("Unable to read configuration from Stdin: %v\n", err)
	}

	//Create FitnessReader
	reader := NewFitnessReader(input.getAppCredentials(config.SessionKey))

	//TODO: Replace hard coded values with pull from storage/passwords
	/*TODO: Determine if the value from storage/passwords has a refresh token.
	Yes: Refresh the existing token.
	No: Get a refresh token and store new token
	*/
	tok := reader.getTokenFromRefreshToken("1/7u5ngLKEF2MiVYHvnWwYKRIb8s3s8u2e8JtHZ2yjUAQ",
		"ya29.Ci8IA_du7mknNus-G_UTfiWB3FHeqdpIqEj_bwaUSvB2lYvsZSuKB7E-2TVuDM44sw",
		"2016-06-21 07:59:23.44961918 -0700 PDT",
		"Bearer")

	startTime, endTime := input.getTimes()
	writer := bufio.NewWriter(os.Stdout)
	input.writeCheckPoint(input.fetchData(reader, tok, startTime, endTime, writer))
}

// getAppCredentials makes a call to the storage/passwords enpoint and retrieves
// an appId and clinetSecret for the application.  The appId is stored in the
// password field of the endpoint data and the appId is in the username.
func (input *GoogleFitnessInput) getAppCredentials(sessionKey string) (string, string) {
	passwords, err := splunk.GetEntities([]string{"storage", "passwords"},
		APP_NAME,
		"nobody",
		sessionKey)

	if err != nil {
		log.Fatalf("Unable to retrieve password entries for TA-GoogleFitness: %v\n",
			err)
	}

	var clientSecret string
	var clientId string

	for _, entry := range passwords.Entries {
		//Because there could/should be multiple stored passwords we need to check
		// the id for `apps.googleusercontent.com` because the id is based on the
		// username.

		if strings.Contains(entry.ID, "apps.googleusercontent.com") {
			for _, key := range entry.Contents.Keys {
				if key.Name == "clear_password" {
					clientSecret = key.Value
				}
				if key.Name == "username" {
					clientId = key.Value
				}
			}
		}
	}
	return clientId, clientSecret
}

//getTimes returns a startTime and an endTime value.  endTime is retrived from
// a checkpoint file, if not it returns the current time.
// The end time is always the current time.
func (input *GoogleFitnessInput) getTimes() (time.Time, time.Time) {
	startTime, err := input.readCheckPoint()
	if err != nil {
		startTime = time.Now()
	}
	endTime := time.Now()
	return startTime, endTime
}

func (input *GoogleFitnessInput) fetchData(
	reader *FitnessReader,
	tok *oauth2.Token,
	startTime time.Time,
	endTime time.Time,
	writer *bufio.Writer) time.Time {

	defer writer.Flush()

	lastOutputTime := startTime

	dataSources := reader.GetDataSources(tok)
	for _, dataSource := range dataSources {
		dataset := reader.GetDataSet(tok, startTime, endTime, *dataSource)

		for _, point := range dataset.Point {
			json, _ := point.MarshalJSON()
			writer.Write(json)
			writer.Flush()

			//find the last time recorded so that we can write that as the checkpoint
			if time.Unix(0, point.EndTimeNanos).After(lastOutputTime) {
				lastOutputTime = time.Unix(0, point.EndTimeNanos)
			}
		}
	}
	return lastOutputTime
}

func (input *GoogleFitnessInput) writeCheckPoint(t time.Time) {

	//TODO: Replace the filename: checkpoint.txt with something that will work with
	// multiple names of the input

	//Encode the time we've been given into bytes
	g, err := t.GobEncode()
	if err != nil {
		log.Fatalf("Unable to encode checkpoint time: %v\n", err)
	}

	//Write the checkpoint
	err = ioutil.WriteFile(input.getCheckPointPath(), g, 0644)
	if err != nil {
		log.Fatalf("Error writing checkpoint file: %v\n", err)
	}
}

func (input *GoogleFitnessInput) readCheckPoint() (time.Time, error) {
	b, err := ioutil.ReadFile(input.getCheckPointPath())
	if err != nil {
		log.Printf("Unable to read checkpoint file:%v\n", err)
		return time.Now(), err
	}
	var t time.Time
	err = t.GobDecode(b)
	if err != nil {
		log.Printf("Unable to decode checkpoint file: %v\n", err)
		return time.Now().Add(-2 * time.Hour), err
	}
	return t, nil
}

func (input *GoogleFitnessInput) getCheckPointPath() string {
	//Get the app directory
	base := os.Getenv("SPLUNK_HOME")
	if base == "" {
		log.Fatal("Unable to find $SPLUNK_HOME")
	}
	path := path.Join(base, "etc/apps/TA-GoogleFitness/checkpoint.txt")
	return path
}
