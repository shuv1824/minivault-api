package main

import (
	"os"

	"github.com/shuv1824/minivault/server"
)

func main() {
	mode := os.Getenv("MODE")

	if mode == "API" {
		server.Run()
	} else if mode == "CLI" {
		server.RunCLI()
	} else {
		panic("Invalid mode")
	}
}
