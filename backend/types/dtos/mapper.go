package dtos

import (
	"sync"

	"unnamed-mmo/backend/types"
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

func (m Mapper) PositionDTOToEntity(positionDTO PositionDTO) *types.Position {
	return types.NewPosition(positionDTO.X, positionDTO.Z)
}

func (m Mapper) PositionToDTO(position types.Position) *PositionDTO {
	x, z := position.GetPosition()
	return &PositionDTO {
		X: x,
		Z: z,
	}
}

func (m Mapper) PlayerToDTO(player types.Player) *PlayerDTO {
	playerStats := player.GetStats()

	playerAbilities := make([]AbilityDTO, 0)
	for _, ability := range player.GetAbilities() {
		playerAbilities = append(playerAbilities, AbilityToDTO(ability))
	}


	return &PlayerDTO {
		Id: player.GetId(),
		Position: *m.PositionToDTO(player.GetPosition()),
		Radius: player.GetRadius(),
		MaxHealth: playerStats.GetMaxHealth(),
		CurrentHealth: playerStats.GetCurrentHealth(),
		ExecutingAction: player.GetExecutingAction(),
		Abilities: playerAbilities,
	}
}

func (m Mapper) ProjectileToDTO(projectile types.Projectile) *ProjectileDTO {
	return &ProjectileDTO {
		Id: projectile.GetId(),
		Caster: projectile.GetCaster(),
		Position: *m.PositionToDTO(projectile.GetPosition()),
		Radius: projectile.GetRadius(),
		Damage: projectile.GetDamage(),
		State: projectile.GetState(),
	}
}

func (m Mapper) GameStateToDTO(gameState types.GameState) *GameStateDTO {
	playerDTOS := make([]PlayerDTO, 0)
	projectileDTOS := make([]ProjectileDTO, 0)

	for _, player := range gameState.GetPlayers() {
		playerDTOS = append(playerDTOS, *m.PlayerToDTO(player))
	}

	for _, projectile := range gameState.GetProjectiles() {
		projectileDTOS = append(projectileDTOS, *m.ProjectileToDTO(projectile))
	}

	return &GameStateDTO {
		Players: playerDTOS,
		Projectiles: projectileDTOS,
	}
}

func AbilityToDTO(ability types.Ability) AbilityDTO{
	return AbilityDTO{
		Id: ability.GetId(),
		Name: ability.GetName(),
		Range: ability.GetRange(),
	}
}

// func (m Mapper) AbilityCastToParameters(abilityCast AbilityCastDTO) types.AbilityInfo {
// 	return 
// }
