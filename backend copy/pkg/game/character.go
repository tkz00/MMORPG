package game

import "tkz00/backend/pkg/utils"

type Character struct {
	position utils.Vector2
}

func (character Character) Position() utils.Vector2 {
	return character.position
}
