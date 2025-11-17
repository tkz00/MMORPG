package entities

import (
	"backend/pkg/utils"
	"math/rand/v2"
	"time"
)

type NPCTemplate struct {
	id             string
	name           string
	startingHealth int
	stats          map[string]int64
	abilities      map[string]*Ability
	aggroRange     float64
	effects        []string
}

func NewNPCTemplate(
	id string,
	name string,
	startingHealth int,
	stats map[string]int64,
	aggroRange float64,
	abilities map[string]*Ability,
	effects []string,
) NPCTemplate {
	return NPCTemplate{
		id:             id,
		name:           name,
		startingHealth: startingHealth,
		stats:          stats,
		aggroRange:     aggroRange,
		abilities:      abilities,
		effects:        effects,
	}
}

func (npcTemplate NPCTemplate) NewNPC(
	id string,
	position utils.Vector2,
	spawnerPosition utils.Vector2,
) *Npc {
	x, z := position.GetPosition()
	availableItems := []struct {
		id       string
		quantity int64
	}{
		{"0", 1},
		{"1", 2},
		{"helm_001", 1},
		{"helm_002", 1},
		{"chest_001", 1},
	}
	randomIndex := rand.IntN(len(availableItems))
	items := map[string]int64{
		availableItems[randomIndex].id: availableItems[randomIndex].quantity,
	}
	return &Npc{
		Character:         CreateCharacter(id, "TODO", x, z, npcTemplate.startingHealth, npcTemplate.startingHealth, npcTemplate.stats, npcTemplate.abilities, items, nil, npcTemplate.effects),
		lastActionDecided: time.Now(),
		spawnerPosition:   spawnerPosition,
		state:             Pacific,
		aggroRange:        npcTemplate.aggroRange,
	}
}
