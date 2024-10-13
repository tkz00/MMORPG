package gameplay

import (
	"tkz00/backend/pkg/game/entities"
	"tkz00/backend/pkg/game/repository"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

func StartGameState() entities.GameState {
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
	updateSpawners(gs)
	updateNpcs(gs, deltaTime)
}

func updatePlayers(gs *entities.GameState, deltaTime float64) {
	for _, player := range gs.Players() {
		if !player.IsAlive() {
			player.ClearActionsQueue()
			return
		}

		player.ExecuteNextAction()
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
			npc.ExecuteNextAction()
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
	projectileHit := false

	for _, player := range gs.Players() {
		if !player.IsAlive() || player.GetId() == projectile.CasterId() {
			continue
		}
		if AreColliding(*player, projectile) {
			player.HealthVariation(-projectile.GetDamage())
			projectileHit = true
			if !player.IsAlive() {
				gs.Players()[projectile.CasterId()].Loot(player.Inventory)
			}
		}
	}

	for _, npc := range gs.NPCs() {
		if !npc.IsAlive() || npc.GetId() == projectile.CasterId() {
			continue
		}
		if AreColliding(*npc.Character, projectile) {
			npc.HealthVariation(-projectile.GetDamage())
			if caster, err := gs.GetCharacterById(projectile.CasterId()); err == nil {
				npc.BecomeAggressive(caster)
			}
			projectileHit = true
			if !npc.IsAlive() {
				gs.Players()[projectile.CasterId()].Loot(npc.Inventory)
			}
		}
	}

	return projectileHit
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
