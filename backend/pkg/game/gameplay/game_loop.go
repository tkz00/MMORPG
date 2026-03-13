package gameplay

import (
	"backend/pkg/game/entities"
	"backend/pkg/game/repository"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

func StartGameState() *entities.GameState {
	gamestate := entities.StartGameState()
	skeletonEnemyAbilities := repository.GetSkeletonEnemyAbilities(gamestate)
	gamestate.SetUpSkeletonEnemies(skeletonEnemyAbilities)

	mapObstacleColliders := repository.GetObstacleColliders()
	gamestate.SetUpObstacles(mapObstacleColliders)

	go saveGameState(gamestate)

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
			continue
		}

		// Remove expired buffs before processing actions
		player.RemoveExpiredBuffs()

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
			newNPC.TriggerEffects(entities.EffectPassive, nil, gs)
		}
	}
}

func updateNpcs(gs *entities.GameState, deltaTime float64) {
	var deadNpcs []string

	for id, npc := range gs.NPCs() {
		if npc.IsAlive() {
			// Remove expired buffs before processing actions
			npc.RemoveExpiredBuffs()

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
func AddPlayer(gs *entities.GameState, conn *websocket.Conn, characterName string) entities.Character {
	var playerId string
	var x, z float64
	var maxHealth, currentHealth int
	var stats map[string]int64
	var abilities map[string]*entities.Ability

	var player *entities.Character

	if player, _ = repository.GetCharacterByName(characterName); player == nil {
		id := uuid.New()
		playerId = id.String()
		x, z = 0, 0
		maxHealth, currentHealth = entities.BASE_MAX_HEALTH, entities.BASE_MAX_HEALTH
		stats = map[string]int64{"damage": 10, "defense": 5}
		abilities = repository.GetPlayersInitialAbilities()
		player = entities.CreateCharacter(playerId, characterName, x, z, maxHealth, currentHealth, stats, abilities, nil, nil, nil)

		player.TriggerEffects(entities.EffectPassive, nil, gs)

		go func(player *entities.Character) {
			if err := repository.SaveCharacter(player); err != nil {
				fmt.Printf("Error saving new character to repository: %v\n", err)
			}
		}(player)
	}

	if player == nil {
		fmt.Println("player == nil")
	}
	gs.AddPlayer(conn, player)

	return *player
}

func saveGameState(gs *entities.GameState) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		for _, player := range gs.Players() {
			if err := repository.SaveCharacter(player); err != nil {
				fmt.Printf("Error saving new character to repository: %v\n", err)
			}
		}
	}
}
