package character

import (
	"unnamed-mmo/backend/pkg/utils"
)

type NPCTemplate struct {
	id				string
	name			string
    startingHealth	int
}

func NewNPCTemplate(id string, name string, startingHealth int) NPCTemplate{
	return NPCTemplate{
		id: id,
		name: name,
		startingHealth: startingHealth,
	}
}

func (npcTemplate NPCTemplate) NewNPC(id string, position utils.Vector2) *Character {
	x, z := position.GetPosition()
	return CreateCharacter(id, x, z, map[string]*Ability{})
}
