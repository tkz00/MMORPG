package entities

import (
	"fmt"
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
		for _, mechanic := range action.ability.mechanics {
			if handler, exists := mechanicHandlers[mechanic.MechanicType]; exists {
				if err := handler(caster, action.castParameters[Target].(string), gs, mechanic.Params); err != nil {
					fmt.Println(err)
				} else {
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
