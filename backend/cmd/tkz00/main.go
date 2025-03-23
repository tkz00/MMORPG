package main

import (
	"tkz00/backend/connection"
	"tkz00/backend/pkg/configurator"
)

func main() {
	configurator.RunSeeds()

	// Start configurator server
	go configurator.Run()

	// Start game server
	const PORT string = "3009"
	server := connection.CreateServer()
	go server.StartConnection(PORT)

	// Block main from exiting by waiting indefinitely
	select {}
}
