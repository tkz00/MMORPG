package spawner

import (
	"time"
	"unnamed-mmo/backend/pkg/game/npc"
	"unnamed-mmo/backend/pkg/utils"
)

type Spawner struct {
	position	utils.Vector2
	radius		float64
	rate		time.Duration
	npcTemplate	npc.NPCTemplate
	lastSpawned	time.Time
}

func NewSpawner(position utils.Vector2, radius float64, rate time.Duration, npcTemplate npc.NPCTemplate) *Spawner {
	return &Spawner{
		position: position,
		radius: radius,
		rate: rate,
		npcTemplate: npcTemplate,
		lastSpawned: time.Now(),
	}
}

func (spawner *Spawner)GetNewNPCs() []npc.NPC {
	elapsedTime := time.Since(spawner.lastSpawned)
	spawnIntervals := int(elapsedTime / spawner.rate)
	if spawnIntervals <= 0 {
		return nil
	}
	spawner.lastSpawned = spawner.lastSpawned.Add(spawner.rate * time.Duration(spawnIntervals))
	
	npcs := make([]npc.NPC, spawnIntervals)
	for i := range npcs {
		npcPosition := utils.RandomCoordinatesInRadius(spawner.position, spawner.radius)
		npcs[i] = spawner.npcTemplate.NewNPC(npcPosition)
	}

	return npcs
}
