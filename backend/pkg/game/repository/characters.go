package repository

import (
	"backend/api/dtos"
	"backend/pkg/game/entities"
	"slices"
	"strings"

	"github.com/lib/pq"
	"github.com/samber/lo"
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
	Items         ItemsMap       `gorm:"type:jsonb"`
	Equipped      EquipmentMap   `gorm:"type:jsonb"`
	Abilities     pq.StringArray `gorm:"type:text[]"`

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

	character := entities.CreateCharacter(
		repoCharacter.Id,
		characterName,
		repoCharacter.X,
		repoCharacter.Z,
		repoCharacter.MaxHealth,
		repoCharacter.CurrentHealth,
		repoCharacter.BaseStats,
		GetAbilitiesByIds(repoCharacter.Abilities),
		repoCharacter.Items,
		lo.MapEntries(repoCharacter.Equipped, func(k string, v string) (entities.EquipmentType, *entities.Equipment) {
			return entities.EquipmentType(k), entities.GetEquipment(v)
		}),
	)

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
		Abilities:     lo.Keys(c.GetAbilities()),
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

func GetAllCharacters() ([]dtos.CharacterDTO, error) {
	var charactersDB []CharacterDB
	if result := DB.Find(&charactersDB); result.Error != nil {
		return nil, result.Error
	}

	characters := make([]dtos.CharacterDTO, len(charactersDB))
	for i, c := range charactersDB {
		characterInventory := dtos.InventoryDTO{}
		equipped := lo.Values(c.Equipped)
		for item, quantity := range c.Items {
			characterInventory.Items = append(
				characterInventory.Items,
				dtos.ItemDTO{Id: item, Quantity: quantity, IsEquipped: lo.Contains(equipped, item)},
			)
		}

		characterAbilities := make([]dtos.AbilityDTO, 0)
		for _, ability := range GetAbilitiesByIds(c.Abilities) {
			abilityDTO := dtos.AbilityToDTO(*ability)
			characterAbilities = append(characterAbilities, abilityDTO)
		}
		slices.SortFunc(characterAbilities, func(a, b dtos.AbilityDTO) int {
			return strings.Compare(a.Id, b.Id)
		})

		characters[i] = dtos.CharacterDTO{
			Id:            c.Id,
			Name:          c.CharacterName,
			MaxHealth:     &c.MaxHealth,
			CurrentHealth: &c.CurrentHealth,
			Stats:         (*map[string]int64)(&c.BaseStats),
			Position:      &dtos.PositionDTO{X: c.X, Z: c.Z},
			Inventory:     &characterInventory,
			Abilities:     characterAbilities,
		}
	}

	return characters, nil
}
