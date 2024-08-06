package types

import (
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

func (gs *GameState) AddPlayer(conn *websocket.Conn) Player {
	id := uuid.New()
	playerId := id.String()
	abilities := map[string]*Ability{
		"1": NewAbility("1", "heal", 10, 3000),
		"0": NewAbility("0", "projectile", 5, 2000),
	}
	player := CreatePlayer(playerId, 0, 0, abilities)
	gs.playerIds[conn] = playerId
	gs.players[playerId] = player

	return *player
}

func (gs *GameState) DeletePlayer(conn *websocket.Conn) {
	playerId := gs.playerIds[conn]
	delete(gs.players, playerId)
	delete(gs.playerIds, conn)
}

func (gs GameState) GetPlayerCount() int {
	return len(gs.players)
}

func (gs GameState) GetPlayers() []Player {
	// WTF is this shit scoob????
	playersSlice := make([]Player, 0, len(gs.players))
    for _, player := range gs.players {
        playersSlice = append(playersSlice, *player)
    }
    return playersSlice
}

func (gs GameState) GetProjectiles() []Projectile {
	projectilesSlice := make([]Projectile, 0, len(gs.projectiles))
    for _, projectile := range gs.projectiles {
        projectilesSlice = append(projectilesSlice, *projectile)
    }
    return projectilesSlice
}

func (gs GameState) MovePlayer(conn *websocket.Conn, position Position) {
	playerId := gs.playerIds[conn]
	moveAction := &MoveAction{
		targetPosition: position,
	}
	gs.players[playerId].EnqueueAction(moveAction)
}

func (gs GameState) UpdateState(deltaTime float64) {
	gs.updatePlayers(deltaTime)
	gs.updateProjectiles(deltaTime)
}

func (gs *GameState) updatePlayers(deltaTime float64) {
	for _, player := range gs.players {
		player.ExecuteNextAction(gs)
		if player.IsMoving() {
			player.UpdatePosition(deltaTime)
		}
	}
}

func (gs GameState) updateProjectiles(deltaTime float64) {
	for key, projectile := range gs.projectiles {
		if projectile.state == Hit {
			delete(gs.projectiles, key)
			continue
		}
		
		if projectile.state == Active && gs.checkCollision(*projectile) {
			projectile.state = Hit;
			continue
		}
		
		isAtMaxRange := projectile.UpdatePosition(deltaTime)

		if isAtMaxRange {
			delete(gs.projectiles, key)
		}
	}
}

func (gs GameState) checkCollision(projectile Projectile) bool {
	playerGotHit := false
	for _, player := range gs.players {
		if player.id != projectile.caster && gs.AreColliding(*player, projectile) {
			player.DealDamage(projectile.damage)
			playerGotHit = true
		}
	}

	return playerGotHit
}

func (gs *GameState) CastAbility(conn *websocket.Conn, abilityInfo AbilityInfo) {
	casterId := gs.playerIds[conn]
	caster := gs.players[casterId]
	caster.CastAbility(gs, abilityInfo)
}

func (gs GameState) AreColliding(player Player, projectile Projectile) bool {
	playerPosition := player.GetPosition()
	projectilePosition := projectile.GetPosition()
	distance := playerPosition.Distance(projectilePosition)
	return distance < (player.GetRadius() + projectile.GetRadius())
}

func (gs *GameState) ResetPlayersState() {
	for _, player := range gs.players {
		player.executingAction = Idle
	}
}
