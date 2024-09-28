package gameplay

import (
	"fmt"
	"tkz00/backend/pkg/game/entities"
	"tkz00/backend/pkg/game/repository"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

func StartGameState() entities.GameState {
	gamestate := entities.StartGameState()
	skeletonEnemyAbilities := repository.GetSkeletonEnemyAbilities(gamestate)
	gamestate.SetUpSkeletonEnemies(skeletonEnemyAbilities)

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
		player.ExecuteNextAction()
		if player.IsMoving() {
			player.UpdatePosition(deltaTime)
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
			npc.Update(deltaTime)
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
		if player.GetId() != projectile.Caster() && AreColliding(*player, projectile) {
			player.HealthVariation(-projectile.GetDamage())
			projectileHit = true
		}
	}

	for _, npc := range gs.NPCs() {
		if npc.GetId() != projectile.Caster() && AreColliding(*npc.Character, projectile) {
			npc.HealthVariation(-projectile.GetDamage())
			if caster, err := gs.GetCharacterById(projectile.Caster()); err == nil {
				npc.BecomeAggressive(caster)
			}
			projectileHit = true
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
	abilities := map[string]*entities.Ability{
		"0": entities.NewAbility("0", "projectile", 5, 2000, entities.Coordinates, entities.Attacking,
			func(caster entities.Character, params entities.AbilityParameters) {
				projectile := entities.CreateProjectile(uuid.NewString(), caster.GetPosition(), params.(entities.CoordinateAbilityParams).Target, 5, caster.GetId())
				gs.AddProjectile(projectile)
			}),
		"1": entities.NewAbility("1", "heal", 7, 3000, entities.Target, entities.CastingHeal,
			func(caster entities.Character, params entities.AbilityParameters) {
				targetId := params.(entities.TargetIdAbilityParams).TargetId

				target, err := gs.GetCharacterById(targetId)
				if err != nil {
					fmt.Println(err)
					return
				}

				target.HealthVariation(10)
			}),
	}
	player := entities.CreateCharacter(playerId, 0, 0, abilities)
	gs.AddPlayer(conn, player)
	return *player
}
