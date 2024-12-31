package dtos

import (
	"slices"
	"sync"

	"tkz00/backend/pkg/game/entities"
	"tkz00/backend/pkg/utils"
)

type Mapper struct {
	entitiesIds []string
}

var mapperInstance *Mapper
var doItOnce sync.Once

func GetMapper() *Mapper {
	doItOnce.Do(func() {
		mapperInstance = &Mapper{
			entitiesIds: make([]string, 0),
		}
	})

	return mapperInstance
}

func (m *Mapper) GameStateToDTO(gameState entities.GameState) *GameStateDTO {
	playerDTOS := make([]CharacterDTO, 0)
	projectileDTOS := make([]ProjectileDTO, 0)
	npcsDTOS := make([]CharacterDTO, 0)
	AoEDTOS := make([]AoEDTO, 0)

	for _, player := range gameState.GetPlayers() {
		playerDTOS = append(playerDTOS, *m.CharacterToDTO(player))
	}

	for _, projectile := range gameState.GetProjectiles() {
		projectileDTOS = append(projectileDTOS, *m.ProjectileToDTO(projectile))
	}

	for _, npcs := range gameState.GetNPCs() {
		npcsDTOS = append(npcsDTOS, *m.CharacterToDTO(*npcs.Character))
	}

	areaEffects := gameState.AreaEffects()
	for AoEId, areaEffect := range areaEffects {
		if !slices.Contains(m.entitiesIds, AoEId) {
			AoEDTOS = append(AoEDTOS, m.AoEToDTO(*areaEffect))
			m.entitiesIds = append(m.entitiesIds, AoEId)
		}
	}
	destroyedAoEs := m.GetDestroyedAoEs(areaEffects)

	return &GameStateDTO{
		Players:           playerDTOS,
		Projectiles:       projectileDTOS,
		Npcs:              npcsDTOS,
		AreaEffects:       AoEDTOS,
		EntitiesToDestroy: destroyedAoEs,
	}
}

func (m Mapper) AoEToDTO(AoE entities.AoE) AoEDTO {
	return AoEDTO{
		Id:       AoE.Id(),
		Caster:   AoE.CasterId(),
		Position: *m.PositionToDTO(AoE.Position()),
		Radius:   AoE.Radius(),
	}
}

func (m *Mapper) GetDestroyedAoEs(areaEffects map[string]*entities.AoE) []string {
	destroyedAoEs := []string{}

	for _, id := range m.entitiesIds {
		if _, exists := areaEffects[id]; !exists {
			destroyedAoEs = append(destroyedAoEs, id)
		}
	}

	// Remove destroyed IDs from m.entitiesIds
	m.entitiesIds = filter(m.entitiesIds, func(id string) bool {
		return areaEffects[id] != nil
	})

	return destroyedAoEs
}

func filter(ids []string, predicate func(string) bool) []string {
	var result []string
	for _, id := range ids {
		if predicate(id) {
			result = append(result, id)
		}
	}
	return result
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
	characterHealth := character.GetHealth()
	characterAbilities := make([]AbilityDTO, 0)
	for _, ability := range character.GetAbilities() {
		abilityDTO := AbilityToDTO(*ability)
		abilityDTO.RemainingCooldown = float64(character.RemainingCooldown(ability)) / 1000
		characterAbilities = append(characterAbilities, abilityDTO)
	}

	characterExecutingAction := character.GetExecutingAction()
	executingActionDirection := characterExecutingAction.Direction()
	executingActionDTO := ExecutingActionDTO{
		Action:    characterExecutingAction.ActionType(),
		Direction: *m.PositionToDTO(executingActionDirection),
	}

	characterInventory := InventoryDTO{Items: character.ChangeLogs()}

	return &CharacterDTO{
		Id:              character.GetId(),
		Position:        *m.PositionToDTO(character.GetPosition()),
		Radius:          character.GetRadius(),
		MaxHealth:       characterHealth.GetMaxHealth(),
		CurrentHealth:   characterHealth.GetCurrentHealth(),
		ExecutingAction: executingActionDTO,
		Abilities:       characterAbilities,
		Inventory:       characterInventory,
	}
}

func (m Mapper) ProjectileToDTO(projectile entities.Projectile) *ProjectileDTO {
	return &ProjectileDTO{
		Id:       projectile.GetId(),
		Caster:   projectile.CasterId(),
		Position: *m.PositionToDTO(projectile.GetPosition()),
		Radius:   projectile.GetRadius(),
		State:    projectile.GetState(),
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
