package repository

import (
	"backend/pkg/game/entities"
	"fmt"
)

func RunSeeds() {
	for _, seedCharacter := range GetCharacterSeeds() {
		if character, _ := GetCharacterByName(seedCharacter.GetName()); character == nil {
			if err := SaveCharacter(seedCharacter); err != nil {
				fmt.Printf("Error saving new character to repository: %v\n", err)
			}
		}
	}
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
			map[entities.EquipmentType]*entities.Equipment{entities.Helmet: entities.GetEquipment("helm_001")}),
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
			map[entities.EquipmentType]*entities.Equipment{entities.Helmet: entities.GetEquipment("helm_002"), entities.Chest: entities.GetEquipment("chest_001")}),
	}
}
