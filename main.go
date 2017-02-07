package main

import (
	"fmt"
	"github.com/meetri/sage/config"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {

	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("Panic: %v", r)
		}
	}()

	processCli()

}

func loadAppConfig() {
	cfg := config.Load("/Users/meetri/.vega.yml")
	cfg.Select("environment")
	log.Fatalf("CONFIG_ROOT = %s", cfg.Find("environment/CONFIG_ROOT"))
}

func listContainers(*cli.Context) (err error) {
	fmt.Printf("Listing your containers\n")

	return
}

func processCli() {

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
		fmt.Printf("hey buddy\n")
		return nil
	}

	app.Run(os.Args)

}
