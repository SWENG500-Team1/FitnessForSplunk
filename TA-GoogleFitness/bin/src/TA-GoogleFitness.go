package main

import (
	"bufio"
	"log"
	"os"

	"github.com/AndyNortrup/GoSplunk"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	input := &FitnessInput{reader: reader, writer: os.Stdout}
	handleArgs(input)
}

func handleArgs(input splunk.ModularInputHandler) {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--scheme":
			input.ReturnScheme()
		case "--validate-arguments":
			success, msg := input.ValidateScheme()
			if !success {
				log.Fatal(msg)
			}
		}
	} else {
		input.StreamEvents()
	}
}
