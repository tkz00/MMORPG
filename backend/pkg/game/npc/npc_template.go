package npc

import (
	"unnamed-mmo/backend/pkg/game/stats"
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

func (npcTemplate NPCTemplate) NewNPC(position utils.Vector2) NPC {
	return NPC{
		health: stats.NewHealth(npcTemplate.startingHealth),
		position: position,
	}
}
