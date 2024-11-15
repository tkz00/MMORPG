package entities

import (
	"fmt"
	"time"
	"tkz00/backend/pkg/utils"
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
	if action.ability.targeting == Target {
		targetId := action.castParameters[Target].(string)
		if MoveIfNotInRange(caster, action, gs) {
			return nil
		}
		for _, mechanic := range action.ability.mechanics {
			if handler, exists := mechanicHandlers[mechanic.MechanicType]; exists {
				if err := handler(caster, targetId, gs, mechanic.Params); err != nil {
					fmt.Println(err)
				} else {
					caster.lastUsed[action.ability.id] = time.Now()
					action.isComplete = true
				}
			} else {
				fmt.Printf("no handler found for effect type: %s/n", mechanic.MechanicType)
			}
		}
	}
	return nil
}

func (action *AbilityCastAction) IsComplete() bool {
	return action.isComplete
}

func MoveIfNotInRange(caster *Character, action *AbilityCastAction, gs *GameState) bool {
	target, _ := gs.GetCharacterById(action.castParameters[Target].(string))
	const epsilon = 1e-9
	if (caster.position.Distance(target.position) - action.ability.rangeValue) > epsilon {
		moveAction := &MoveAction{
			TargetPosition: utils.ClosestPositionInRange(caster.position, target.position, action.ability.rangeValue),
		}
		caster.PrependAction(moveAction)
		return true
	}
	return false
}
