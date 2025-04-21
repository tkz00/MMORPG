package dtos

import (
	"sync"

	"backend/pkg/game/entities"
	"backend/pkg/utils"
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
		projectileDTOS = append(projectileDTOS, *GetProjectileDiff(&projectile))
	}

	for _, npcs := range gameState.GetNPCs() {
		npcsDTOS = append(npcsDTOS, *m.CharacterToDTO(*npcs.Character))
	}

	areaEffects := gameState.AreaEffects()
	for _, AoE := range areaEffects {
		AoEDTOS = append(AoEDTOS, *GetAoEDiff(AoE))
	}

	return &GameStateDTO{
		Players:     playerDTOS,
		Projectiles: projectileDTOS,
		Npcs:        npcsDTOS,
		AreaEffects: AoEDTOS,
	}
}

func GetAoEDiff(AoE *entities.AoE) *AoEDTO {
	diff := &AoEDTO{Id: AoE.Id()}
	if AoE.AoELastTickState.Position == nil ||
		!AoE.Position().Equals(*AoE.AoELastTickState.Position) {
		diff.Position = PositionToDTO(AoE.Position())
		AoE.AoELastTickState.Position = utils.NewVector2(AoE.Position().GetPosition())
	}
	if AoE.AoELastTickState.Radius == nil ||
		AoE.Radius() != *AoE.AoELastTickState.Radius {
		radius := AoE.Radius()
		diff.Radius = &radius
		AoE.AoELastTickState.Radius = &radius
	}
	return diff
}

func (m Mapper) PositionDTOToEntity(positionDTO PositionDTO) *utils.Vector2 {
	return utils.NewVector2(positionDTO.X, positionDTO.Z)
}

func PositionToDTO(position utils.Vector2) *PositionDTO {
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
		Direction: *PositionToDTO(executingActionDirection),
	}

	characterInventory := InventoryDTO{Items: character.ChangeLogs()}

	return &CharacterDTO{
		Id:              character.GetId(),
		Position:        *PositionToDTO(character.GetPosition()),
		Radius:          character.GetRadius(),
		MaxHealth:       characterHealth.GetMaxHealth(),
		CurrentHealth:   characterHealth.GetCurrentHealth(),
		ExecutingAction: executingActionDTO,
		Abilities:       characterAbilities,
		Inventory:       characterInventory,
	}
}

func GetProjectileDiff(p *entities.Projectile) *ProjectileDTO {
	diff := &ProjectileDTO{Id: p.GetId()}
	if p.ProjectileLastTickState.Position == nil ||
		!p.GetPosition().Equals(*p.ProjectileLastTickState.Position) {
		diff.Position = PositionToDTO(p.GetPosition())
		p.ProjectileLastTickState.Position = utils.NewVector2(p.GetPosition().GetPosition())
	}
	if p.ProjectileLastTickState.State == nil ||
		p.State() != *p.ProjectileLastTickState.State {
		diff.State = p.GetState()
		state := p.State()
		p.ProjectileLastTickState.State = &state
	}
	if p.ProjectileLastTickState.Radius == nil ||
		p.GetRadius() != *p.ProjectileLastTickState.Radius {
		radius := p.GetRadius()
		diff.Radius = &radius
		p.ProjectileLastTickState.Radius = &radius
	}
	return diff
}

func AbilityToDTO(ability entities.Ability) AbilityDTO {
	return AbilityDTO{
		Id:       ability.Id(),
		Name:     ability.Name(),
		Range:    ability.Range(),
		Cooldown: ability.Cooldown(),
	}
}
