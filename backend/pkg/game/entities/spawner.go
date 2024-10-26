package entities

import (
	"math/rand/v2"
	"time"
	"tkz00/backend/pkg/utils"

	"github.com/google/uuid"
)

type Spawner struct {
	position    utils.Vector2
	radius      float64
	rate        time.Duration
	npcTemplate NPCTemplate
	lastSpawned time.Time
	activeNPCs  []*Npc
	maxNPCs     int
}

func NewSpawner(position utils.Vector2, radius float64, rate time.Duration, npcTemplate NPCTemplate) *Spawner {
	return &Spawner{
		position:    position,
		radius:      radius,
		rate:        rate,
		npcTemplate: npcTemplate,
		lastSpawned: time.Now(),
		activeNPCs:  []*Npc{},
		maxNPCs:     3,
	}
}

func (spawner *Spawner) GetNewNPCs() []*Npc {
	elapsedTime := time.Since(spawner.lastSpawned)
	spawnIntervals := int(elapsedTime / spawner.rate)
	if spawnIntervals <= 0 {
		return nil
	}
	npcsToSpawn := spawnIntervals
	if len(spawner.activeNPCs)+npcsToSpawn > spawner.maxNPCs {
		npcsToSpawn = spawner.maxNPCs - len(spawner.activeNPCs)
	}
	if npcsToSpawn <= 0 {
		return nil
	}
	spawner.lastSpawned = spawner.lastSpawned.Add(spawner.rate * time.Duration(spawnIntervals))

	npcs := make([]*Npc, npcsToSpawn)
	for i := range npcs {
		newNpcId := uuid.NewString()
		npcPosition := utils.RandomCoordinatesInRadius(spawner.position, spawner.radius)
		npcs[i] = spawner.npcTemplate.NewNPC(newNpcId, npcPosition, spawner.position)
		npcs[i].SubscribeToRemoval(func() {
			spawner.HandleNPCDeath(npcs[i])
		})
		if rand.IntN(2) == 0 {
			npcs[i].AddItem(NewItem("1", "small health potion", Mechanic{MechanicType: "heal", Params: map[string]interface{}{"amount": 40}}), 1)
		} else {
			npcs[i].AddItem(NewItem("2", "leather"), 2)
		}
	}

	spawner.activeNPCs = append(spawner.activeNPCs, npcs...)

	return npcs
}

// HandleNPCDeath removes the NPC from the activeNPCs list when it dies.
func (spawner *Spawner) HandleNPCDeath(npc *Npc) {
	for i, activeNPC := range spawner.activeNPCs {
		if activeNPC == npc {
			// Remove the NPC from the slice
			spawner.activeNPCs = append(spawner.activeNPCs[:i], spawner.activeNPCs[i+1:]...)
			break
		}
	}
}
