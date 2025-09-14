package entities

import (
	"backend/pkg/utils"
	"time"
)

type NPCTemplate struct {
	id             string
	name           string
	startingHealth int
	stats          map[string]int64
	abilities      map[string]*Ability
	aggroRange     float64
}

func NewNPCTemplate(
	id string,
	name string,
	startingHealth int,
	stats map[string]int64,
	aggroRange float64,
	abilities map[string]*Ability,
) NPCTemplate {
	return NPCTemplate{
		id:             id,
		name:           name,
		startingHealth: startingHealth,
		stats:          stats,
		aggroRange:     aggroRange,
		abilities:      abilities,
	}
}

func (npcTemplate NPCTemplate) NewNPC(
	id string,
	position utils.Vector2,
	spawnerPosition utils.Vector2,
) *Npc {
	x, z := position.GetPosition()
	return &Npc{
		Character:         CreateCharacter(id, "TODO", x, z, npcTemplate.stats, npcTemplate.abilities),
		lastActionDecided: time.Now(),
		spawnerPosition:   spawnerPosition,
		state:             Pacific,
		aggroRange:        npcTemplate.aggroRange,
	}
}
