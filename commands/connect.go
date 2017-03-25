package commands

import (
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/docker/docker/api/types"
	"github.com/meetri/sage/clients"
	"github.com/meetri/sage/config"
	"github.com/meetri/sage/core"
	"github.com/meetri/sage/proxy"
	"github.com/meetri/ymltree"
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

		matched := 0
		var matched_container types.Container
		var matched_host map[string]string

		for _, cdata := range ac {

			if cdata != nil {

				containers := cdata.(ymltree.Map).Find("container")
				hostdetails := cdata.(ymltree.Map).Find("host").(map[string]string)

				for _, container := range containers.([]types.Container) {

					terms := config.GetContainerTerms(container)

					match := false

					if len(c.Args()) == 0 {
						match = true
					} else {
						for _, arg := range c.Args() {
							if config.MatchTerms(arg, terms) {
								match = true
								break
							}
						}
					}

					if match {
						matched++
						matched_container = container
						matched_host = hostdetails
					}

				}

			}

		}

		if matched == 1 {
			_ = matched_container
			dockercli := proxy.CreateDockerCli(matched_host)
			dockercli.Connect(matched_container.ID)
		} else {
			fmt.Printf("Found %d results that match your query, please refine your search\n", matched)
		}

	}

	return

}
