package commands

import (
	"fmt"
	"strings"
	//"github.com/davecgh/go-spew/spew"
	"github.com/meetri/sage/core"
	"github.com/meetri/ymltree"
	"github.com/urfave/cli"
)

func RegisterEnvCommand(app *cli.App) {

	app.Commands = append(app.Commands, cli.Command{
		Name:  "env",
		Usage: "env -s host1",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "s,server", Usage: "set default docker server"},
			cli.BoolFlag{Name: "d,dump"},
		},
		Action: SetEnvironment,
	})

}

func SetEnvironment(c *cli.Context) (err error) {

	hostlist := core.Config().Hosts()
	hosts := hostlist.Find("hosts")
	certpath := hostlist.Find("certpath").(string)

	srchTerm := c.Args()[0]

	if hosts != nil {

		for _, v := range hosts.([]interface{}) {

			hn := v.(ymltree.Map).FindDefault("host", "")
			alias := v.(ymltree.Map).FindDefault("alias", hn)
			cp := v.(ymltree.Map).FindDefault("certpath", certpath)

			if strings.Contains(alias, srchTerm) {
				fmt.Printf("export DOCKER_HOST=%s\n", hn)
				if len(cp) > 0 {
					fmt.Printf("export DOCKER_TLS_VERIFY=1\n")
					fmt.Printf("export DOCKER_CERT_PATH=%s\n", cp)
				}

			}
		}
	}

	return

}
