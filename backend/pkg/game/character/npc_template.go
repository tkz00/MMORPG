package character

import (
	"time"
	"unnamed-mmo/backend/pkg/utils"
)

type NPCTemplate struct {
	id             string
	name           string
	startingHealth int
	abilities      map[string]*Ability
}

func NewNPCTemplate(id string, name string, startingHealth int, abilities map[string]*Ability) NPCTemplate {
	return NPCTemplate{
		id:             id,
		name:           name,
		startingHealth: startingHealth,
		abilities:      abilities,
	}
}

func (npcTemplate NPCTemplate) NewNPC(id string, position utils.Vector2, spawnerPosition utils.Vector2) *Npc {
	x, z := position.GetPosition()
	return &Npc{
		Character:         CreateCharacter(id, x, z, npcTemplate.abilities),
		lastActionDecided: time.Now(),
		spawnerPosition:   spawnerPosition,
		state:             Pacific,
	}
}
