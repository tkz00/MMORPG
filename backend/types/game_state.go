package types

import (
	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

type GameState struct {
	playerIds map[*websocket.Conn]string
	Players   map[string]*Player
}

func StartGameState() GameState {
	return GameState{
		playerIds: make(map[*websocket.Conn]string),
		Players:   make(map[string]*Player),
	}
}

func (gs *GameState) AddPlayer(conn *websocket.Conn) {
	id := uuid.New()
	playerId := id.String()
	gs.playerIds[conn] = playerId
	gs.Players[playerId] = CreatePlayer(20.0, 10.0)
}

func (gs *GameState) DeletePlayer(conn *websocket.Conn) {
	playerId := gs.playerIds[conn]
	delete(gs.Players, playerId)
	delete(gs.playerIds, conn)
}

func (gs GameState) GetPlayerCount() int {
	return len(gs.Players)
}

func (gs GameState) MovePlayer(conn *websocket.Conn, positionMsg []byte) {
	position := CreatePosition(positionMsg)
	playerId := gs.playerIds[conn]
	gs.Players[playerId].SetPosition(position)
}
