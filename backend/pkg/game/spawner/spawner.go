package spawner

import (
	"time"
	"unnamed-mmo/backend/pkg/game/npc"
	"unnamed-mmo/backend/pkg/utils"
)

type Spawner struct {
	position	utils.Vector2
	radius		float64
	rate		time.Time
	npc			npc.NPCTemplate
}
