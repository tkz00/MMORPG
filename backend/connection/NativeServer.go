package connection

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"unnamed-mmo/backend/types"

	"golang.org/x/net/websocket"
)

const TICKER_TIME = 50 * time.Millisecond

type NativeServer struct {
	clients      map[*websocket.Conn]bool
	addClient    chan *websocket.Conn
	removeClient chan *websocket.Conn
	broadcast    chan []byte
	gameState    types.GameState
}

func (ws *NativeServer) newServer() Server {
	return &NativeServer{
		clients:      make(map[*websocket.Conn]bool),
		addClient:    make(chan *websocket.Conn),
		removeClient: make(chan *websocket.Conn),
		broadcast:    make(chan []byte),
		gameState:    types.StartGameState(),
	}
}

func (ws NativeServer) StartConnection(port string) {
	http.Handle("/ws", websocket.Handler(ws.handleWebSocket))
	go ws.readLoop()
	go ws.broadcastGameState()

	log.Println("WebSocket server running on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func (server *NativeServer) readLoop() {
	for {
		select {
		case client := <-server.addClient:
			playerDTO := server.gameState.AddPlayer(client)
			server.clients[client] = true
			response := types.CreateWebSocketResponse(playerDTO)
			message := response.Serialize()
	        err := websocket.Message.Send(client, message)
	        if err != nil {
	        	log.Println("Error broadcasting message:", err)
	        	return
	        }
			log.Println("Client connected", client.RemoteAddr())
		case client := <-server.removeClient:
			server.gameState.DeletePlayer(client)
			delete(server.clients, client)
			log.Println("Client disconnected", client.RemoteAddr())
		case message := <-server.broadcast: // THIS IS AN OBSERVER
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

func (server *NativeServer) broadcastGameState() {
	ticker := time.NewTicker(TICKER_TIME)
	defer ticker.Stop()

	for range ticker.C {
		server.gameState.UpdateState()
		gameStateDTO := server.gameState.GetGameState()
		webSocketResponse := types.CreateWebSocketResponse(gameStateDTO) 		
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

        var message types.WebSocketMessage
        if err := json.Unmarshal(data, &message); err != nil {
            fmt.Println("Error decoding message:", err)
            return
        }

        switch message.ActionType {
        case "Position":
            server.handlePlayerMovement(conn, message.Body.(types.PositionDTO))
        case "AbilityCast":
            server.handleAbilityCast(conn, message.Body.(types.AbilityCastDTO))
        default:
            fmt.Println("Unknown message type:", message.ActionType)
        }
    }
}

func (server *NativeServer) handlePlayerMovement(client *websocket.Conn, positionDTO types.PositionDTO) {
	position := *types.GetMapper().PositionDTOToEntity(positionDTO)

	server.gameState.MovePlayer(client, position)
}

func (server *NativeServer) handleAbilityCast(client *websocket.Conn, abilityCastDTO types.AbilityCastDTO) {

}

// {"x":3.6,"z":19.01}
