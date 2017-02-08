package main

import (
	"fmt"
	"github.com/urfave/cli"
)

func listContainers(*cli.Context) (err error) {
	fmt.Printf("Listing your containers\n")
	return
}

func setupAppDefinitions() *cli.App {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Load configuration from `FILE`",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "ls",
			Aliases: []string{"l"},
			Usage:   "list running containers",
			Action:  listContainers,
		},
	}

	app.Name = "sage"
	app.Usage = "Docker Orchestrator"
	app.Action = func(c *cli.Context) error {
		return nil
	}

	return app

}
