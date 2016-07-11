package main

import (
	"fmt"
	"testing"
	"time"
)

func TestRetrieval(t *testing.T) {

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

	startTime := time.Now().Add(-6 * time.Hour)
	endTime := time.Now()
	latestTime := outputData(reader, tok, startTime, endTime)
	fmt.Printf("Last recorded time: %v\n", latestTime)
}
