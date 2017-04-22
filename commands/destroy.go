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

func RegisterDestroyCommand(app *cli.App) {

	app.Commands = append(app.Commands, cli.Command{
		Name:   "destroy",
		Usage:  "stop and remove container",
		Action: DestroyContainer,
	})

}

func DestroyContainer(c *cli.Context) (err error) {

	hostlist := core.Config().Hosts()
	certpath_global := hostlist.FindDefault("certpath", "")
	timeout := hostlist.FindDefaultInt("timeout", 5)
	hosts := hostlist.Find("hosts")

	if hosts != nil {
		ac := clients.GetAllContainers(hosts.([]interface{}), certpath_global, timeout)
		mh, mc, err := config.SearchContainers(c.Args(), ac)

		if err == nil {
			_ = mh
			orch, err := proxy.Create("docker", mh)
			if err == nil {
				orch.Destroy(orch.GetId(mc))
			} else {
				fmt.Printf("%s\n", err)
			}

		} else {
			fmt.Printf("%s", err)
		}

	}

	return

}
