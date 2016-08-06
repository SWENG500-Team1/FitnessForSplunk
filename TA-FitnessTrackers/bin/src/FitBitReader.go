package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type FitbitReader struct {
	startTime, endTime time.Time
}

func NewFitbitReader(startTime, endTime time.Time) (*FitbitReader, error) {
	if startTime.After(endTime) {
		return &FitbitReader{}, errors.New("Start time after end time not allowed")
	}

	//Establish what the last date we will query is.  We only want to get a
	// complete day so if input.endTime is today's date, subtract one day.
	lastDate := time.Date(endTime.Year(),
		endTime.Month(),
		endTime.Day(),
		0, 0, 0, 0, time.Local)

	//If the last date is today then back up a date
	if lastDate.Year() == time.Now().Year() &&
		lastDate.Month() == time.Now().Month() &&
		lastDate.Day() == time.Now().Day() {

		lastDate = lastDate.AddDate(0, 0, -1)
	}
	//we can't request dates in the future.  If the last date is after the current
	// time back it up to yesterday.
	if lastDate.After(time.Now()) {
		lastDate = time.Now().AddDate(0, 0, -1)
	}

	return &FitbitReader{startTime: startTime, endTime: lastDate}, nil
}

func (input *FitbitReader) getData(
	client *http.Client,
	writer *bufio.Writer,
	user User) time.Time {

	var dateFormat string = "2006-01-02"
	var lastDate time.Time
	//Make a request for every complete day between startDate and endDate
	for rd := input.startTime; rd.Before(input.endTime) || rd.Equal(input.endTime); rd = rd.AddDate(0, 0, 1) {
		//Make a request to the API endpoint:
		// https://api.fitbit.com/1/user/[user-id]/activities/date/[date].json
		requestString := "https://api.fitbit.com/1/user/-/activities/date/" +
			rd.Format(dateFormat) + ".json"

		response, err := client.Get(requestString)
		if err != nil {
			log.Fatalf("Error making request to fitbit: %v", err)
		}

		if response.StatusCode != http.StatusOK {
			log.Printf("Non-200 status code from fitbit request: %v, %v",
				response.StatusCode,
				requestString)
		}

		input.decodeAndPrint(response.Body, writer, user.name, rd)
		lastDate = rd
	}
	return lastDate
}

func (input *FitbitReader) decodeAndPrint(reader io.Reader,
	writer *bufio.Writer,
	username string, date time.Time) {

	decoder := json.NewDecoder(reader)
	summary := &FitbitActivitySummary{}
	err := decoder.Decode(summary)
	if err != nil {
		log.Printf("Error decoding fitbit summary: %v\n", err)
		return
	}

	us := NewFitbitOutput(username, date, summary)

	b, _ := json.Marshal(us)
	writer.WriteString(fmt.Sprintf("%s\n", b))
	writer.Flush()
}

//userStruct: A a struct to glue username and Activity Summaries together
type FitbitOutput struct {
	Source string    `json:"Source"`
	User   string    `json:"User"`
	Date   time.Time `json:"Date"`
	*FitbitActivitySummary
}

func NewFitbitOutput(user string, date time.Time, summary *FitbitActivitySummary) *FitbitOutput {
	return &FitbitOutput{
		Source: strategyFitbit,
		User:   user,
		Date:   date,
		FitbitActivitySummary: summary,
	}
}

type FitbitActivitySummary struct {
	Activities []FitbitActivity `json:"activities"`
	Goals      FitbitGoal       `json:"goals"`
	Summary    FitbitSummary    `json:"summary"`
}

type FitbitActivity struct {
	ActivityId       float32 `json:"activityId"`
	ActivityParentId float32 `json:"activityParentId"`
	Calories         float32 `json:"calories"`
	Description      string  `json:"description"`
	Distance         float32 `json:"distance"`
	Duration         float32 `json:"duration"`
	HasStartTime     bool    `json:"hasStartTime"`
	IsFavorite       bool    `json:"isFavorite"`
	LogId            float32 `json:"logId"`
	Name             string  `json:"name"`
	StartTime        string  `json:"startTime"`
	Steps            float32 `json:"steps"`
}

type FitbitGoal struct {
	CaloriesOut float32 `json:"caloriesOut"`
	Distance    float32 `json:"distance"`
	Floors      float32 `json:"floors"`
	Steps       float32 `json:"steps"`
}

type FitbitSummary struct {
	ActivityCalories     float32          `json:"activityCalories"`
	CaloriesBMR          float32          `json:"caloriesBMR"`
	CaloriesOut          float32          `json:"caloriesOut"`
	Distances            []FitbitDistance `json:"distances"`
	Elevation            float32          `json:"elevation"`
	FairlyActiveMinutes  float32          `json:"fairlyActiveMinutes"`
	Floors               float32          `json:"floors"`
	LightlyActiveMinutes float32          `json:"lightlyActiveMinutes"`
	MarginalCalories     float32          `json:"marginalCalories"`
	SedentaryMinutes     float32          `json:"sedentaryMinutes"`
	Steps                float32          `json:"steps"`
	VeryActiveMinutes    float32          `json:"veryActiveMinutes"`
}

type FitbitDistance struct {
	Activity string  `json:"activity"`
	Distance float32 `json:"distance"`
}
