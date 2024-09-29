package connection

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"tkz00/backend/api/dtos"
	"tkz00/backend/pkg/game/entities"
	"tkz00/backend/pkg/game/gameplay"
	"tkz00/backend/pkg/utils"

	"golang.org/x/net/websocket"
)

const TICKER_TIME = 50 * time.Millisecond

type NativeServer struct {
	clients      map[*websocket.Conn]bool
	addClient    chan *websocket.Conn
	removeClient chan *websocket.Conn
	broadcast    chan []byte
	gameState    entities.GameState
}

func (ws *NativeServer) newServer() Server {
	gamestate := gameplay.StartGameState()

	return &NativeServer{
		clients:      make(map[*websocket.Conn]bool),
		addClient:    make(chan *websocket.Conn),
		removeClient: make(chan *websocket.Conn),
		broadcast:    make(chan []byte),
		gameState:    gamestate,
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
			player := gameplay.AddPlayer(&server.gameState, client)
			server.clients[client] = true
			response := CreateWebSocketResponse(*dtos.GetMapper().CharacterToDTO(player))
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

// the broadcast function should just broadcast, the updating of the state should be handled somewhere else
// actually, native server shouldn't know anything about game state, it should only receive messages that it should send to the clients, but how then would client that connect to the server be convereted to players?
func (server *NativeServer) broadcastGameState() {
	ticker := time.NewTicker(TICKER_TIME)
	defer ticker.Stop()
	previousUpdateTime := time.Now()

	for range ticker.C {
		currentUpdateTime := time.Now()
		deltaTime := currentUpdateTime.Sub(previousUpdateTime)
		previousUpdateTime = currentUpdateTime
		gameplay.UpdateState(&server.gameState, deltaTime.Seconds())

		gameStateDTO := *dtos.GetMapper().GameStateToDTO(server.gameState)
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
		case "Position":
			server.handlePlayerMovement(conn, message.Body.(dtos.PositionDTO))
		case "AbilityCast":
			server.handleAbilityCast(conn, message.Body.(dtos.AbilityCastDTO))
		case "Respawn":
			server.handleRespawn(conn)
		default:
			fmt.Printf("Unknown message type: %s\n", message.ActionType)
		}
	}
}

func (server *NativeServer) handlePlayerMovement(client *websocket.Conn, positionDTO dtos.PositionDTO) {
	position := *dtos.GetMapper().PositionDTOToEntity(positionDTO)

	// this should be just one "action", one line of code, it gives access to two different things while the entry point should be single
	player := server.gameState.GetPlayerByConn(client)
	player.EnqueueMovementAction(position)
}

func (server *NativeServer) handleAbilityCast(client *websocket.Conn, abilityCastDTO dtos.AbilityCastDTO) {
	player := server.gameState.GetPlayerByConn(client)
	entities.EnqueueAbilityCast(server.gameState, player, abilityCastDTO)
}

func (server *NativeServer) handleRespawn(client *websocket.Conn) {
	player := server.gameState.GetPlayerByConn(client)
	if !player.IsAlive() {
		player.HealthVariation(player.GetMaxHealth())
		player.MoveTowards(*utils.NewVector2(0, 0))
		player.SetPosition(*utils.NewVector2(0, 0))
		response := CreateWebSocketResponse(*dtos.GetMapper().CharacterToDTO(*player))
		response.ActionType = "Respawn"
		message := response.Serialize()
		websocket.Message.Send(client, message)
	}
}
