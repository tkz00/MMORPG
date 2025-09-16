package repository

import (
	"backend/pkg/game/entities"

	"gorm.io/gorm"
)

const CharactersFileName = "characters.json"

type CharacterDB struct {
	Id            string `gorm:"primaryKey"`
	CharacterName string
	MaxHealth     int
	CurrentHealth int
	BaseStats     StatsMap `gorm:"type:jsonb"`
	X             float64
	Z             float64
	Items         ItemsMap     `gorm:"type:jsonb"`
	Equipped      EquipmentMap `gorm:"type:jsonb"`
	// Abilities     []string

	gorm.Model
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
	character := entities.CreateCharacter(repoCharacter.Id, characterName, repoCharacter.X, repoCharacter.Z, repoCharacter.MaxHealth, repoCharacter.CurrentHealth, repoCharacter.BaseStats, GetPlayersInitialAbilities())
	for id, qty := range repoCharacter.Items {
		character.Inventory.AddItem(id, qty)
	}
	for _, id := range repoCharacter.Equipped {
		character.EquipItem(id)
	}
	return character, nil
}

func FromEntity(c *entities.Character) (*CharacterDB, error) {
	x, z := c.GetPosition().GetPosition()
	equipped := make(EquipmentMap)
	for k, v := range c.GetEquipped() {
		equipped[string(k)] = v.Id()
	}

	return &CharacterDB{
		Id:            c.GetId(),
		CharacterName: c.GetName(),
		MaxHealth:     c.GetMaxHealth(),
		CurrentHealth: c.GetCurrentHealth(),
		BaseStats:     StatsMap(c.GetBaseStats()),
		Items:         ItemsMap(c.GetInventory()),
		Equipped:      equipped,
		X:             x,
		Z:             z,
	}, nil
}

func SaveCharacter(c *entities.Character) error {
	dbChar, err := FromEntity(c)
	if err != nil {
		return err
	}

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
