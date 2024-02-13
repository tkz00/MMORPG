package main

import (
	"unnamed-mmo/backend/connection"
)

func main() {
	const PORT string = "3009"
	server := connection.CreateServer()
	server.StartConnection(PORT)
}

