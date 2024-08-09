package npc

import (
	"unnamed-mmo/backend/pkg/game/stats"
	"unnamed-mmo/backend/pkg/utils"
)

type NPC struct {
	health	  		stats.Health
	position  		utils.Vector2
	to        		utils.Vector2
	direction 		utils.Vector2
}
