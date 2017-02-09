package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/meetri/sage/configurator"
	"log"
)

/*
func setup() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("Panic: %v", r)
		}
	}()

	app := setupAppDefinitions()
	app.Run(os.Args)

	AppConfig.Data.Save()
}*/

var configMap configurator.ConfigMap

func main() {

	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("Panic: %v", r)
		}
	}()

	//TODO: Get rid of this .. .debug mode
	_ = fmt.Println
	_ = spew.Dump

	var err error
	if configMap, err = configurator.Load("/Users/meetri/.sage.yaml"); err != nil {
		log.Fatalf("ERROR:" + err.Error())
	}

	out, _ := configMap.Select("main")
	fmt.Printf("%s", out.Export())

}
