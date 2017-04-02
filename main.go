package main

import (
	"github.com/meetri/sage/commands"
	"github.com/meetri/sage/core"
	"github.com/urfave/cli"
	//"log"
	"os"
)

func init() {
	core.Config()
}

func setupAppDefinitions() *cli.App {

	app := cli.NewApp()

	app.Name = "sage"
	app.Version = "0.1.0"
	app.Usage = "sage [OPTIONS] [ARGUMENTS]"
	app.Action = func(c *cli.Context) error {
		return nil
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "host",
			Value:       "all",
			Usage:       "run on selected hosts",
			Destination: &commands.SelectedHosts,
		},
	}

	commands.RegisterListCommands(app)
	commands.RegisterConnectCommands(app)
	commands.RegisterLogsCommand(app)
	commands.RegisterInspectCommand(app)
	commands.RegisterDockerCommand(app)
	commands.RegisterRestartCommand(app)
	commands.RegisterStopCommand(app)
	commands.RegisterStartCommand(app)
	commands.RegisterRemoveCommand(app)

	return app

}

func main() {

	/*
		defer func() {
			if r := recover(); r != nil {
				log.Fatalf("Panic: %v", r)
			}
		}()*/

	app := setupAppDefinitions()
	app.Run(os.Args)

}
