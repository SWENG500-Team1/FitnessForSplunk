package main

import (
	"bufio"
	"errors"
	"net/http"
	"time"
)

const STRATEGY_GOOGLE string = "GoogleFitness"
const STRATEGY_FITBIT string = "FitBit"
const STRATEGY_MICROSOFT string = "Microsoft"
const STRATEGY_PARAM_NAME string = "FitnessService"

type FitnessReader interface {
	//getData takes a start and end time, and HTTP client for communication with
	//the service an output channle to return data for writing to the command line,
	// and a wait group to make sure things stay open until we are done with all
	// of the go routines.
	// Returns a time of the last data retrived.
	getData(
		client *http.Client,
		output *bufio.Writer,
		username string) time.Time
}

func readerFactory(strategy string, startTime time.Time, endTime time.Time) (FitnessReader, error) {
	switch {
	case strategy == STRATEGY_GOOGLE:
		reader := &GoogleFitnessReader{startTime: startTime, endTime: endTime}
		return reader, nil
	default:
		return nil, errors.New("Unsupported reader requested: " + string(strategy))
	}
}
