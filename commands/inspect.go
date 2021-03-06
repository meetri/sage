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

func RegisterInspectCommand(app *cli.App) {

	app.Commands = append(app.Commands, cli.Command{
		Name:    "inspect",
		Aliases: []string{"inspect", "i"},
		Usage:   "inspect container",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "details", Usage: "show all containers"},
			cli.BoolFlag{Name: "f,follow"},
			cli.StringFlag{Name: "since", Usage: "Show logs since timestamp"},
			cli.IntFlag{Name: "tail", Usage: "Number of lines t o show from the end of the logs ( default \"all\")"},
			cli.BoolFlag{Name: "t,timestamps", Usage: "Show timestamps"},
		},
		Action: InspectContainer,
	})

}

func InspectContainer(c *cli.Context) (err error) {

	hostlist := core.Config().Hosts()
	certpath_global := hostlist.FindDefault("certpath", "")
	timeout := hostlist.FindDefaultInt("timeout", 5)
	hosts := hostlist.Find("hosts")

	if hosts != nil {
		ac := clients.GetAllContainers(hosts.([]interface{}), certpath_global, timeout)
		mh, mc, err := config.SearchContainers(c.Args(), ac)

		if err == nil {
			dockercli := proxy.CreateDockerCli(mh)
			dockercli.Inspect(mc.ID)
		} else {
			fmt.Printf("%s", err)
		}

	}

	return

}
