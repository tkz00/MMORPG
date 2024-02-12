package main

import (
	"fmt"
	"log"
	"moreno-gaming/backend/src/models"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

// Define a simple WebSocket server struct
type WebSocketServer struct {
	clients      map[*websocket.Conn]bool
	addClient    chan *websocket.Conn
	removeClient chan *websocket.Conn
	broadcast    chan []byte
}

func (server *WebSocketServer) broadcastStateLoop() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Send a message containing the current time to all clients
			currentTime := time.Now().Format(time.RFC3339)
			for client := range server.clients {
				if err := websocket.Message.Send(client, currentTime); err != nil {
					log.Println("Error sending update to client:", err)
				}
			}
		}
	}
}

func (server *WebSocketServer) readLoop() {
	for {
		select {
		case client := <-server.addClient:
			server.clients[client] = true
			log.Println("Client connected", client.RemoteAddr())
		case client := <-server.removeClient:
			delete(server.clients, client)
			log.Println("Client disconnected", client.RemoteAddr())
		case message := <-server.broadcast:
			for client := range server.clients {
				err := websocket.Message.Send(client, string(message))
				if err != nil {
					log.Println("Error broadcasting message:", err)
					return
				}
			}
		}
	}
}

func (server *WebSocketServer) handleWebSocket(conn *websocket.Conn) {
	server.addClient <- conn
	defer func() { server.removeClient <- conn }()
	for {
		var message string
		err := websocket.Message.Receive(conn, &message)
		if err != nil {
			log.Println("Error reading message from client:", err)
			break
		}
		server.broadcast <- []byte(message)
		fmt.Println(string(message))
	}
}

func main() {
	player := models.Player {
		Position: models.Position{
	    	X: 20.0,
	    	Z: 20.0,
		},
	}
	fmt.Println(player)
	server := WebSocketServer{
		clients:      make(map[*websocket.Conn]bool),
		addClient:    make(chan *websocket.Conn),
		removeClient: make(chan *websocket.Conn),
		broadcast:    make(chan []byte),
	}

	http.Handle("/ws", websocket.Handler(server.handleWebSocket))

	go server.broadcastStateLoop()
	go server.readLoop()

	log.Println("WebSocket server running on :3009")
	log.Fatal(http.ListenAndServe(":3009", nil))
}

