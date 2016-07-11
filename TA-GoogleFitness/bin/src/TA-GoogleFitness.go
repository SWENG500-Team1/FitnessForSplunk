package main

import (
	"encoding/xml"
	"os"
)

func main() {
	input := &GoogleFitnessInput{}
	handleArgs(input)
}

type Scheme struct {
	XMLName               xml.Name   `xml:"scheme"`
	Title                 string     `xml:"title"`
	Description           string     `xml:"description"`
	UseExternalValidation bool       `xml:"use_external_validation"`
	StreamingMode         string     `xml:"streaming_mode"`
	Args                  []Argument `xml:"args"`
}

type Argument struct {
	XMLName     xml.Name `xml:"arg"`
	Name        string   `xml:"name,attr"`
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
}

func handleArgs(input ModularInputHandler) {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--scheme":
			input.ReturnScheme()
			return
		case "--validate-arguments":
			input.ValidateScheme()
			return
		}
	}
	input.StreamEvents()
}

type ModularInputHandler interface {
	ReturnScheme()
	ValidateScheme()
	StreamEvents()
}
