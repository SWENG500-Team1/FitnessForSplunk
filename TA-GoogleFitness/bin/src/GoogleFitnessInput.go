package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"golang.org/x/oauth2"
)

type GoogleFitnessInput struct{}

func (input *GoogleFitnessInput) ReturnScheme() {
	arguments := append([]Argument{}, Argument{
		Name:        "client_id",
		Title:       "Googl API Client ID",
		Description: "This ID identifies the application to Google.",
	})

	scheme := &Scheme{
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

	//TODO: Replace hard coded values with pull from arguments
	reader := NewFitnessReader(
		"616872666934-ctkc2btlhme0or0vmar8mlaidt2g1j16.apps.googleusercontent.com",
		"-CTNssDbQMnU5G6UVjkKcioA")

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
	input.writeCheckPoint(outputData(reader, tok, startTime, endTime))
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

func outputData(
	reader *FitnessReader,
	tok *oauth2.Token,
	startTime time.Time,
	endTime time.Time) time.Time {

	lastOutputTime := startTime

	dataSources := reader.GetDataSources(tok)
	for _, dataSource := range dataSources {
		dataset := reader.GetDataSet(tok, startTime, endTime, *dataSource)

		for _, point := range dataset.Point {
			json, _ := point.MarshalJSON()
			fmt.Println(string(json))

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
	err = ioutil.WriteFile(getCheckPointPath(), g, 0644)
	if err != nil {
		log.Fatalf("Error writing checkpoint file: %v\n", err)
	}
}

func (input *GoogleFitnessInput) readCheckPoint() (time.Time, error) {
	b, err := ioutil.ReadFile(getCheckPointPath())
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

func getCheckPointPath() string {
	//Get the app directory
	base := os.Getenv("SPLUNK_HOME")
	if base == "" {
		log.Fatal("Unable to find $SPLUNK_HOME")
	}
	path := path.Join(base, "etc/apps/TA-GoogleFitness/checkpoint.txt")
	return path
}
