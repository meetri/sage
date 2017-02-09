package main

import (
	"fmt"
	"github.com/meetri/sage/commands"
	"github.com/meetri/sage/config"
	"github.com/urfave/cli"
	"log"
	"os"
)

var AppConfig config.Tree

func init() {

	if err := AppConfig.SmartLoad(fmt.Sprintf("%s/.sage.yaml", os.Getenv("HOME"))); err != nil {
		log.Fatalf("burp")
	} else {
		AppConfig.Select("main")
	}

}

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
			Name:    "list",
			Aliases: []string{"ls", "ps"},
			Usage:   "list running containers",
			Action:  commands.ListContainers,
		},
	}

	app.Name = "sage"
	app.Usage = "Docker Orchestrator"
	app.Action = func(c *cli.Context) error {
		return nil
	}

	return app

}
