package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"google.golang.org/api/fitness/v1"
)

const googleOauthTimeFormat = "2006-01-02 15:04:05.00000000 -0700 MST"
const fitbitOauthTimeFormat = "2006-01-02T15:04:05.000000000-07:00"
const sessionTimeFormat string = "2006-01-02T15:04:05.00Z"

type GoogleFitnessReader struct {
	startTime time.Time
	endTime   time.Time
	username  string
}

func (input *GoogleFitnessReader) getData(
	client *http.Client,
	writer *bufio.Writer,
	user User) time.Time {

	lastOutputTime := input.startTime

	dataSources := input.getDataSources(client)
	for _, dataSource := range dataSources {
		dataset := input.getDataSet(client, *dataSource)

		for _, point := range dataset.Point {
			type event struct {
				Username  string            `json:"username"`
				DataPoint fitness.DataPoint `json:"DataPoint"`
			}
			e := &event{Username: user.Name, DataPoint: *point}
			encoder := json.NewEncoder(writer)
			encoder.Encode(e)
			writer.Flush()

			//find the last time recorded so that we can write that as the checkpoint
			if time.Unix(0, point.EndTimeNanos).After(lastOutputTime) {
				input.endTime = time.Unix(0, point.EndTimeNanos)
			}
		}
	}

	input.getSessions(client, writer)
	return lastOutputTime
}

func (input *GoogleFitnessReader) getDataSources(client *http.Client) []*fitness.DataSource {

	service, err := fitness.New(client)
	if err != nil {
		log.Fatalf("Unable to create DataSource service: %v\n", err)
	}
	dataSourceService := fitness.NewUsersDataSourcesService(service)
	call := dataSourceService.List("me")
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error getting DataSources for User: %v\t Error %v\n",
			input.username, err)
	}

	return response.DataSource
}

func (input *GoogleFitnessReader) getDataSet(
	client *http.Client,
	dataSource fitness.DataSource) *fitness.Dataset {

	dataSetId := strconv.FormatInt(input.startTime.UnixNano(), 10) + "-" +
		strconv.FormatInt(input.endTime.UnixNano(), 10)

	service, err := fitness.New(client)
	if err != nil {
		log.Fatalf("Error creating DataSet Service: %v\n", err)
	}

	dataSetService := fitness.NewUsersDataSourcesDatasetsService(service)
	request := dataSetService.Get("me", dataSource.DataStreamId, dataSetId)
	resp, err := request.Do()
	if err != nil {
		log.Fatalf("Error Getting DataSet: %v\n", err)
	}

	return resp
}

func (input *GoogleFitnessReader) getSessions(client *http.Client,
	writer *bufio.Writer) error {

	s, err := fitness.New(client)
	if err != nil {
		log.Fatalf("Unable to create Google Fitness Session: %v\n", err)
	}
	sessionService := fitness.NewUsersSessionsService(s)
	sessionCall := sessionService.List("me")
	sessionCall.StartTime(input.startTime.Format(sessionTimeFormat))
	sessionCall.EndTime(input.endTime.Format(sessionTimeFormat))

	list, err := sessionCall.Do()
	if err != nil {
		return err
	}

	type SessionData struct {
		Username string          `json:"username"`
		Session  fitness.Session `json:"session"`
	}
	encoder := json.NewEncoder(writer)

	for _, session := range list.Session {
		event := &SessionData{Username: input.username, Session: *session}

		encoder.Encode(event)
		writer.Flush()
	}
	return nil
}
