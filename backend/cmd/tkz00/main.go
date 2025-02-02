package main

import (
	"time"
	"tkz00/backend/connection"
	"tkz00/backend/pkg/configurator"
)

func main() {
	// Start configurator server
	go configurator.Run()

	// need a configurator setup method that's blocking and then run it in parallel to the game
	time.Sleep(5 * time.Second)

	// Start game server
	const PORT string = "3009"
	server := connection.CreateServer()
	go server.StartConnection(PORT)

	// Block main from exiting by waiting indefinitely
	select {}
}
