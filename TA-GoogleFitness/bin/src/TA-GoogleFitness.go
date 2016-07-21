package main

import (
	"log"
	"os"

	"github.com/AndyNortrup/GoSplunk"
)

func main() {
	input, err := NewGoogleFitnessInput(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatalf("Unable to create GoogleFitnessInput: %v", err)
	}
	handleArgs(input)
}

func handleArgs(input splunk.ModularInputHandler) {
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
