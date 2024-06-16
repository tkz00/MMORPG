package types

import (
	"fmt"
	"unnamed-mmo/backend/utils"
)

type DTO interface {
	GetType() string
}

type GameStateDTO struct {
	Players 	[]PlayerDTO `json:"players"`
	Projectiles []ProjectileDTO `json:"projectiles"`
}

func (g GameStateDTO) GetType() string {
	return "GameState"
}

type PlayerDTO struct {
	Id       		string      `json:"id"`
	MaxHealth 		int			`json:"maxHealth"`	
	CurrentHealth 	int			`json:"currentHealth"`	
	Position 		PositionDTO `json:"position"`
}

func (p PlayerDTO) GetType() string {
	return "Player"
}

type PositionDTO struct {
	X float32 `json:"x"`
	Z float32 `json:"z"`
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
	return "Position"
}

type ProjectileDTO struct {
	Id			string 		`json:"id"`
	Caster 		string 		`json:"caster"`
	Position 	PositionDTO `json:"position"`
	Damage		int 		`json:"damage"`
}

func (p ProjectileDTO) GetType() string {
	return "Projectile"
}
