package types

import (
	"fmt"
	"unnamed-mmo/backend/utils"
)

type DTO interface {
	GetType() string
}

type GameStateDTO struct {
	Players []PlayerDTO `json:"players"`
}

func (g GameStateDTO) GetType() string {
	return utils.GetTypeName(g)
}

type PlayerDTO struct {
	Id       string      `json:"id"`
	Position PositionDTO `json:"position"`
}

func (p PlayerDTO) GetType() string {
	return utils.GetTypeName(p)
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
	return utils.GetTypeName(p)
}