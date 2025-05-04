package configurator

import (
	"encoding/json"
	"fmt"
	"os"
)

const AbilitiesFileName = "abilities.json"

func SaveAbilitiesToFile(abilities map[string]ConfiguratorAbility) error {
	file, err := os.OpenFile(AbilitiesFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(abilities)
}

func LoadAbilitiesFromFile() (map[string]ConfiguratorAbility, error) {
	file, err := os.Open(AbilitiesFileName)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	abilities := make(map[string]ConfiguratorAbility)
	if err := json.NewDecoder(file).Decode(&abilities); err != nil {
		return nil, fmt.Errorf("error decoding abilities from JSON: %v", err)
	}

	return abilities, nil
}

const PlayerAbilitiesFileName = "playerAbilities.json"

func SavePlayerInitialAbilities(abilitiesIds []string) error {
	file, err := os.OpenFile(PlayerAbilitiesFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(abilitiesIds)
}

func LoadPlayerAbilitiesIds() ([]string, error) {
	file, err := os.Open(PlayerAbilitiesFileName)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var playerAbilitiesIds []string
	if err := json.NewDecoder(file).Decode(&playerAbilitiesIds); err != nil {
		return nil, fmt.Errorf("error decoding abilities from JSON: %v", err)
	}

	return playerAbilitiesIds, nil
}
