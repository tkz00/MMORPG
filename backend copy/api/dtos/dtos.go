package dtos

import (
	"fmt"
	"tkz00/backend/pkg/game"
	"tkz00/backend/pkg/utils"
)

type DTO interface {
	GetType() string
}

type GameStateDTO struct {
	Players     []CharacterDTO  `json:"players"`
	Projectiles []ProjectileDTO `json:"projectiles"`
	Npcs        []CharacterDTO  `json:"npcs"`
}

func (g GameStateDTO) GetType() string {
	return "GameState"
}

type CharacterDTO struct {
	Id string `json:"id"`
	// MaxHealth     int         `json:"maxHealth"`
	// CurrentHealth int         `json:"currentHealth"`
	// Radius        float64     `json:"radius"`
	Position PositionDTO `json:"position"`
	// ExecutingAction ExecutingActionDTO `json:"executingAction"`
	// Abilities       []AbilityDTO       `json:"abilities"`
	// Inventory       InventoryDTO       `json:"inventory"`
}

func (p CharacterDTO) GetType() string {
	return "Player"
}

type PositionDTO struct {
	X float64 `json:"x"`
	Z float64 `json:"z"`
}

func CreatePositionDTO(data []byte) *PositionDTO {
	var positionDTO PositionDTO
	err := utils.FromJSON(data, &positionDTO)

	if err != nil {
		fmt.Println("Error creating PositionDTO", err)
		return &PositionDTO{}
	}

	return &positionDTO
}

func (p PositionDTO) GetType() string {
	return "position"
}

type ProjectileDTO struct {
	Id       string      `json:"id"`
	Caster   string      `json:"caster"`
	Position PositionDTO `json:"position"`
	Radius   float64     `json:"radius"`
	Damage   int         `json:"damage"`
	State    string      `json:"state"`
}

func (p ProjectileDTO) GetType() string {
	return "Projectile"
}

// type ExecutingActionDTO struct {
// 	Action    entities.Action `json:"action"`
// 	Direction PositionDTO     `json:"direction"`
// }

// type InventoryDTO struct {
// 	Items []entities.ItemChange `json:"items"`
// }

type UseItemDTO struct {
	Id string `json:"id"`
}

func (u UseItemDTO) GetType() string {
	return "use_item"
}

func GameStateToDTO(gameState game.GameState) GameStateDTO {
	playerDTOS := make([]CharacterDTO, 0)

	for id, player := range gameState.Players() {
		playerDTOS = append(playerDTOS, CharacterToDTO(id, player))
	}

	return GameStateDTO{
		Players: playerDTOS,
	}
}

func CharacterToDTO(id string, character game.Character) CharacterDTO {
	return CharacterDTO{
		Id:       id,
		Position: PositionToDTO(character.Position()),
	}
}

func PositionToDTO(position utils.Vector2) PositionDTO {
	x, z := position.GetPosition()
	return PositionDTO{
		X: x,
		Z: z,
	}
}

func PositionDTOToEntity(positionDTO PositionDTO) utils.Vector2 {
	return utils.NewVector2(positionDTO.X, positionDTO.Z)
}
