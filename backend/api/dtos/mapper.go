package dtos

import (
	"sync"

	"unnamed-mmo/backend/pkg/game"
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

func (m Mapper) PositionDTOToEntity(positionDTO PositionDTO) *game.Position {
	return game.NewPosition(positionDTO.X, positionDTO.Z)
}

func (m Mapper) PositionToDTO(position game.Position) *PositionDTO {
	x, z := position.GetPosition()
	return &PositionDTO {
		X: x,
		Z: z,
	}
}

func (m Mapper) PlayerToDTO(player game.Player) *PlayerDTO {
	playerStats := player.GetStats()
	playerAbilities := make([]AbilityDTO, 0)
	for _, ability := range player.GetAbilities() {
		abilityDTO := AbilityToDTO(*ability)
		abilityDTO.RemainingCooldown = float64(player.RemainingCooldown(ability)) / 1000
		playerAbilities = append(playerAbilities, abilityDTO)
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

func (m Mapper) ProjectileToDTO(projectile game.Projectile) *ProjectileDTO {
	return &ProjectileDTO {
		Id: projectile.GetId(),
		Caster: projectile.GetCaster(),
		Position: *m.PositionToDTO(projectile.GetPosition()),
		Radius: projectile.GetRadius(),
		Damage: projectile.GetDamage(),
		State: projectile.GetState(),
	}
}

func (m Mapper) GameStateToDTO(gameState game.GameState) *GameStateDTO {
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

func AbilityToDTO(ability game.Ability) AbilityDTO{
	return AbilityDTO{
		Id: ability.Id(),
		Name: ability.Name(),
		Range: ability.Range(),
		Cooldown: ability.Cooldown(),
	}
}
