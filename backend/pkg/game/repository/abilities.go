package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"tkz00/backend/pkg/configurator"
	"tkz00/backend/pkg/game/entities"
)

const ABILITIES_FILE_NAME = "abilities.json"

func LoadAbilitiesFromFile(filename string) (map[string]*entities.Ability, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	configuratorAbilities := make(map[string]configurator.ConfiguratorAbility)

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuratorAbilities)
	if err != nil {
		return nil, fmt.Errorf("error decoding abilities from JSON: %v", err)
	}

	abilities := make(map[string]*entities.Ability, len(configuratorAbilities))
	for i, configuratorAbilities := range configuratorAbilities {
		abilities[i] = ConvertFromConfiguratorAbility(configuratorAbilities)
	}

	return abilities, nil
}

func GetSkeletonEnemyAbilities(gs *entities.GameState) map[string]*entities.Ability {
	abilities, err := LoadAbilitiesFromFile(ABILITIES_FILE_NAME)
	if err != nil {
		fmt.Printf("Error loading abilities: %v\n", err)
	}

	return map[string]*entities.Ability{
		"0": abilities["0"],
	}
}

func GetPlayerAbilities() map[string]*entities.Ability {
	abilities, err := LoadAbilitiesFromFile(ABILITIES_FILE_NAME)
	if err != nil {
		fmt.Printf("Error loading abilities: %v\n", err)
	}

	return map[string]*entities.Ability{
		"1": abilities["1"],
		"2": abilities["2"],
		"3": abilities["3"],
		// "4": abilities["4"],
	}
}

func ConvertFromConfiguratorAbility(ability configurator.ConfiguratorAbility) *entities.Ability {
	// Recursively convert map[string]interface{} to Mechanic structs
	for _, mechanic := range ability.Mechanics {
		parseMechanics(mechanic.Params)
	}

	return entities.NewAbility(
		ability.ID,
		ability.Name,
		*ability.RangeValue,
		*ability.Cooldown,
		*ability.Targeting,
		*ability.CharacterState,
		ability.Mechanics...,
	)
}

// Recursive function to convert params maps to nested Mechanics
func parseMechanics(params map[string]interface{}) {
	if rawMechanics, ok := params["on_hit_mechanics"]; ok {
		if mechanicList, ok := rawMechanics.([]interface{}); ok {
			var mechanics []entities.Mechanic
			for _, m := range mechanicList {
				if mechMap, ok := m.(map[string]interface{}); ok {
					mechanic := entities.Mechanic{
						MechanicType: mechMap["mechanic_type"].(string),
						Params:       mechMap["params"].(map[string]interface{}),
					}
					// Recursively parse nested mechanics
					parseMechanics(mechanic.Params)
					mechanics = append(mechanics, mechanic)
				}
			}
			params["on_hit_mechanics"] = mechanics
		}
	} else if rawMechanics, ok := params["execute_after_delay_mechanics"]; ok {
		if mechanicList, ok := rawMechanics.([]interface{}); ok {
			var mechanics []entities.Mechanic
			for _, m := range mechanicList {
				if mechMap, ok := m.(map[string]interface{}); ok {
					mechanic := entities.Mechanic{
						MechanicType: mechMap["mechanic_type"].(string),
						Params:       mechMap["params"].(map[string]interface{}),
					}
					// Recursively parse nested mechanics
					parseMechanics(mechanic.Params)
					mechanics = append(mechanics, mechanic)
				}
			}
			params["execute_after_delay_mechanics"] = mechanics
		}
	}
}
