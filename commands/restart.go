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

func RegisterRestartCommand(app *cli.App) {

	app.Commands = append(app.Commands, cli.Command{
		Name:   "restart",
		Usage:  "restart container",
		Action: RestartContainer,
	})

}

func RestartContainer(c *cli.Context) (err error) {

	hostlist := core.Config().Hosts()
	certpath_global := hostlist.FindDefault("certpath", "")
	timeout := hostlist.FindDefaultInt("timeout", 5)
	hosts := hostlist.Find("hosts")

	if hosts != nil {
		ac := clients.GetAllContainers(hosts.([]interface{}), certpath_global, timeout)
		mh, mc, err := config.SearchContainers(c.Args(), ac)

		if err == nil {
			_ = mh
			terms := config.GetContainerTerms(mc)

			orch, err := proxy.Create(terms["orchestration"], mh)
			if err == nil {
				orch.Restart(orch.GetId(mc))
			} else {
				fmt.Printf("%s\n", err)
			}

		} else {
			fmt.Printf("%s", err)
		}

	}

	return

}
