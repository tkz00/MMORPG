package dtos

import (
	"sync"

	"tkz00/backend/pkg/game/entities"
	"tkz00/backend/pkg/utils"
)

type Mapper struct{}

var mapperInstance *Mapper
var doItOnce sync.Once

func GetMapper() *Mapper {
	doItOnce.Do(func() {
		mapperInstance = &Mapper{}
	})

	return mapperInstance
}

func (m Mapper) PositionDTOToEntity(positionDTO PositionDTO) *utils.Vector2 {
	return utils.NewVector2(positionDTO.X, positionDTO.Z)
}

func (m Mapper) PositionToDTO(position utils.Vector2) *PositionDTO {
	x, z := position.GetPosition()
	return &PositionDTO{
		X: x,
		Z: z,
	}
}

func (m Mapper) CharacterToDTO(character entities.Character) *CharacterDTO {
	playerHealth := character.GetHealth()
	playerAbilities := make([]AbilityDTO, 0)
	for _, ability := range character.GetAbilities() {
		abilityDTO := AbilityToDTO(*ability)
		abilityDTO.RemainingCooldown = float64(character.RemainingCooldown(ability)) / 1000
		playerAbilities = append(playerAbilities, abilityDTO)
	}

	characterExecutingAction := character.GetExecutingAction()
	executingActionDirection := characterExecutingAction.Direction()
	executingActionDTO := ExecutingActionDTO{Action: characterExecutingAction.ActionType(), Direction: *m.PositionToDTO(executingActionDirection)}

	return &CharacterDTO{
		Id:              character.GetId(),
		Position:        *m.PositionToDTO(character.GetPosition()),
		Radius:          character.GetRadius(),
		MaxHealth:       playerHealth.GetMaxHealth(),
		CurrentHealth:   playerHealth.GetCurrentHealth(),
		ExecutingAction: executingActionDTO,
		Abilities:       playerAbilities,
	}
}

func (m Mapper) ProjectileToDTO(projectile entities.Projectile) *ProjectileDTO {
	return &ProjectileDTO{
		Id:       projectile.GetId(),
		Caster:   projectile.CasterId(),
		Position: *m.PositionToDTO(projectile.GetPosition()),
		Radius:   projectile.GetRadius(),
		Damage:   projectile.GetDamage(),
		State:    projectile.GetState(),
	}
}

func (m Mapper) GameStateToDTO(gameState entities.GameState) *GameStateDTO {
	playerDTOS := make([]CharacterDTO, 0)
	projectileDTOS := make([]ProjectileDTO, 0)
	npcsDTOS := make([]CharacterDTO, 0)

	for _, player := range gameState.GetPlayers() {
		playerDTOS = append(playerDTOS, *m.CharacterToDTO(player))
	}

	for _, projectile := range gameState.GetProjectiles() {
		projectileDTOS = append(projectileDTOS, *m.ProjectileToDTO(projectile))
	}

	for _, npcs := range gameState.GetNPCs() {
		npcsDTOS = append(npcsDTOS, *m.CharacterToDTO(*npcs.Character))
	}

	return &GameStateDTO{
		Players:     playerDTOS,
		Projectiles: projectileDTOS,
		Npcs:        npcsDTOS,
	}
}

func AbilityToDTO(ability entities.Ability) AbilityDTO {
	return AbilityDTO{
		Id:       ability.Id(),
		Name:     ability.Name(),
		Range:    ability.Range(),
		Cooldown: ability.Cooldown(),
	}
}
