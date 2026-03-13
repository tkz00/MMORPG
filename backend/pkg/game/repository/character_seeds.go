package repository

import (
	"backend/pkg/game/entities"
	"fmt"

	"github.com/samber/lo"
)

func RunSeeds() {
	// move outside of character_seeds
	for _, seedEffect := range GetEffectSeeds() {
		if effect, _ := GetEffectById(seedEffect.Id()); effect == nil {
			if err := SaveEffect(&seedEffect); err != nil {
				fmt.Printf("Error saving new seed to repository: %v\n", err)
			}
		}
	}

	for _, seedCharacter := range GetCharacterSeeds() {
		if character, _ := GetCharacterByName(seedCharacter.GetName()); character == nil {
			if err := SaveCharacter(seedCharacter); err != nil {
				fmt.Printf("Error saving new character to repository: %v\n", err)
			}
		}
	}
}

// Delete this?
func GetEffectSeeds() []entities.Effect {
	return lo.Values(entities.ExistingEffects)
}

func GetCharacterSeeds() []*entities.Character {
	playersInitialAbilities := GetPlayersInitialAbilities()
	return []*entities.Character{
		entities.CreateCharacter(
			"barbarian",
			"barbarian",
			5.0,
			0.0,
			100,
			100,
			map[string]int64{"damage": 20, "defense": 5},
			playersInitialAbilities,
			map[string]int64{"helm_001": 1, "1": 5},
			map[entities.EquipmentType]*entities.Equipment{entities.Helmet: entities.GetEquipment("helm_001")},
			[]string{"spell_vampirism"}),
		entities.CreateCharacter(
			"paladin",
			"paladin",
			-5,
			0,
			100,
			100,
			map[string]int64{"damage": 8, "defense": 20},
			playersInitialAbilities,
			map[string]int64{"helm_002": 1, "chest_001": 1},
			map[entities.EquipmentType]*entities.Equipment{entities.Helmet: entities.GetEquipment("helm_002"), entities.Chest: entities.GetEquipment("chest_001")},
			[]string{"iron_will"}),
	}
}
