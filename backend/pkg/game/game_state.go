package game

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"

	"unnamed-mmo/backend/pkg/game/character"
	"unnamed-mmo/backend/pkg/game/spawner"
	"unnamed-mmo/backend/pkg/utils"
)

type GameState struct {
	playerIds 	map[*websocket.Conn]string
	players   	map[string]*character.Player
	projectiles map[string]*Projectile
	spawners 	map[string]*spawner.Spawner
	npcs		map[string]character.NPC
}

func StartGameState() GameState {
	skeletonNPCTemplate := character.NewNPCTemplate("0", "skeleton", 25)

	spawners := map[string]*spawner.Spawner{
		"skeleton_spawner_0": spawner.NewSpawner(
			*utils.NewVector2(0, 15),
			2.5,
			4 * time.Second,
			skeletonNPCTemplate,
		),
	}

	return GameState{
		playerIds: 		make(map[*websocket.Conn]string),
		players:   		make(map[string]*character.Player),
		projectiles: 	make(map[string]*Projectile),
		spawners: 		spawners,
		npcs: 			make(map[string]character.NPC),
	}
}

func (gs *GameState) AddPlayer(conn *websocket.Conn) character.Player {
	id := uuid.New()
	playerId := id.String()
	abilities := map[string]*character.Ability{
		"0": character.NewAbility("0", "projectile", 5, 2000, character.Coordinates,
		func(caster character.Player, params character.AbilityParameters) {
			projectileId := uuid.NewString()
			gs.projectiles[projectileId] = CreateProjectile(uuid.NewString(), caster.GetPosition(), params.(character.CoordinateAbilityParams).Target, 5, caster.GetId())
		}),
		"1": character.NewAbility("1", "heal", 7, 3000, character.Target,
		func(caster character.Player, params character.AbilityParameters) {
			target := gs.players[params.(character.TargetIdAbilityParams).TargetId]
			target.HealthVariation(-10)
		}),
	}
	player := character.CreatePlayer(playerId, 0, 0, abilities)
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

func (gs GameState) GetPlayers() []character.Player {
	// WTF is this shit scoob????
	playersSlice := make([]character.Player, 0, len(gs.players))
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

func (gs GameState) MovePlayer(conn *websocket.Conn, position utils.Vector2) {
	playerId := gs.playerIds[conn]
	moveAction := &character.MoveAction{
		TargetPosition: position,
	}
	gs.players[playerId].ClearActionsQueue()
	gs.players[playerId].EnqueueAction(moveAction)
}

func (gs GameState) UpdateState(deltaTime float64) {
	gs.updatePlayers(deltaTime)
	gs.updateProjectiles(deltaTime)
	gs.updateSpawners()
}

func (gs *GameState) updatePlayers(deltaTime float64) {
	for _, player := range gs.players {
		player.ExecuteNextAction()
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
		if player.GetId() != projectile.caster && gs.AreColliding(*player, projectile) {
			player.HealthVariation(projectile.damage)
			playerGotHit = true
		}
	}

	return playerGotHit
}

func (gs *GameState) EnqueueAbilityCast(conn *websocket.Conn, abilityInfo character.AbilityInfo) {
	casterId := gs.playerIds[conn]
	caster := gs.players[casterId]
	ability := caster.GetAbilities()[abilityInfo.GetId()]
	abilityAction := ability.CreateAction(abilityInfo, 
	func(targetId string) utils.Vector2 {
		if caster.GetPosition().Distance(gs.players[targetId].GetPosition()) > ability.Range() {
			return caster.ClosestPositionInRange(gs.players[targetId].GetPosition(), ability.Range())
		} else {
			return caster.GetPosition()
		}
	})
	caster.EnqueueAbilityCast(abilityAction)
}

func (gs GameState) AreColliding(player character.Player, projectile Projectile) bool {
	playerPosition := player.GetPosition()
	projectilePosition := projectile.GetPosition()
	distance := playerPosition.Distance(projectilePosition)
	return distance < (player.GetRadius() + projectile.GetRadius())
}

// How do I solve this?
// func (gs *GameState) ResetPlayersState() {
// 	for _, player := range gs.players {
// 		player.executingAction = character.Idle
// 	}
// }

func (gs *GameState) updateSpawners() {
	for _, spawner := range gs.spawners {
		newNPCs := spawner.GetNewNPCs()
		for _, newNPC := range newNPCs {
			gs.npcs[uuid.NewString()] = newNPC
		}
	}
}
