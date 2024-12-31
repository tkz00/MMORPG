package gameplay

import (
	"tkz00/backend/pkg/game/entities"
	"tkz00/backend/pkg/game/repository"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

func StartGameState() *entities.GameState {
	gamestate := entities.StartGameState()
	skeletonEnemyAbilities := repository.GetSkeletonEnemyAbilities(gamestate)
	gamestate.SetUpSkeletonEnemies(skeletonEnemyAbilities)

	mapObstacleColliders := repository.GetObstacleColliders()
	gamestate.SetUpObstacles(mapObstacleColliders)

	return gamestate
}

func UpdateState(gs *entities.GameState, deltaTime float64) {
	updatePlayers(gs, deltaTime)
	updateProjectiles(gs, deltaTime)
	updateAreaEffects(gs, deltaTime)
	updateSpawners(gs)
	updateNpcs(gs, deltaTime)
	gs.UpdateMechanics(deltaTime)
}

func updatePlayers(gs *entities.GameState, deltaTime float64) {
	for _, player := range gs.Players() {
		if !player.IsAlive() {
			player.ClearActionsQueue()
			return
		}

		player.ExecuteNextAction(gs)
		if player.IsMoving() {
			if !player.UpdatePosition(deltaTime, gs.GetObstacleColliders()) {
				player.ClearActionsQueue()
			}
		}
	}
}

func updateProjectiles(gs *entities.GameState, deltaTime float64) {
	for key, projectile := range gs.Projectiles() {
		if projectile.State() == entities.Hit {
			gs.RemoveProjectile(key)
			continue
		}

		if projectile.State() == entities.Active && checkCollision(gs, *projectile) {
			projectile.SetState(entities.Hit)
			continue
		}

		isAtMaxRange := projectile.UpdatePosition(deltaTime)

		if isAtMaxRange {
			gs.RemoveProjectile(key)
		}
	}
}

func updateAreaEffects(gs *entities.GameState, deltaTime float64) {
	var AoEsToDelete []string
	for AoEID, AoE := range gs.AreaEffects() {
		AoE.Tick(gs, deltaTime)

		if AoE.RemainingDurationMs() <= 0 {
			AoEsToDelete = append(AoEsToDelete, AoEID)
		}
	}
	for _, id := range AoEsToDelete {
		delete(gs.AreaEffects(), id)
	}
}

func updateSpawners(gs *entities.GameState) {
	for _, spawner := range gs.Spawners() {
		newNPCs := spawner.GetNewNPCs()
		for _, newNPC := range newNPCs {
			gs.NPCs()[newNPC.GetId()] = newNPC
		}
	}
}

func updateNpcs(gs *entities.GameState, deltaTime float64) {
	var deadNpcs []string

	for id, npc := range gs.NPCs() {
		if npc.IsAlive() {
			npc.ExecuteNextAction(gs)
			npc.UpdateBehaviour()
			if npc.IsMoving() {
				if !npc.UpdatePosition(deltaTime, gs.GetObstacleColliders()) {
					npc.ClearActionsQueue()
				}
			}
		} else {
			npc.Remove()
			deadNpcs = append(deadNpcs, id)
		}
	}

	for _, id := range deadNpcs {
		delete(gs.NPCs(), id)
	}
}

// AreColliding surely must go in a separate collisions module, I don't know about checkCollisions, but from what I'm seeing it must be refactored
func checkCollision(gs *entities.GameState, projectile entities.Projectile) bool {
	for _, player := range gs.Players() {
		if !player.IsAlive() || player.GetId() == projectile.CasterId() {
			continue
		}
		if AreColliding(*player, projectile) {
			// execute projectile mechanics
			projectile.Hit(player, gs)
			return true
		}
	}

	for _, npc := range gs.NPCs() {
		if !npc.IsAlive() || npc.GetId() == projectile.CasterId() {
			continue
		}
		if AreColliding(*npc.Character, projectile) {
			// execute projectile mechanics
			projectile.Hit(npc.Character, gs)
			return true
		}
	}

	return false
}

func AreColliding(player entities.Character, projectile entities.Projectile) bool {
	playerPosition := player.GetPosition()
	projectilePosition := projectile.GetPosition()
	distance := playerPosition.Distance(projectilePosition)
	return distance < (player.GetRadius() + projectile.GetRadius())
}

// Should this be here? Where should this be?
func AddPlayer(gs *entities.GameState, conn *websocket.Conn) entities.Character {
	id := uuid.New()
	playerId := id.String()
	abilities := repository.GetPlayerAbilities(gs)
	player := entities.CreateCharacter(playerId, 0, 0, abilities)
	gs.AddPlayer(conn, player)
	return *player
}
