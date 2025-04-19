package entities

import (
	"backend/pkg/utils"
	"fmt"
	"time"
)

type CharacterAction interface {
	Execute(player *Character, gs *GameState) error
	IsComplete() bool
}

type MoveAction struct {
	TargetPosition utils.Vector2
	isComplete     bool
}

func (a *MoveAction) Execute(character *Character, _ *GameState) error {
	if !a.isComplete {
		character.MoveTowards(a.TargetPosition)
		a.isComplete = character.position == a.TargetPosition // Adjust this check based on your movement logic
	}
	return nil
}

func (a *MoveAction) IsComplete() bool {
	return a.isComplete
}

type AbilityCastAction struct {
	ability        Ability
	castParameters map[Targeting]interface{}
	isComplete     bool
}

func (action *AbilityCastAction) Execute(caster *Character, gs *GameState) error {
	switch action.ability.targeting {
	case Target:
		targetId := action.castParameters[Target].(string)
		target, _ := gs.GetCharacterById(targetId)
		if MoveIfNotInRange(caster, target.position, action.ability.rangeValue, gs) {
			return nil
		}
		for _, mechanic := range action.ability.mechanics {
			if handler, exists := mechanicHandlers[mechanic.MechanicType]; exists {
				mechanic.Params["target_id"] = targetId
				resolveParameters(
					mechanic,
					caster.id,
					gs,
				)
				if err := handler(caster, gs, mechanic.Params); err != nil {
					fmt.Println(err)
				} else {
					caster.lastUsed[action.ability.id] = time.Now()
					action.isComplete = true

					target, _ := gs.GetCharacterById(action.castParameters[Target].(string))
					switch action.ability.characterState {
					case Attacking:
						if caster.id == target.id {
							caster.executingAction = ExecutingAction{Attacking, caster.executingAction.direction}
						} else {
							normalizedCastAbilityVector := utils.Normalize(caster.position, target.position)
							caster.executingAction = ExecutingAction{Attacking, normalizedCastAbilityVector}
						}
					case CastingHeal:
						if caster.id == target.id {
							caster.executingAction = ExecutingAction{CastingHeal, caster.executingAction.direction}
						} else {
							normalizedCastAbilityVector := utils.Normalize(caster.position, target.position)
							caster.executingAction = ExecutingAction{CastingHeal, normalizedCastAbilityVector}
						}
					}
				}
			} else {
				fmt.Printf("no handler found for effect type: %s\n", mechanic.MechanicType)
			}
		}
	case Coordinates:
		targetCoordinates := action.castParameters[Coordinates].(utils.Vector2)
		for _, mechanic := range action.ability.mechanics {
			if mechanic.MechanicType == "create_AoE" {
				if MoveIfNotInRange(caster, targetCoordinates, action.ability.rangeValue, gs) {
					return nil
				}
			}
			if handler, exists := mechanicHandlers[mechanic.MechanicType]; exists {
				mechanic.Params["target_coordinates"] = targetCoordinates
				mechanic.Params["range"] = action.ability.rangeValue
				mechanic.Params["initial_coordinates"] = caster.position
				if err := handler(caster, gs, mechanic.Params); err != nil {
					fmt.Println(err)
				} else {
					caster.lastUsed[action.ability.id] = time.Now()
					action.isComplete = true
					normalizedCastAbilityVector := utils.Normalize(caster.position, targetCoordinates)
					switch action.ability.characterState {
					case Attacking:
						caster.executingAction = ExecutingAction{Attacking, normalizedCastAbilityVector}
					case CastingHeal:
						caster.executingAction = ExecutingAction{CastingHeal, normalizedCastAbilityVector}
					}
				}
			} else {
				fmt.Printf("no handler found for effect type: %s\n", mechanic.MechanicType)
			}
		}
	}
	return nil
}

func (action *AbilityCastAction) IsComplete() bool {
	return action.isComplete
}

func MoveIfNotInRange(
	caster *Character,
	targetPosition utils.Vector2,
	rangeValue float64,
	gs *GameState,
) bool {
	const epsilon = 1e-9
	if (caster.position.Distance(targetPosition) - rangeValue) > epsilon {
		moveAction := &MoveAction{
			TargetPosition: utils.ClosestPositionInRange(
				caster.position,
				targetPosition,
				rangeValue,
			),
		}
		caster.PrependAction(moveAction)
		return true
	}
	return false
}
