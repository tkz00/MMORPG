package connection

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"backend/api/dtos"
	"backend/config"
	"backend/pkg/game/entities"
	"backend/pkg/game/gameplay"
	"backend/pkg/utils"

	"golang.org/x/net/websocket"
)

type AddClientData struct {
	Client    *websocket.Conn
	Character string
}

type NativeServer struct {
	clients      map[*websocket.Conn]bool
	addClient    chan AddClientData
	removeClient chan *websocket.Conn
	broadcast    chan []byte
	gameState    *entities.GameState
}

func (ws *NativeServer) newServer() Server {
	gamestate := gameplay.StartGameState()

	return &NativeServer{
		clients:      make(map[*websocket.Conn]bool),
		addClient:    make(chan AddClientData),
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
		case clientData := <-server.addClient:
			player := gameplay.AddPlayer(server.gameState, clientData.Client, clientData.Character)
			server.clients[clientData.Client] = true

			response := CreateWebSocketResponse(dtos.CharacterToDTO(player))
			message := response.Serialize()
			err := websocket.Message.Send(clientData.Client, message)
			if err != nil {
				log.Println("Error broadcasting message:", err)
				return
			}

			// Send initial gamestate to client
			gameStateDTO := dtos.GameStateToDTO(*server.gameState)
			webSocketResponse := CreateWebSocketResponse(gameStateDTO)
			if err := websocket.Message.Send(clientData.Client, webSocketResponse.Serialize()); err != nil {
				log.Println("Error broadcasting initial message:", err)
				return
			}

			log.Println("Client connected", clientData.Client.RemoteAddr())
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
	ticker := time.NewTicker(config.TICKER_TIME)
	defer ticker.Stop()
	previousUpdateTime := time.Now()

	for range ticker.C {
		currentUpdateTime := time.Now()
		deltaTime := currentUpdateTime.Sub(previousUpdateTime)
		previousUpdateTime = currentUpdateTime
		gameplay.UpdateState(server.gameState, deltaTime.Seconds())

		gameStateDTO := dtos.GameStateDiff(*server.gameState)
		webSocketResponse := CreateWebSocketResponse(gameStateDTO)
		server.broadcast <- webSocketResponse.Serialize()
	}
}

func (server *NativeServer) handleWebSocket(conn *websocket.Conn) {
	// Get query parameter "character"
	character := conn.Request().URL.Query().Get("character")

	server.addClient <- AddClientData{conn, character}
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
		case "ability_cast":
			server.handleAbilityCast(conn, message.Body.(dtos.AbilityCastDTO))
		case "respawn":
			server.handleRespawn(conn)
		case "use_item":
			server.handleUseItem(conn, message.Body.(dtos.UseItemDTO))
		case "equip_item":
			server.handleEquipItem(conn, message.Body.(dtos.EquipItemDTO))
		case "unequip_item":
			server.handleUnequipItem(conn, message.Body.(dtos.UnequipItemDTO))
		default:
			fmt.Printf("Unknown message type: %s\n", message.ActionType)
		}
	}
}

func (server *NativeServer) handlePlayerMovement(
	client *websocket.Conn,
	positionDTO dtos.PositionDTO,
) {
	position := *dtos.PositionDTOToEntity(positionDTO)

	// this should be just one "action", one line of code, it gives access to two different things while the entry point should be single
	player := server.gameState.GetPlayerByConn(client)
	player.EnqueueMovementAction(position)
}

func (server *NativeServer) handleAbilityCast(
	client *websocket.Conn,
	abilityCastDTO dtos.AbilityCastDTO,
) {
	player := server.gameState.GetPlayerByConn(client)
	castParameters := make(map[entities.Targeting]interface{})
	for key, value := range abilityCastDTO.AbilityParameters {
		switch key {
		case dtos.TargetPosition:
			// Cast value to map[string]interface{} first
			valueMap, ok := value.(map[string]interface{})
			if !ok {
				// Handle the error if it's not a map[string]interface{}
				fmt.Printf("expected map[string]interface{} but got %T", value)
			}

			// Now extract and cast the individual elements to float64
			x, xOk := valueMap["x"].(float64)
			z, zOk := valueMap["z"].(float64)
			if !xOk || !zOk {
				// Handle the error if x or z are not float64
				fmt.Printf(
					"expected x and z to be float64 but got %T and %T",
					valueMap["x"],
					valueMap["z"],
				)
			}

			// Now you can use x and z as float64 values to create the vector
			coordinates := *utils.NewVector2(x, z)
			castParameters[entities.Coordinates] = coordinates

		case dtos.TargetId:
			castParameters[entities.Target] = value
		}
	}
	player.EnqueueAbilityCastAction(abilityCastDTO.Id, castParameters)
}

func (server *NativeServer) handleRespawn(client *websocket.Conn) {
	player := server.gameState.GetPlayerByConn(client)
	if !player.IsAlive() {
		player.HealthVariation(player.GetMaxHealth())
		player.MoveTowards(*utils.NewVector2(0, 0))
		player.SetPosition(*utils.NewVector2(0, 0))
		response := CreateWebSocketResponse(dtos.CharacterToDTO(*player))
		response.ActionType = "respawn"
		message := response.Serialize()
		websocket.Message.Send(client, message)
	}
}

func (server *NativeServer) handleUseItem(client *websocket.Conn, useItemDTO dtos.UseItemDTO) {
	player := server.gameState.GetPlayerByConn(client)
	player.UseItem(useItemDTO.ItemId, useItemDTO.TargetId, server.gameState)
}

func (server *NativeServer) handleEquipItem(
	client *websocket.Conn,
	equipItemDTO dtos.EquipItemDTO,
) {
	player := server.gameState.GetPlayerByConn(client)
	player.EquipItem(equipItemDTO.ItemId)
}

func (server *NativeServer) handleUnequipItem(
	client *websocket.Conn,
	unequipItemDTO dtos.UnequipItemDTO,
) {
	player := server.gameState.GetPlayerByConn(client)
	player.UnequipItem(unequipItemDTO.ItemId)
}
