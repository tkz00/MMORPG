package types

import (
	"encoding/json"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

type GameState struct {
	playerIds map[*websocket.Conn]string
	players   map[string]*Player
}

func StartGameState() GameState {
	return GameState{
		playerIds: make(map[*websocket.Conn]string),
		players:   make(map[string]*Player),
	}
}

func (gs *GameState) AddPlayer(conn *websocket.Conn) {
	id := uuid.New()
	playerId := id.String()
	gs.playerIds[conn] = playerId
	gs.players[playerId] = CreatePlayer(20.0, 10.0)
}

func (gs *GameState) DeletePlayer(conn *websocket.Conn) {
	playerId := gs.playerIds[conn]
	delete(gs.players, playerId)
	delete(gs.playerIds, conn)
}

func (gs GameState) GetPlayerCount() int {
	return len(gs.players)
}

func (gs GameState) MovePlayer(conn *websocket.Conn, positionMsg []byte) {
	position := CreatePosition(positionMsg)
	playerId := gs.playerIds[conn]
	gs.players[playerId].SetPosition(position)
}

func (gs GameState) GetGameState() string {
	gameDTO := GameDTO{
		Players: make([]PlayerDTO, 0),
	}

	for playerId, player := range(gs.players) {
		gameDTO.Players = append(gameDTO.Players, player.ToDTO(playerId))
	}

	gameStateJSON, _ := json.Marshal(gameDTO)

	return string(gameStateJSON)
}
