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
		Radius: player.GetRadius(),
		MaxHealth: player.stats.maxHealth,
		CurrentHealth: player.stats.currentHealth,
		ExecutingAction: player.executingAction,
	}
}

func (m Mapper) ProjectileToDTO(projectile Projectile) *ProjectileDTO {
	return &ProjectileDTO {
		Id: projectile.id,
		Caster: projectile.caster,
		Position: *m.PositionToDTO(projectile.position),
		Radius: projectile.GetRadius(),
		Damage: projectile.damage,
		State: projectile.state.String(),
	}
}

func (m Mapper) GameStateToDTO(gameState GameState) *GameStateDTO {
	playerDTOS := make([]PlayerDTO, 0)
	projectileDTOS := make([]ProjectileDTO, 0)

	for _, player := range gameState.players {
		playerDTOS = append(playerDTOS, *m.PlayerToDTO(*player))
	}

	for _, projectile := range gameState.projectiles {
		projectileDTOS = append(projectileDTOS, *m.ProjectileToDTO(*projectile))
	}

	return &GameStateDTO {
		Players: playerDTOS,
		Projectiles: projectileDTOS,
	}
}
