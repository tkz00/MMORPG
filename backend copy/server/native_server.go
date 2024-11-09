package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"tkz00/backend/api/dtos"
	"tkz00/backend/pkg/game"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

const TICKER_TIME = 50 * time.Millisecond

type NativeServer struct {
	clients      map[*websocket.Conn]string
	addClient    chan *websocket.Conn
	removeClient chan *websocket.Conn
	broadcast    chan []byte
	gameState    game.GameState
}

func CreateServer() NativeServer {
	return NativeServer{
		clients:      make(map[*websocket.Conn]string),
		addClient:    make(chan *websocket.Conn),
		removeClient: make(chan *websocket.Conn),
		broadcast:    make(chan []byte),
		gameState:    game.NewGameState(),
	}
}

func (ws NativeServer) StartConnection(port string) {
	http.Handle("/ws", websocket.Handler(ws.handleWebSocket))
	go ws.readLoop()
	go ws.serverLoop()

	log.Println("WebSocket server running on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func (server *NativeServer) readLoop() {
	for {
		select {
		case client := <-server.addClient:
			id := uuid.New().String()
			server.clients[client] = id
			server.gameState = server.gameState.AddPlayer(id)
			response := CreateWebSocketResponse(dtos.CharacterToDTO(id, server.gameState.Players()[id]))
			message := response.Serialize()
			err := websocket.Message.Send(client, message)
			if err != nil {
				log.Println("Error broadcasting message:", err)
				return
			}
		case client := <-server.removeClient:
			id := server.clients[client]
			delete(server.clients, client)
			server.gameState = server.gameState.RemovePlayer(id)
		case message := <-server.broadcast:
			for client := range server.clients {
				err := websocket.Message.Send(client, message)
				if err != nil {
					log.Println("Error broadcasting message:", err)
					return
				}
			}
		}
	}
}

func (server *NativeServer) serverLoop() {
	ticker := time.NewTicker(TICKER_TIME)
	defer ticker.Stop()
	// previousUpdateTime := time.Now()

	for range ticker.C {
		// currentUpdateTime := time.Now()
		// deltaTime := currentUpdateTime.Sub(previousUpdateTime)
		// previousUpdateTime = currentUpdateTime
		// gameplay.UpdateState(&server.gameState, deltaTime.Seconds())

		gameStateDTO := dtos.GameStateToDTO(server.gameState)
		webSocketResponse := CreateWebSocketResponse(gameStateDTO)
		server.broadcast <- webSocketResponse.Serialize()
	}
}

func (server *NativeServer) handleWebSocket(conn *websocket.Conn) {
	server.addClient <- conn
	defer func() {
		conn.Close()
		server.removeClient <- conn
	}()
	for {
		var data []byte
		err := websocket.Message.Receive(conn, &data)
		if err != nil {
			log.Println("Error reading message from client:", err)
			break
		}

		var message WebSocketMessage
		if err := json.Unmarshal(data, &message); err != nil {
			fmt.Println("Error decoding message:", err)
			return
		}

		switch message.ActionType {
		case "position":
			server.handlePlayerMovement(conn, message.Body.(dtos.PositionDTO))
		// case "ability_cast":
		// 	server.handleAbilityCast(conn, message.Body.(dtos.AbilityCastDTO))
		// case "respawn":
		// 	server.handleRespawn(conn)
		// case "use_item":
		// 	server.handleUseItem(conn, message.Body.(dtos.UseItemDTO))
		default:
			fmt.Printf("Unknown message type: %s\n", message.ActionType)
		}
	}
}

func (server *NativeServer) handlePlayerMovement(client *websocket.Conn, positionDTO dtos.PositionDTO) {
	position := dtos.PositionDTOToEntity(positionDTO)
	server.gameState = server.gameState.MovePlayer(server.clients[client], position)
}

// func (server *NativeServer) handleAbilityCast(client *websocket.Conn, abilityCastDTO dtos.AbilityCastDTO) {
// }

// func (server *NativeServer) handleRespawn(client *websocket.Conn) {}

// func (server *NativeServer) handleUseItem(client *websocket.Conn, useItemDTO dtos.UseItemDTO) {
// }
