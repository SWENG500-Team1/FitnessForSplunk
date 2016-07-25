package main

import (
	"bufio"
	"net/http"
	"time"
)

type FitnessReader interface {
	//getData takes a start and end time, and HTTP client for communication with
	//the service an output channle to return data for writing to the command line,
	// and a wait group to make sure things stay open until we are done with all
	// of the go routines.
	// Returns a time of the last data retrived.
	getData(
		client *http.Client,
		output *bufio.Writer) time.Time
}
