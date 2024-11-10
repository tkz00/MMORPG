package entities

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
