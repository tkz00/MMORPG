package main

import (
	"tkz00/backend/connection"
	"tkz00/backend/pkg/configurator"
)

func main() {
	// Start game server
	const PORT string = "3009"
	server := connection.CreateServer()
	go server.StartConnection(PORT) // Run the game server in a goroutine

	// Start configurator server
	go configurator.Run() // Run the configurator server in a goroutine

	// Block main from exiting by waiting indefinitely
	select {}
}
