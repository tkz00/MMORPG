package repository

import (
	"backend/pkg/game/entities"
)

const CharactersFileName = "characters.json"

type CharacterDB struct {
	Id            string `gorm:"primaryKey"`
	CharacterName string
	MaxHealth     int
	CurrentHealth int
	// BaseStats     map[string]int64
	X float64
	Z float64
	// Items         map[string]int64
	// Equipped      map[entities.EquipmentType]string
	// Abilities     []string

	// gorm.Model
}

func (CharacterDB) TableName() string {
	return "characters"
}

func GetCharacterByName(characterName string) (*entities.Character, error) {
	var repoCharacter CharacterDB
	if err := DB.Where("character_name = ?", characterName).First(&repoCharacter).Error; err != nil {
		return nil, err
	}

	// character := entities.CreateCharacter(repoCharacter.Id, characterName, repoCharacter.X, repoCharacter.Z, repoCharacter.MaxHealth, repoCharacter.CurrentHealth, repoCharacter.BaseStats, GetAbilitiesByIds(repoCharacter.Abilities))
	character := entities.CreateCharacter(repoCharacter.Id, characterName, repoCharacter.X, repoCharacter.Z, repoCharacter.MaxHealth, repoCharacter.CurrentHealth, make(map[string]int64), GetPlayersInitialAbilities())
	// for id, qty := range repoCharacter.Items {
	// 	character.Inventory.AddItem(id, qty)
	// }
	// for _, id := range repoCharacter.Equipped {
	// 	character.EquipItem(id)
	// }
	return character, nil
}

// func GetCharacters() (map[string]CharacterDB, error) {
// 	file, err := os.Open(CharactersFileName)
// 	if err != nil {
// 		return nil, fmt.Errorf("error opening file: %v", err)
// 	}
// 	defer file.Close()

// 	characters := make(map[string]CharacterDB)
// 	if err := json.NewDecoder(file).Decode(&characters); err != nil {
// 		return nil, fmt.Errorf("error decoding characters from JSON: %v", err)
// 	}

// 	return characters, nil
// }

func FromEntity(c *entities.Character) (*CharacterDB, error) {
	x, z := c.GetPosition().GetPosition()

	return &CharacterDB{
		Id:            c.GetId(),
		CharacterName: c.GetName(),
		MaxHealth:     c.GetMaxHealth(),
		CurrentHealth: c.GetCurrentHealth(),
		// BaseStats:     c.GetBaseStats(),
		X: x,
		Z: z,
	}, nil
}

func SaveCharacter(c *entities.Character) error {
	dbChar, err := FromEntity(c)
	if err != nil {
		return err
	}

	// Insert or update (upsert-like behavior)
	return DB.Save(dbChar).Error
}

// func SaveCharacter(newCharacter *entities.Character) error {
// 	characters := make(map[string]CharacterDB)

// 	// try reading existing file
// 	if f, err := os.Open(CharactersFileName); err == nil {
// 		defer f.Close()
// 		if err := json.NewDecoder(f).Decode(&characters); err != nil && !errors.Is(err, io.EOF) {
// 			return fmt.Errorf("error decoding characters from JSON: %w", err)
// 		}
// 	} else if !errors.Is(err, os.ErrNotExist) {
// 		return fmt.Errorf("error opening characters file: %w", err)
// 	}

// 	// add/replace character
// 	x, z := newCharacter.GetPosition().GetPosition()
// 	equipped := lo.MapValues(
// 		newCharacter.GetEquipped(),
// 		func(eq *entities.Equipment, slot entities.EquipmentType) string {
// 			if eq == nil {
// 				return "" // or some placeholder
// 			}
// 			return eq.Id()
// 		},
// 	)

// 	characters[newCharacter.GetId()] = CharacterDB{
// 		newCharacter.GetId(),
// 		newCharacter.GetName(),
// 		newCharacter.GetMaxHealth(),
// 		newCharacter.GetCurrentHealth(),
// 		newCharacter.GetBaseStats(),
// 		x, z,
// 		newCharacter.GetInventory(),
// 		equipped,
// 		lo.Keys(newCharacter.GetAbilities()),
// 	}

// 	// rewrite full file
// 	f, err := os.Create(CharactersFileName)
// 	if err != nil {
// 		return fmt.Errorf("error opening file for write: %w", err)
// 	}
// 	defer f.Close()

// 	enc := json.NewEncoder(f)
// 	enc.SetIndent("", "  ")
// 	if err := enc.Encode(characters); err != nil {
// 		return fmt.Errorf("error encoding characters to JSON: %w", err)
// 	}

// 	return nil
// }
