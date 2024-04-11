package types

import (
	"sync"
)

type Mapper struct {}

var mapperInstance *Mapper
var doItOnce sync.Once

func GetMapper() *Mapper {
	doItOnce.Do(func () {
		mapperInstance = &Mapper{}
	})

	return mapperInstance
}

func (m Mapper) PositionDTOToEntity(positionDTO PositionDTO) *Position {
	return &Position {
		x: positionDTO.X,
		z: positionDTO.Z,
	}
}

func (m Mapper) PositionToDTO(position Position) *PositionDTO {
	return &PositionDTO {
		X: position.x,
		Z: position.z,
	}
}

func (m Mapper) PlayerDTOToEntity(playerDTO PlayerDTO) *Player {
	initialPosition := *m.PositionDTOToEntity(playerDTO.Position)

	return &Player {
		position: initialPosition,
		to: initialPosition,
		stats: PlayerStats{
			maxHealth: playerDTO.MaxHealth,
			currentHealth: playerDTO.CurrentHealth,
		},
	}
}

func (m Mapper) PlayerToDTO(player Player) *PlayerDTO {
	return &PlayerDTO {
		Id: player.id,
		Position: *m.PositionToDTO(player.position),
		MaxHealth: player.stats.maxHealth,
		CurrentHealth: player.stats.currentHealth,
	}
}

func (m Mapper) GameStateToDTO(gameState GameState) *GameStateDTO {
	playerDTOS := make([]PlayerDTO, 0)

	for _, player := range gameState.players {
		playerDTOS = append(playerDTOS, *m.PlayerToDTO(*player))
	}

	return &GameStateDTO {
		Players: playerDTOS,
	}
}