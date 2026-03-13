package connection

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"backend/api/dtos"
	"backend/config"
	"backend/pkg/game/entities"
	"backend/pkg/game/gameplay"
	"backend/pkg/utils"

	"github.com/gorilla/websocket"
	"github.com/samber/lo"
)

// connWriter serializes writes to a single WebSocket connection.
// gorilla/websocket connections do not support concurrent writes;
// this type ensures all sends go through a single dedicated goroutine.
type connWriter struct {
	ch chan []byte
}

func newConnWriter(conn *websocket.Conn) *connWriter {
	cw := &connWriter{ch: make(chan []byte, 256)}
	go func() {
		for msg := range cw.ch {
			if err := conn.WriteMessage(websocket.BinaryMessage, msg); err != nil {
				return
			}
		}
	}()
	return cw
}

func (cw *connWriter) send(data []byte) {
	select {
	case cw.ch <- data:
	default:
		log.Println("Warning: send buffer full, dropping message")
	}
}

type addClientReply struct {
	err    error
	writer *connWriter
}

type AddClientData struct {
	Client    *websocket.Conn
	Character string
	Reply     chan addClientReply
}

type NativeServer struct {
	clients             map[*websocket.Conn]*connWriter
	connectedCharacters map[*websocket.Conn]string
	addClient           chan AddClientData
	removeClient        chan *websocket.Conn
	tick                chan float64
	gameState           *entities.GameState
	upgrader            websocket.Upgrader
}

func NewServer() *NativeServer {
	gamestate := gameplay.StartGameState()

	return &NativeServer{
		clients:             make(map[*websocket.Conn]*connWriter),
		connectedCharacters: make(map[*websocket.Conn]string),
		addClient:           make(chan AddClientData),
		removeClient:        make(chan *websocket.Conn),
		tick:                make(chan float64),
		gameState:           gamestate,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

func (ws NativeServer) StartConnection(port string) {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.handleWebSocket(w, r)
	})
	go ws.gameLoop()
	go ws.tickLoop()

	log.Println("WebSocket server running on 127.0.0.1:" + port)
	log.Fatal(http.ListenAndServe("127.0.0.1:"+port, nil))
}

func (server *NativeServer) gameLoop() {
	for {
		select {
		case clientData := <-server.addClient:
			if lo.Contains(lo.Values(server.connectedCharacters), clientData.Character) {
				log.Println("Error, character already in use")
				clientData.Reply <- addClientReply{err: errors.New("character already in use")}
				continue
			}
			player := gameplay.AddPlayer(server.gameState, clientData.Client, clientData.Character)
			server.connectedCharacters[clientData.Client] = clientData.Character

			writer := newConnWriter(clientData.Client)

			response := CreateWebSocketResponse(dtos.CharacterToDTO(player))
			writer.send(response.Serialize())

			// Send initial game state to new client
			gameStateDTO := dtos.GameStateToDTO(*server.gameState)
			webSocketResponse := CreateWebSocketResponse(gameStateDTO)
			writer.send(webSocketResponse.Serialize())

			// Add to broadcast set only after initial messages are queued so the
			// client receives them in order before any tick updates.
			server.clients[clientData.Client] = writer
			clientData.Reply <- addClientReply{writer: writer}
			log.Println("Client connected", clientData.Client.RemoteAddr())

		case client := <-server.removeClient:
			if writer, ok := server.clients[client]; ok {
				close(writer.ch)
				delete(server.clients, client)
			}
			delete(server.connectedCharacters, client)
			server.gameState.DeletePlayer(client)
			log.Println("Client disconnected", client.RemoteAddr())

		case deltaTime := <-server.tick:
			gameplay.UpdateState(server.gameState, deltaTime)

			gameStateDTO := dtos.GameStateDiff(*server.gameState)
			webSocketResponse := CreateWebSocketResponse(gameStateDTO)
			message := webSocketResponse.Serialize()

			for _, writer := range server.clients {
				writer.send(message)
			}
		}
	}
}

func (server *NativeServer) tickLoop() {
	ticker := time.NewTicker(config.TICKER_TIME)
	defer ticker.Stop()
	previousUpdateTime := time.Now()

	for range ticker.C {
		currentUpdateTime := time.Now()
		deltaTime := currentUpdateTime.Sub(previousUpdateTime)
		server.tick <- deltaTime.Seconds()
		previousUpdateTime = currentUpdateTime
	}
}

func (server *NativeServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	character := r.URL.Query().Get("character")

	conn, err := server.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}

	reply := make(chan addClientReply)
	server.addClient <- AddClientData{Client: conn, Character: character, Reply: reply}

	result := <-reply
	if result.err != nil {
		errMsg := map[string]string{"error": "character already in use"}
		data, _ := json.Marshal(errMsg)
		if err := conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
			log.Println("Error sending rejection message:", err)
		}
		conn.Close()
		return
	}

	writer := result.writer

	defer func() {
		conn.Close()
		server.removeClient <- conn
	}()
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message from client: ", err, conn.RemoteAddr())
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
			server.handleRespawn(conn, writer)
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

func (server *NativeServer) handleRespawn(client *websocket.Conn, writer *connWriter) {
	player := server.gameState.GetPlayerByConn(client)
	if !player.IsAlive() {
		player.HealthVariation(player.GetMaxHealth())
		player.MoveTowards(*utils.NewVector2(0, 0))
		player.SetPosition(*utils.NewVector2(0, 0))
		response := CreateWebSocketResponse(dtos.CharacterToDTO(*player))
		response.ActionType = "respawn"
		writer.send(response.Serialize())
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
