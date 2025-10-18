package configurator

import (
	"encoding/json"
	"fmt"
	"os"
)

func getAbilitiesFilePath() string {
	if path := os.Getenv("ABILITIES_FILE_PATH"); path != "" {
		return path
	}
	return "abilities.json"
}

func SaveAbilitiesToFile(abilities map[string]ConfiguratorAbility) error {
	abilitiesFileName := getAbilitiesFilePath()
	file, err := os.OpenFile(abilitiesFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(abilities)
}

func LoadAbilitiesFromFile() (map[string]ConfiguratorAbility, error) {
	abilitiesFileName := getAbilitiesFilePath()
	file, err := os.Open(abilitiesFileName)
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

func getInitialAbilitiesFilePath() string {
	if path := os.Getenv("INITIAL_ABILITIES_FILE_PATH"); path != "" {
		return path
	}
	return "playersInitialAbilities.json"
}

func SavePlayersInitialAbilities(abilitiesIds []string) error {
	playersInitialAbilitiesFileName := getInitialAbilitiesFilePath()
	file, err := os.OpenFile(
		playersInitialAbilitiesFileName,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0644,
	)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(abilitiesIds)
}

func LoadPlayersInitialAbilitiesIds() ([]string, error) {
	playersInitialAbilitiesFileName := getInitialAbilitiesFilePath()
	file, err := os.Open(playersInitialAbilitiesFileName)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var playersInitialAbilitiesIds []string
	if err := json.NewDecoder(file).Decode(&playersInitialAbilitiesIds); err != nil {
		return nil, fmt.Errorf("error decoding abilities from JSON: %v", err)
	}

	return playersInitialAbilitiesIds, nil
}
