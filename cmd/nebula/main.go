package main

import (
	"log"
	"os"

	"github.com/puppetlabs/nebula/pkg/cmd"
)

func main() {
	command, err := cmd.NewRootCommand()
	if err != nil {
		log.Fatal(err)
	}

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
