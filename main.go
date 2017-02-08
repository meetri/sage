package main

import (
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/meetri/sage/config"
	"text/tabwriter"
	//"github.com/urfave/cli"
	"context"
	"log"
	"os"
)

var appConfig config.Map

func hiDocker() {

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
	//fmt.Fprintf(w, "Goodbye\tMoon\tYou'll always be my boo\n")
	//fmt.Fprintf(w, "What a good thing\tDon't hate the playa hate the game\tJust think about good things\n")

	//need to come back to this
	_ = containers

}

func loadAppConfig() {

	var appConfig config.Tree

	if err := appConfig.SmartLoad(fmt.Sprintf("%s/.sage.yaml", os.Getenv("HOME"))); err != nil {
		log.Fatalf("burp")
	}

	appConfig.Select("main")
	appConfig.Sel.Save()
	//spew.Dump(appConfig.Sel)
}

func main() {

	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("Panic: %v", r)
		}
	}()

	app := setupAppDefinitions()
	app.Run(os.Args)

	loadAppConfig()

	// hiDocker()

}
