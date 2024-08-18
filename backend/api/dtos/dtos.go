package dtos

import (
	"fmt"
	"unnamed-mmo/backend/pkg/game/character"
	"unnamed-mmo/backend/pkg/utils"
)

type DTO interface {
	GetType() string
}

type GameStateDTO struct {
	Players 	[]PlayerDTO 	`json:"players"`
	Projectiles []ProjectileDTO `json:"projectiles"`
}

func (g GameStateDTO) GetType() string {
	return "GameState"
}

type PlayerDTO struct {
	Id       		string      		`json:"id"`
	MaxHealth 		int					`json:"maxHealth"`
	CurrentHealth 	int					`json:"currentHealth"`
	Radius 			float64 			`json:"radius"`	
	Position 		PositionDTO 		`json:"position"`
	ExecutingAction	character.Action	`json:"executingAction"`
	Abilities		[]AbilityDTO		`json:"abilities"`
	
}

func (p PlayerDTO) GetType() string {
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
	return "Position"
}

type ProjectileDTO struct {
	Id			string 		`json:"id"`
	Caster 		string 		`json:"caster"`
	Position 	PositionDTO `json:"position"`
	Radius 		float64 	`json:"radius"`
	Damage		int 		`json:"damage"`
	State		string 		`json:"state"`
}

func (p ProjectileDTO) GetType() string {
	return "Projectile"
}
