package main

import (
	"log"
	"personal-project-core/services/wire"
)

func main() {
	// load config -> connect database -> new service context -> new server
	server, err := wire.InitializeServer()
	if err != nil {
		log.Fatalf("error initializing server: %v", err)
	}
	server.SetupRoutes()
	server.Run()
}
