package commands

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/urfave/cli"
	"os"
	"text/tabwriter"
)

func ListContainers(*cli.Context) (err error) {

	cli, err := client.NewClient("tcp://192.168.59.104:2376", "v1.21", nil, nil)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 4, ' ', 00)
	fmt.Fprintf(w, "CONTAINER ID\tNAME\tIMAGE\tCOMMAND\tSTATUS\n")
	for _, container := range containers {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", container.ID[0:16], container.Names[0], container.Image, container.Command, container.Status)
	}
	w.Flush()
	return

}
