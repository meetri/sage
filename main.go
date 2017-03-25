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

	commands.RegisterListCommands(app)
	commands.RegisterConnectCommands(app)

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
