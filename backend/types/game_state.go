package types

import (
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

func (gs *GameState) AddPlayer(conn *websocket.Conn) string {
	id := uuid.New()
	playerId := id.String()
	gs.playerIds[conn] = playerId
	gs.players[playerId] = CreatePlayer(0, 0)

	return playerId
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
	gs.players[playerId].MoveTowards(position)
}

func (gs GameState) UpdateState() {
	for _, player := range gs.players {
		if player.IsMoving() {
			player.UpdatePosition()
		}
	}
}

func (gs GameState) GetGameState() GameDTO {
	var players []PlayerDTO

	for playerId, player := range gs.players {
		players = append(players, player.ToDTO(playerId))
	}

	gameDTO := GameDTO{
		Players: players,
	}

	return gameDTO
}
