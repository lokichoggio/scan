package main

import (
	"log"

	"scan-eth/cmd/scan/app"
)

func main() {
	command := app.NewScanCommand()
	if err := command.Execute(); err != nil {
		log.Fatalf("cmd Execute error: %s", err)
	}
}
