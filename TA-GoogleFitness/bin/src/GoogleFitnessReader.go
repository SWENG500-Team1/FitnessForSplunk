package main

import (
	"bufio"
	"log"
	"net/http"
	"strconv"
	"time"

	"google.golang.org/api/fitness/v1"
)

const oauth_time_format = "2006-01-02 15:04:05.00000000 -0700 MST"

type GoogleFitnessReader struct {
	startTime time.Time
	endTime   time.Time
}

func (input *GoogleFitnessReader) getData(
	client *http.Client,
	writer *bufio.Writer) time.Time {

	lastOutputTime := input.startTime

	dataSources := input.getDataSources(client)
	for _, dataSource := range dataSources {
		dataset := input.getDataSet(client, input.startTime, input.endTime, *dataSource)

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
