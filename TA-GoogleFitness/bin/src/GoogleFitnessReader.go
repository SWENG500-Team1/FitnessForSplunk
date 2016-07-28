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

const oauth_time_format = "2006-01-02 15:04:05.00000000 -0700 MST"
const sessionTimeFormat string = "2006-01-02T15:04:05.00Z"

type GoogleFitnessReader struct {
	startTime time.Time
	endTime   time.Time
	username  string
}

func (input *GoogleFitnessReader) getData(
	client *http.Client,
	writer *bufio.Writer,
	username string) time.Time {

	lastOutputTime := input.startTime

	dataSources := input.getDataSources(client)
	for _, dataSource := range dataSources {
		dataset := input.getDataSet(client, input.startTime, input.endTime, *dataSource)

		for _, point := range dataset.Point {
			type event struct {
				Username  string            `json:"username"`
				DataPoint fitness.DataPoint `json:"DataPoint"`
			}
			e := &event{Username: username, DataPoint: *point}
			encoder := json.NewEncoder(writer)
			encoder.Encode(e)
			writer.Flush()

			//find the last time recorded so that we can write that as the checkpoint
			if time.Unix(0, point.EndTimeNanos).After(lastOutputTime) {
				lastOutputTime = time.Unix(0, point.EndTimeNanos)
			}
		}
	}

	input.getSessions(client, input.startTime, lastOutputTime, writer)
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
		log.Fatalf("Error getting DataSources: %v\n", err)
	}

	return response.DataSource
}

func (input *GoogleFitnessReader) getDataSet(
	client *http.Client,
	startTime time.Time,
	endTime time.Time,
	dataSource fitness.DataSource) *fitness.Dataset {

	dataSetId := strconv.FormatInt(startTime.UnixNano(), 10) + "-" +
		strconv.FormatInt(endTime.UnixNano(), 10)

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
	startTime, endTime time.Time, writer *bufio.Writer) error {

	s, err := fitness.New(client)
	if err != nil {
		log.Fatalf("Unable to create Google Fitness Session: %v\n", err)
	}
	sessionService := fitness.NewUsersSessionsService(s)
	sessionCall := sessionService.List("me")
	sessionCall.StartTime(startTime.Format(sessionTimeFormat))
	sessionCall.EndTime(endTime.Format(sessionTimeFormat))

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
