package main

import (
	"tkz00/backend/connection"
)

func main() {
	const PORT string = "3009"
	server := connection.CreateServer()
	server.StartConnection(PORT)
}
