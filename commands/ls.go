package commands

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/docker/docker/api/types"
	"github.com/meetri/sage/clients"
	"github.com/meetri/sage/config"
	"github.com/meetri/sage/core"
	"github.com/meetri/ymltree"
	"github.com/urfave/cli"
	"os"
	"text/tabwriter"
)

func RegisterListCommands(app *cli.App) {

	app.Commands = append(app.Commands, cli.Command{
		Name:    "list",
		Aliases: []string{"ls", "ps"},
		Usage:   "list running containers",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "a,all", Usage: "show all containers"},
			cli.BoolFlag{Name: "q,quiet"},
			cli.BoolFlag{Name: "n,name"},
			cli.BoolFlag{Name: "d,dump"},
		},
		Action: ListContainers,
	})

}

func ListContainers(c *cli.Context) (err error) {

	hostlist := core.Config().Hosts()
	certpath_global := hostlist.FindDefault("certpath", "")
	timeout := hostlist.FindDefaultInt("timeout", 5)
	hosts := hostlist.Find("hosts")

	//Allow overrides in ~/.sage/config.yml
	//termlist := []string{"id", "hostalias", "name", "network", "address", "image", "command", "status", "state", "ports"}
	termlist := []string{"id", "hostalias", "name", "network", "image", "status"}
	termmap := map[string]string{
		"id":        "CONTAINER ID",
		"hostalias": "HOST ALIAS",
		"name":      "NAME",
		"network":   "NETWORK",
		"address":   "ADDRESS",
		"image":     "IMAGE",
		"command":   "COMMAND",
		"status":    "STATUS",
		"state":     "STATE",
		"ports":     "PORTS",
		"labels":    "LABELS",
	}

	if hosts != nil {
		ac := clients.GetAllContainers(hosts.([]interface{}), certpath_global, timeout)

		w := tabwriter.NewWriter(os.Stdout, 1, 1, 4, ' ', 00)

		if !c.Bool("dump") {
			fmt.Fprintf(w, config.WriteTerms(termmap, termlist))
		}

		for _, cdata := range ac {

			if cdata != nil {

				containers := cdata.(ymltree.Map).Find("container")
				hostname := cdata.(ymltree.Map).FindDefault("hostname", "")
				hostalias := cdata.(ymltree.Map).FindDefault("hostalias", hostname)

				for _, container := range containers.([]types.Container) {

					terms := config.GetContainerTerms(container)
					terms["hostalias"] = hostalias
					terms["host"] = hostalias

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
						if c.Bool("dump") {
							spew.Dump(container)
						} else {
							termsT := config.TruncateContainerTerms(terms)
							fmt.Fprintf(w, config.WriteTerms(termsT, termlist))
						}
					}

				}

			}

		}
		w.Flush()

	}

	return

}
