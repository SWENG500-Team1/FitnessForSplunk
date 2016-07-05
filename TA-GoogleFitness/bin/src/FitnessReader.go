package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/fitness/v1"
)

type FitnessReader struct {
	conf         *oauth2.Config
	client       *http.Client
	refreshToken string
	accessToken  string
	expires      time.Time
	tokenType    string
}

func NewFitnessReader(clientID string,
	clientSecret string,
) *FitnessReader {

	fit := &FitnessReader{}

	fit.conf = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes: []string{fitness.FitnessActivityReadScope,
			fitness.FitnessBodyReadScope},
		Endpoint: google.Endpoint,
	}

	return fit
}

func (fit *FitnessReader) getTokenFromAccessCode(accessCode string) *oauth2.Token {
	tok, err := fit.conf.Exchange(oauth2.NoContext, accessCode)
	if err != nil {
		log.Fatal(err)
	}
	return tok
}

func (fit *FitnessReader) getTokenFromRefreshToken(
	refreshToken string,
	accessToken string,
	expiryStr string,
	tokenType string) *oauth2.Token {

	const example = "2006-01-02 15:04:05.00000000 -0700 MST"
	expires, err := time.Parse(example, expiryStr)
	if err != nil {
		log.Fatal(err)
	}

	tok := new(oauth2.Token)
	tok.RefreshToken = refreshToken
	tok.AccessToken = accessToken
	tok.Expiry = expires
	tok.TokenType = tokenType

	return tok
}

func (fit *FitnessReader) getClient(tok *oauth2.Token) *http.Client {
	if fit.client == nil {
		context := oauth2.NoContext
		fit.client = fit.conf.Client(context, tok)
	}
	return fit.client
}

func (fit *FitnessReader) GetDataSources(tok *oauth2.Token) []*fitness.DataSource {
	if fit.client == nil {
		fit.getClient(tok)
	}
	service, err := fitness.New(fit.client)
	if err != nil {
		log.Fatal(err)
	}
	dataSourceService := fitness.NewUsersDataSourcesService(service)
	call := dataSourceService.List("me")
	response, err := call.Do()
	if err != nil {
		log.Fatal(err)
	}

	return response.DataSource
}

func (fit *FitnessReader) GetDataSet(tok *oauth2.Token,
	startTime time.Time,
	endTime time.Time,
	dataSource fitness.DataSource) *fitness.Dataset {

	if fit.client == nil {
		fit.getClient(tok)
	}

	dataSetId := strconv.FormatInt(startTime.UnixNano(), 10) + "-" +
		strconv.FormatInt(endTime.UnixNano(), 10)

	service, err := fitness.New(fit.getClient(tok))
	if err != nil {
		log.Fatal(err)
	}

	dataSetService := fitness.NewUsersDataSourcesDatasetsService(service)
	request := dataSetService.Get("me", dataSource.DataStreamId, dataSetId)
	resp, err := request.Do()
	if err != nil {
		log.Fatal(err)
	}

	return resp
}
