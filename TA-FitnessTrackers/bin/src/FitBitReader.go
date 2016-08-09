package main

import (
	"bufio"
	"encoding/json"
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
	return &FitbitReader{startTime: startTime, endTime: endTime}, nil
}

const dateFormat string = "2006-01-02"
const timeFormat string = "15:04"

func (input *FitbitReader) getData(
	client *http.Client,
	writer *bufio.Writer,
	user User) time.Time {

	//Get the user's time-zone
	tz := input.getTimeZone(client)

	for loopStart := input.startTime; loopStart.Before(input.endTime); loopStart = loopStart.Add(24 * time.Hour) {
		var loopEnd time.Time
		if input.endTime.Sub(loopStart).Hours() < 24 {
			loopEnd = input.endTime
		} else {
			loopEnd = loopStart.Add(23 * time.Hour).Add(59 * time.Minute)
		}
		//Make a request to the API endpoint:
		// https://api.fitbit.com/1/user/[user-id]/activities/date/[date].json
		requestString := "https://api.fitbit.com/1/user/-/activities/steps/date/" +
			loopStart.Format(dateFormat) + "/" +
			loopEnd.Format(dateFormat) + "/" +
			"1min/time/" + loopStart.Format(timeFormat) +
			"/" + loopEnd.Format(timeFormat) + ".json"

		response, err := client.Get(requestString)
		if err != nil {
			log.Fatalf("Error making request to fitbit: %v", err)
		}

		if response.StatusCode != http.StatusOK {
			log.Printf("Non-200 status code from fitbit request: %v, %v",
				response.StatusCode,
				requestString)
		}
		defer response.Body.Close()
		input.decodeAndPrint(response.Body, writer, user.Name, tz)
	}

	// TODO: Replace this with a value from the data structure
	return time.Now()
}

func (input *FitbitReader) decodeAndPrint(reader io.Reader,
	writer *bufio.Writer,
	username, timeZone string) {

	decoder := json.NewDecoder(reader)
	summary := &FitbitIntrdayActivityResponse{}
	err := decoder.Decode(summary)
	if err != nil {
		log.Printf("Error decoding fitbit summary: %v\n", err)
		return
	}

	for _, dataPoint := range summary.Data.Dataset {
		if dataPoint.Value > 0 {

			output := &FitbitOutput{
				Source:   strategyFitbit,
				User:     username,
				DateTime: summary.Summary[0].Date + " " + dataPoint.Time + " " + timeZone,
				Value:    dataPoint.Value,
			}
			b, _ := json.Marshal(output)
			writer.WriteString(fmt.Sprintf("%s\n", b))
			writer.Flush()
		}
	}
}

//getTimeZone makes a call to the fitbit profile endpoint so that we can get
// the user's timezone for proper time series indexing
func (input *FitbitReader) getTimeZone(client *http.Client) string {
	resp, err := client.Get("https://api.fitbit.com/1/user/-/profile.json")
	if err != nil {
		log.Fatalf("Failed to get user profile information: %v\n", err)
	}
	defer resp.Body.Close()

	type profileData struct {
		Timezone string `json:"timezone"`
	}

	type profile struct {
		Data profileData `json:"user"`
	}

	f := &profile{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(f)
	if err != nil {
		log.Fatalf("Failed to decode user profile data: %v", err)
	}

	local, err := time.LoadLocation(f.Data.Timezone)
	if err != nil {
		log.Fatalf("Failed to convert timezone to local: %v\n", err)
	}
	t := time.Date(0, 0, 0, 0, 0, 0, 0, local)
	return string(t.Format("-0700"))
}

//userStruct: A a struct to glue username and Activity Summaries together
type FitbitOutput struct {
	Source   string `json:"Source"`
	User     string `json:"User"`
	DateTime string `json:"Date"`
	Value    int    `json:"Value"`
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

/*{
    "activities-log-steps":[
        {"dateTime":"2014-09-05","value":1433}
    ],
    "activities-log-steps-intraday":{
        "datasetInterval":1,
        "dataset":[
            {"time":"00:00:00","value":0},
            {"time":"00:01:00","value":0},
            {"time":"00:02:00","value":0},
            {"time":"00:03:00","value":0},
            {"time":"00:04:00","value":0},
            {"time":"00:05:00","value":287},
            {"time":"00:06:00","value":287},
            {"time":"00:07:00","value":287},
            {"time":"00:08:00","value":287},
            {"time":"00:09:00","value":287},
            {"time":"00:10:00","value":0},
            {"time":"00:11:00","value":0},

        ]
    }
}*/

type FitbitIntrdayActivityResponse struct {
	Summary FitbitIntradayActivitesSteps   `json:"activities-steps"`
	Data    FitbitIntradayActivityLogSteps `json:"activities-steps-intraday"`
}

type FitbitIntradayActivitesSteps []struct {
	Date  string `json:"dateTime"`
	Value string `json:"value"`
}

type FitbitIntradayActivityLogSteps struct {
	DataSetInterval int                     `json:"datasetInterval"`
	DatSetType      string                  `json:"dataSetType"`
	Dataset         []FitbitIntradayDataset `json:"dataset"`
}

type FitbitIntradayDataset struct {
	Time  string `json:"time"`
	Value int    `json:"value"`
}
