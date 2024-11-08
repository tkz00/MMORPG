package main

import "tkz00/backend/server"

func main() {
	const PORT string = "3009"
	server := server.CreateServer()
	server.StartConnection(PORT)
}
