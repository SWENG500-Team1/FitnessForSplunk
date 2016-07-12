package main

import (
	"os"

	"github.com/AndyNortrup/GoSplunk"
)

func main() {
	input := &GoogleFitnessInput{}
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
