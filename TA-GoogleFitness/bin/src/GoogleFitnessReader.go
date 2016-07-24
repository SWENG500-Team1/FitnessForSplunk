package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"google.golang.org/api/fitness/v1"
)

const oauth_time_format = "2006-01-02 15:04:05.00000000 -0700 MST"

func GetDataSources(client *http.Client) []*fitness.DataSource {

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

func GetDataSet(
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
