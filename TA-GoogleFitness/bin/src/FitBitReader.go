package main

import (
	"bufio"
	"net/http"
	"time"
)

type FitBitReader struct {
	startTime time.Time
	endTime   time.Time
}

func getData(
	client *http.Client,
	output *bufio.Writer,
	username string) time.Time {

	return time.Now()
}
