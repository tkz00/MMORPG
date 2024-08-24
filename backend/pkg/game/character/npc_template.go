package character

import (
	"time"
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

func (npcTemplate NPCTemplate) NewNPC(id string, position utils.Vector2, spawnerPosition utils.Vector2) *Npc {
	x, z := position.GetPosition()
	return &Npc{
		Character: CreateCharacter(id, x, z, map[string]*Ability{}),
		lastActionDecided: time.Now(),
		spawnerPosition: spawnerPosition,
	}
}
