package commands

import (
	//"fmt"
	"strings"
	//"github.com/davecgh/go-spew/spew"
	//"github.com/meetri/sage/clients"
	//"github.com/meetri/sage/config"
	"github.com/meetri/sage/core"
	"github.com/meetri/sage/proxy"
	"github.com/meetri/ymltree"
	"github.com/urfave/cli"
)

func RegisterDockerCommand(app *cli.App) {

	app.Commands = append(app.Commands, cli.Command{
		Name:    "docker",
		Aliases: []string{"d", "doc"},
		Usage:   "docker proxy",
		Action:  DockerProxy,
	})

}

func DockerProxy(c *cli.Context) (err error) {

	hostlist := core.Config().Hosts()
	certpath_global := hostlist.FindDefault("certpath", "")
	hosts := hostlist.Find("hosts")

	if hosts != nil {

		for _, host := range hosts.([]interface{}) {
			match := host.(ymltree.Map).FindDefault("host", "")
			match = host.(ymltree.Map).FindDefault("alias", match)
			if strings.Contains(match, SelectedHosts) || SelectedHosts == "all" {
				proxconf := make(map[string]string)
				proxconf["hostname"] = host.(ymltree.Map).FindDefault("host", "")
				proxconf["certpath"] = host.(ymltree.Map).FindDefault("certpath", certpath_global)
				proxconf["binary"] = host.(ymltree.Map).FindDefault("binary", "docker")
				docprox := proxy.CreateDockerCli(proxconf)
				arglist := strings.Split(c.Args()[0], " ")
				docprox.Proxy(arglist...)
			}
		}

	}

	return

}
