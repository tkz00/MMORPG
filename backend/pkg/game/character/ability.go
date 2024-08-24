package character

import (
	"unnamed-mmo/backend/pkg/utils"
)

type Targeting int

const (
	Target Targeting = iota
	NoTarget
	Coordinates
)

type Ability struct {
	id             string
	name           string
	rangeValue     float64
	cooldown       int64
	targeting      Targeting
	characterState Action
	Mechanic       func(caster Character, params AbilityParameters)
}

func NewAbility(
	id string,
	name string,
	rangeValue float64,
	cooldown int64,
	targeting Targeting,
	characterState Action,
	mechanic func(caster Character, params AbilityParameters)) *Ability {
	return &Ability{
		id:             id,
		name:           name,
		rangeValue:     rangeValue,
		cooldown:       cooldown,
		targeting:      targeting,
		characterState: characterState,
		Mechanic:       mechanic,
	}
}

func (ability Ability) Id() string {
	return ability.id
}

func (ability Ability) Name() string {
	return ability.name
}

func (ability Ability) Range() float64 {
	return ability.rangeValue
}

func (ability Ability) Cooldown() int64 {
	return ability.cooldown
}

func (ability Ability) Targeting() Targeting {
	return ability.targeting
}

func (ability Ability) Cast(caster Character, params AbilityParameters) {
	ability.Mechanic(caster, params)
}

func (ability Ability) CreateAction(abilityInfo AbilityInfo, targetPositionCallback func(targetId string) utils.Vector2) CastAbilityAction {
	var abilityParams AbilityParameters
	switch ability.targeting {
	case Target:
		targetId, _ := abilityInfo.GetTargetId()
		abilityParams = TargetIdAbilityParams{
			TargetId:               targetId,
			TargetPositionCallback: targetPositionCallback,
		}
	case Coordinates:
		targetPosition, _ := abilityInfo.GetTargetPosition()
		abilityParams = CoordinateAbilityParams{
			Target: targetPosition,
		}
	}
	return CastAbilityAction{
		ability: ability,
		params:  abilityParams,
	}
}
