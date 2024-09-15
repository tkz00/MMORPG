package character

import (
	"time"
	"tkz00/backend/pkg/utils"
)

type NPCTemplate struct {
	id             string
	name           string
	startingHealth int
	abilities      map[string]*Ability
	aggroRange     float64
}

func NewNPCTemplate(id string, name string, startingHealth int, aggroRange float64, abilities map[string]*Ability) NPCTemplate {
	return NPCTemplate{
		id:             id,
		name:           name,
		startingHealth: startingHealth,
		aggroRange:     aggroRange,
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
		aggroRange:        npcTemplate.aggroRange,
	}
}
