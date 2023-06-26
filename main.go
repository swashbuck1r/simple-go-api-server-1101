package main

import (
	"log"

	"github.com/swashbuck1r/simple-go-api-server/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
