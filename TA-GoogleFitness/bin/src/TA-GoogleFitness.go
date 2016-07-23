package main

import (
	"bufio"
	"os"

	"github.com/AndyNortrup/GoSplunk"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	input := &GoogleFitnessInput{reader: reader, writer: os.Stdout}
	handleArgs(input)
}

func handleArgs(input splunk.ModularInputHandler) {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--scheme":
			input.ReturnScheme()
		case "--validate-arguments":
			input.ValidateScheme()
		}
	} else {
		input.StreamEvents()
	}
}
