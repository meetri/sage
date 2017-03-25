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

func RegisterConnectCommands(app *cli.App) {

	app.Commands = append(app.Commands, cli.Command{
		Name:    "connect",
		Aliases: []string{"co", "hijack"},
		Usage:   "connect to a running container",
		Action:  ConnectContainer,
	})

}

func ConnectContainer(c *cli.Context) (err error) {

	hostlist := core.Config().Hosts()
	certpath_global := hostlist.FindDefault("certpath", "")
	timeout := hostlist.FindDefaultInt("timeout", 5)
	hosts := hostlist.Find("hosts")

	if hosts != nil {
		ac := clients.GetAllContainers(hosts.([]interface{}), certpath_global, timeout)
		mh, mc, err := config.SearchContainers(c.Args(), ac)

		if err == nil {
			dockercli := proxy.CreateDockerCli(mh)
			dockercli.Connect(mc.ID)
		} else {
			fmt.Printf("%s", err)
		}

	}

	return

}
