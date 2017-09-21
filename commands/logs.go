package commands

import (
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/meetri/sage/clients"
	"github.com/meetri/sage/config"
	"github.com/meetri/sage/core"
	"github.com/meetri/sage/proxy"
	"github.com/urfave/cli"
)

func RegisterLogsCommand(app *cli.App) {

	app.Commands = append(app.Commands, cli.Command{
		Name:    "logs",
		Aliases: []string{"log", "l"},
		Usage:   "view container logs",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "details", Usage: "show all containers"},
			cli.BoolFlag{Name: "f,follow"},
			cli.StringFlag{Name: "since", Usage: "Show logs since timestamp"},
			cli.StringFlag{Name: "tail", Usage: "Number of lines t o show from the end of the logs ( default \"all\")"},
			cli.BoolFlag{Name: "t,timestamps", Usage: "Show timestamps"},
		},
		Action: LogContainer,
	})

}

func LogContainer(c *cli.Context) (err error) {

	hostlist := core.Config().Hosts()
	certpath_global := hostlist.FindDefault("certpath", "")
	timeout := hostlist.FindDefaultInt("timeout", 5)
	hosts := hostlist.Find("hosts")

	if hosts != nil {
		ac := clients.GetAllContainers(hosts.([]interface{}), certpath_global, timeout)
		mh, mc, err := config.SearchContainers(c.Args(), ac)

		if err == nil {
			orch, err := proxy.Create("docker", mh)
			if err == nil {
				args := make(map[string]interface{})
				args["follow"] = c.Bool("follow")
				args["tail"] = c.String("tail")
				orch.Logs(mc.ID, args)
			} else {
				fmt.Printf("%s\n", err)
			}
		} else {
			fmt.Printf("%s", err)
		}

	}

	return

}
