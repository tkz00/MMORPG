package character

import (
	"tkz00/backend/pkg/utils"
)

type CharacterAction interface {
	Execute(player *Character) error
	IsComplete() bool
}

type MoveAction struct {
	TargetPosition utils.Vector2
	isComplete     bool
}

func (a *MoveAction) Execute(character *Character) error {
	if !a.isComplete {
		character.MoveTowards(a.TargetPosition)
		a.isComplete = character.position == a.TargetPosition // Adjust this check based on your movement logic
	}
	return nil
}

func (a *MoveAction) IsComplete() bool {
	return a.isComplete
}

type CastAbilityAction struct {
	ability    Ability
	params     AbilityParameters
	isComplete bool
}

func (a *CastAbilityAction) Execute(caster *Character) error {
	if a.ability.targeting == Target {
		targetPosition, err := a.params.GetCastingCoordinates()
		if err != nil {
			return err
		}
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
	GetTargetCoordinates() (utils.Vector2, error)
	GetCastingCoordinates() (utils.Vector2, error)
}

type CoordinateAbilityParams struct {
	Target utils.Vector2
}

func (p CoordinateAbilityParams) GetTargetCoordinates() (utils.Vector2, error) {
	return p.Target, nil
}

func (p CoordinateAbilityParams) GetCastingCoordinates() (utils.Vector2, error) {
	return p.Target, nil
}

type TargetIdAbilityParams struct {
	TargetId                   string
	TargetCoordinatesCallback  func(targetId string) (utils.Vector2, error)
	CastingCoordinatesCallback func(targetId string) (utils.Vector2, error)
}

func (p TargetIdAbilityParams) GetTargetCoordinates() (utils.Vector2, error) {
	return p.TargetCoordinatesCallback(p.TargetId)
}

func (p TargetIdAbilityParams) GetCastingCoordinates() (utils.Vector2, error) {
	return p.CastingCoordinatesCallback(p.TargetId)
}
