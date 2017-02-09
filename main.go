package main

import (
	//"github.com/davecgh/go-spew/spew"
	"log"
	"os"
)

func main() {

	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("Panic: %v", r)
		}
	}()

	app := setupAppDefinitions()
	app.Run(os.Args)

	AppConfig.Data.Save()
	// hiDocker()

}
