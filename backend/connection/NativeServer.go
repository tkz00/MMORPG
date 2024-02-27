package connection

import (
	"log"
	"net/http"
	"time"

	"unnamed-mmo/backend/types"

	"golang.org/x/net/websocket"
)


const TICKER_TIME = time.Second

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
			playerId := server.gameState.AddPlayer(client)
			server.clients[client] = true
			err := websocket.Message.Send(client, playerId)
			if err != nil {
				log.Println("Error broadcasting message:", err)
				return
			}
			log.Println("Client connected", client.RemoteAddr())
		case client := <-server.removeClient:
			server.gameState.DeletePlayer(client)
			delete(server.clients, client)
			log.Println("Client disconnected", client.RemoteAddr())
		case message := <-server.broadcast: //THIS IS AN OBSERVER
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

func (server *NativeServer) broadcastGameState() {
	ticker := time.NewTicker(TICKER_TIME)
	defer ticker.Stop()

	for range ticker.C {
		server.gameState.UpdateState()
		gameStateBytes := server.gameState.GetGameState()
		server.broadcast <- gameStateBytes
	}
}

func (server *NativeServer) handleWebSocket(conn *websocket.Conn) {
	server.addClient <- conn
	defer func() { server.removeClient <- conn }()
	for {
		var message string
		err := websocket.Message.Receive(conn, &message)
		if err != nil {
			log.Println("Error reading message from client:", err)
			break
		}
		// fmt.Println(message)
		server.gameState.MovePlayer(conn, []byte(message))
	}
}

// {"x":3.6,"z":19.01}
