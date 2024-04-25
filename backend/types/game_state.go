package types

import (
	"unnamed-mmo/backend/utils"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

type GameState struct {
	playerIds map[*websocket.Conn]string
	players   map[string]*Player
	projectiles map[string]*Projectile
}

func StartGameState() GameState {
	return GameState{
		playerIds: make(map[*websocket.Conn]string),
		players:   make(map[string]*Player),
		projectiles: make(map[string]*Projectile),
	}
}

func (gs *GameState) AddPlayer(conn *websocket.Conn) PlayerDTO {
	id := uuid.New()
	playerId := id.String()
	player := CreatePlayer(0, 0, playerId)
	gs.playerIds[conn] = playerId
	gs.players[playerId] = player

	return *GetMapper().PlayerToDTO(*player)
}

func (gs *GameState) DeletePlayer(conn *websocket.Conn) {
	playerId := gs.playerIds[conn]
	delete(gs.players, playerId)
	delete(gs.playerIds, conn)
}

func (gs GameState) GetPlayerCount() int {
	return len(gs.players)
}

func (gs GameState) MovePlayer(conn *websocket.Conn, position Position) {
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

func (gs *GameState) CastAbility(conn *websocket.Conn, abilityDirection Position, abilityName string) {
	gs.projectiles[gs.playerIds[conn]] = CreateProjectile(abilityDirection)
}

func (gs GameState) GetGameState() GameStateDTO {
	return *GetMapper().GameStateToDTO(gs)
}

func (gs GameState) AreColliding(player1 Player, player2 Player) bool {
	position1 := player1.GetPosition()
	position2 := player2.GetPosition()

	diffX, diffZ := utils.GetDiff(position1.x, position1.z, position2.x, position2.z)
	playersDistance := utils.GetDistance(diffX, diffZ)

	return playersDistance < player1.GetRadius()+player2.GetRadius()
}
