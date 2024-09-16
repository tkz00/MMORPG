package gameplay

import (
	"tkz00/backend/pkg/game/entities"
)

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

func checkCollision(gs *entities.GameState, projectile entities.Projectile) bool {
	projectileHit := false
	for _, player := range gs.Players() {
		if player.GetId() != projectile.Caster() && gs.AreColliding(*player, projectile) {
			player.HealthVariation(-projectile.GetDamage())
			projectileHit = true
		}
	}

	for _, npc := range gs.NPCs() {
		if npc.GetId() != projectile.Caster() && gs.AreColliding(*npc.Character, projectile) {
			npc.HealthVariation(-projectile.GetDamage())
			if caster, err := gs.GetCharacterById(projectile.Caster()); err == nil {
				npc.BecomeAggressive(caster)
			}
			projectileHit = true
		}
	}

	return projectileHit
}
