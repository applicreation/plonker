package main

import (
	"github.com/applicreation/plonker/command"
	"github.com/applicreation/plonker/config"
	"github.com/applicreation/plonker/connection"
	"github.com/mitchellh/cli"
	"log"
	"os"
)

func main() {
	conf := config.Config{}

	conf.Load()

	db := connection.GetConnection(&conf.Connection)

	defer db.Close()

	conn := connection.Connection{
		conf.Connection,
		db,
	}

	plonker := &cli.CLI{
		Args: os.Args[1:],
		Commands: map[string]cli.CommandFactory{
			"dry-run": func() (cli.Command, error) {
				return &command.DryRunCommand{
					Config:     &conf,
					Connection: &conn,
				}, nil
			},
		},
	}

	exitCode, err := plonker.Run()
	if err != nil {
		log.Printf("Error executing CLI: %s", err.Error())
		os.Exit(1)
	}

	os.Exit(exitCode)
}
