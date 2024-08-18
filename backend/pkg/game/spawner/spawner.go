package spawner

import (
	"time"
	"unnamed-mmo/backend/pkg/game/character"
	"unnamed-mmo/backend/pkg/utils"
)

type Spawner struct {
	position	utils.Vector2
	radius		float64
	rate		time.Duration
	npcTemplate	character.NPCTemplate
	lastSpawned	time.Time
	activeNPCs  []character.NPC
	maxNPCs     int
}

func NewSpawner(position utils.Vector2, radius float64, rate time.Duration, npcTemplate character.NPCTemplate) *Spawner {
	return &Spawner{
		position: position,
		radius: radius,
		rate: rate,
		npcTemplate: npcTemplate,
		lastSpawned: time.Now(),
		activeNPCs:  []character.NPC{},
		maxNPCs:     3,
	}
}

func (spawner *Spawner)GetNewNPCs() []character.NPC {
	elapsedTime := time.Since(spawner.lastSpawned)
	spawnIntervals := int(elapsedTime / spawner.rate)
	if spawnIntervals <= 0 {
		return nil
	}
	npcsToSpawn := spawnIntervals
	if len(spawner.activeNPCs) + npcsToSpawn > spawner.maxNPCs {
		npcsToSpawn = spawner.maxNPCs - len(spawner.activeNPCs)
	}
	if npcsToSpawn <= 0 {
		return nil
	}
	spawner.lastSpawned = spawner.lastSpawned.Add(spawner.rate * time.Duration(spawnIntervals))
	
	npcs := make([]character.NPC, npcsToSpawn)
	for i := range npcs {
		npcPosition := utils.RandomCoordinatesInRadius(spawner.position, spawner.radius)
		npcs[i] = spawner.npcTemplate.NewNPC(npcPosition)
		npcs[i].RegisterOnDeathCallback(spawner.HandleNPCDeath)
	}

	spawner.activeNPCs = append(spawner.activeNPCs, npcs...)

	return npcs
}

// HandleNPCDeath removes the NPC from the activeNPCs list when it dies.
func (spawner *Spawner) HandleNPCDeath(npc *character.NPC) {
	for i, activeNPC := range spawner.activeNPCs {
		if &activeNPC == npc {
			// Remove the NPC from the slice
			spawner.activeNPCs = append(spawner.activeNPCs[:i], spawner.activeNPCs[i+1:]...)
			break
		}
	}
}

