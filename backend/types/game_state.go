package types

import (
	"golang.org/x/net/websocket"
)

type GameState struct {
	players map[*websocket.Conn]Player
}

func StartGameState() GameState {
	return GameState{
		players: make(map[*websocket.Conn]Player),
	}
}

func (gs *GameState) AddPlayer(conn *websocket.Conn) {
	gs.players[conn] = CreatePlayer(20.0, 10.0)
}

func (gs *GameState) DeletePlayer(conn *websocket.Conn) {
	delete(gs.players, conn)
}
func (gs GameState) GetPlayerCount() int {
	return len(gs.players)
}
func (gs GameState) MovePlayer(conn *websocket.Conn, position Position) {
	gs.players[conn].SetPosition(position)
}
