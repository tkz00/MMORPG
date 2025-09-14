package repository

import (
	"backend/pkg/game/entities"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/samber/lo"
)

const CharactersFileName = "characters.json"

type Character struct {
	Id            string                            `json:"id"`
	CharacterName string                            `json:"character_name"`
	MaxHealth     int                               `json:"max_health"`
	CurrentHealth int                               `json:"current_health"`
	BaseStats     map[string]int64                  `json:"base_stats"`
	X             float64                           `json:"x_pos"`
	Z             float64                           `json:"z_pos"`
	Items         map[string]int64                  `json:"items"`
	Equipped      map[entities.EquipmentType]string `json:"equipped"`
	Abilities     []string                          `json:"abilities"`
}

func GetCharacterByName(characterName string) (*entities.Character, error) {
	characters, err := GetCharacters()

	if err != nil {
		return nil, err
	}

	if character, ok := lo.Find(lo.Values(characters), func(c Character) bool {
		return c.CharacterName == characterName
	}); ok {
		return entities.CreateCharacter(character.Id, characterName, character.X, character.Z, character.BaseStats, GetAbilitiesByIds(character.Abilities)), nil
	}

	return nil, fmt.Errorf("character with name %s not found", characterName)
}

func GetCharacters() (map[string]Character, error) {
	file, err := os.Open(CharactersFileName)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	characters := make(map[string]Character)
	if err := json.NewDecoder(file).Decode(&characters); err != nil {
		return nil, fmt.Errorf("error decoding characters from JSON: %v", err)
	}

	return characters, nil
}

func SaveCharacter(newCharacter *entities.Character) error {
	characters := make(map[string]Character)

	// try reading existing file
	if f, err := os.Open(CharactersFileName); err == nil {
		defer f.Close()
		if err := json.NewDecoder(f).Decode(&characters); err != nil && !errors.Is(err, io.EOF) {
			return fmt.Errorf("error decoding characters from JSON: %w", err)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("error opening characters file: %w", err)
	}

	// add/replace character
	x, z := newCharacter.GetPosition().GetPosition()
	characters[newCharacter.GetId()] = Character{
		newCharacter.GetId(),
		newCharacter.GetName(),
		newCharacter.GetMaxHealth(),
		newCharacter.GetCurrentHealth(),
		newCharacter.GetBaseStats(),
		x, z,
		newCharacter.GetInventory(),
		make(map[entities.EquipmentType]string),
		lo.Keys(newCharacter.GetAbilities()),
	}

	// rewrite full file
	f, err := os.Create(CharactersFileName)
	if err != nil {
		return fmt.Errorf("error opening file for write: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(characters); err != nil {
		return fmt.Errorf("error encoding characters to JSON: %w", err)
	}

	return nil
}
