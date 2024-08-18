package character

import (
	"unnamed-mmo/backend/pkg/utils"
)

type CharacterAction interface {
	Execute(player *Player) error
    IsComplete() bool
}

type MoveAction struct {
    TargetPosition utils.Vector2
    isComplete     bool
}

func (a *MoveAction) Execute(player *Player) error {
    if !a.isComplete {
        player.MoveTowards(a.TargetPosition)
        a.isComplete = player.position == a.TargetPosition // Adjust this check based on your movement logic
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

func (a *CastAbilityAction) Execute(caster *Player) error {
	if a.ability.targeting == Target {
		targetPosition := a.params.GetTargetPosition()
		if caster.GetPosition() != targetPosition {
			moveTarget := targetPosition
            caster.MoveTowards(moveTarget)
			return nil
        }
    }

	caster.CastAbility(a.ability, a.params)
	a.isComplete = true
    return nil
}

func (a *CastAbilityAction) IsComplete() bool {
    return a.isComplete
}

type AbilityParameters interface {
	GetTargetPosition() utils.Vector2
}

type CoordinateAbilityParams struct {
	Target utils.Vector2
}

func (p CoordinateAbilityParams) GetTargetPosition() utils.Vector2 {
    return p.Target
}

type TargetIdAbilityParams struct {
	TargetId string
	TargetPositionCallback func(targetId string) utils.Vector2
}

func (p TargetIdAbilityParams) GetTargetPosition() utils.Vector2 {
    return p.TargetPositionCallback(p.TargetId)
}
