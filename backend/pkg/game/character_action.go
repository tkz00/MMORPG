package game

import (
	"time"
	"unnamed-mmo/backend/pkg/utils"

	"github.com/google/uuid"
)

type CharacterAction interface {
	Execute(player *Player, gameState *GameState) error
    IsComplete() bool
}

type MoveAction struct {
    targetPosition utils.Vector2
    isComplete     bool
}

func (a *MoveAction) Execute(player *Player, gameState *GameState) error {
    if !a.isComplete {
        player.MoveTowards(a.targetPosition)
        a.isComplete = player.position == a.targetPosition // Adjust this check based on your movement logic
    }
    return nil
}

func (a *MoveAction) IsComplete() bool {
    return a.isComplete
}

type HealAction struct {
    targetId   	string
	ability 	Ability
    isComplete 	bool
}

func (a *HealAction) Execute(player *Player, gameState *GameState) error {
    target := gameState.players[a.targetId]
	target.health.HealthVariation(-10)
	player.executingAction = CastingHeal
	// this code is duplicated and should be unified some way
	now := time.Now()
	player.lastUsed[a.ability.id] = now
	a.isComplete = true
    return nil
}

func (a *HealAction) IsComplete() bool {
    return a.isComplete
}

type ProjectileAction struct {
    targetPosition 	utils.Vector2
	ability 		Ability
    isComplete 		bool
}

func (a *ProjectileAction) Execute(player *Player, gameState *GameState) error {
	// is this ID necessary?
	projectileId := uuid.New().String()
	gameState.projectiles[projectileId] = CreateProjectile(projectileId, player.position, a.targetPosition, a.ability.rangeValue, player.id)
	player.executingAction = Attacking
	// this code is duplicated and should be unified some way
	now := time.Now()
	player.lastUsed[a.ability.id] = now
	a.isComplete = true
    return nil
}

func (a *ProjectileAction) IsComplete() bool {
    return a.isComplete
}
