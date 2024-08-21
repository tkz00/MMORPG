package spawner

import (
	"time"
	"unnamed-mmo/backend/pkg/game/character"
	"unnamed-mmo/backend/pkg/utils"

	"github.com/google/uuid"
)

type Spawner struct {
	position	utils.Vector2
	radius		float64
	rate		time.Duration
	npcTemplate	character.NPCTemplate
	lastSpawned	time.Time
	activeNPCs  []*character.Character
	maxNPCs     int
}

func NewSpawner(position utils.Vector2, radius float64, rate time.Duration, npcTemplate character.NPCTemplate) *Spawner {
	return &Spawner{
		position: position,
		radius: radius,
		rate: rate,
		npcTemplate: npcTemplate,
		lastSpawned: time.Now(),
		activeNPCs:  []*character.Character{},
		maxNPCs:     3,
	}
}

func (spawner *Spawner)GetNewNPCs() []*character.Character {
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
	
	npcs := make([]*character.Character, npcsToSpawn)
	for i := range npcs {
		newNpcId := uuid.NewString()
		npcPosition := utils.RandomCoordinatesInRadius(spawner.position, spawner.radius)
		npcs[i] = spawner.npcTemplate.NewNPC(newNpcId, npcPosition)
		// npcs[i].RegisterOnDeathCallback(spawner.HandleNPCDeath)
	}

	spawner.activeNPCs = append(spawner.activeNPCs, npcs...)

	return npcs
}

// HandleNPCDeath removes the NPC from the activeNPCs list when it dies.
func (spawner *Spawner) HandleNPCDeath(npc *character.Character) {
	for i, activeNPC := range spawner.activeNPCs {
		if activeNPC == npc {
			// Remove the NPC from the slice
			spawner.activeNPCs = append(spawner.activeNPCs[:i], spawner.activeNPCs[i+1:]...)
			break
		}
	}
}

