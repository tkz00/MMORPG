package character

import (
	"unnamed-mmo/backend/pkg/game/stats"
	"unnamed-mmo/backend/pkg/utils"
)

type NPC struct {
	health	  		stats.Health
	position  		utils.Vector2
	to        		utils.Vector2
	direction 		utils.Vector2
	onDeath     	func(*NPC) // Callback to notify when the NPC dies
}

func (n *NPC) Die() {
	if n.onDeath != nil {
		n.onDeath(n)
	}
	// Trigger any other death logic here
}

func (n *NPC) RegisterOnDeathCallback(callback func(*NPC)) {
	n.onDeath = callback
}
