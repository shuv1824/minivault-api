package main

import (
	"os"
	"time"

	"github.com/shuv1824/minivault/server"
)

func main() {
	startTime := time.Now()
	mode := os.Getenv("MODE")

	if mode == "API" {
		server.Run(startTime)
	} else if mode == "CLI" {
		server.RunCLI()
	} else {
		panic("Invalid mode")
	}
}
