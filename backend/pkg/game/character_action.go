package game

import (
	"unnamed-mmo/backend/pkg/utils"
)

type CharacterAction interface {
	Execute(player *Player) error
    IsComplete() bool
}

type MoveAction struct {
    targetPosition utils.Vector2
    isComplete     bool
}

func (a *MoveAction) Execute(player *Player) error {
    if !a.isComplete {
        player.MoveTowards(a.targetPosition)
        a.isComplete = player.position == a.targetPosition // Adjust this check based on your movement logic
    }
    return nil
}

func (a *MoveAction) IsComplete() bool {
    return a.isComplete
}

type CastAbilityAction struct {
	ability 	Ability
	params 		AbilityParameters
	isComplete 	bool
}

func (a *CastAbilityAction) Execute(player *Player) error {
	player.CastAbility(a.ability, a.params)
	a.isComplete = true
    return nil
}

func (a *CastAbilityAction) IsComplete() bool {
    return a.isComplete
}

type AbilityParameters interface { }

type CoordinateAbilityParams struct {
	AbilityParameters
	target utils.Vector2
}

type TargetIdAbilityParams struct {
	AbilityParameters
	targetId string
}
