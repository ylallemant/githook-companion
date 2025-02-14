package main

import (
	"log"

	"github.com/ylallemant/githooks-butler/pkg/cli"
)

func main() {
	if err := cli.Command().Execute(); err != nil {
		log.Fatalf("error during execution: %v", err)
	}
}
