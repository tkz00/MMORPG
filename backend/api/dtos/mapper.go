package dtos

import (
	"backend/pkg/game/entities"
	"backend/pkg/utils"
)

func GameStateDiff(gameState entities.GameState) GameStateDTO {
	playerDTOS := make([]CharacterDTO, 0)
	projectileDTOS := make([]ProjectileDTO, 0)
	npcsDTOS := make([]CharacterDTO, 0)
	AoEDTOS := make([]AoEDTO, 0)

	for _, player := range gameState.GetPlayers() {
		playerDTOS = append(playerDTOS, *GetCharacterDiff(&player))
	}

	for _, projectile := range gameState.GetProjectiles() {
		projectileDTOS = append(projectileDTOS, *GetProjectileDiff(&projectile))
	}

	for _, npcs := range gameState.GetNPCs() {
		npcsDTOS = append(npcsDTOS, *GetCharacterDiff(npcs.Character))
	}

	for _, AoE := range gameState.AreaEffects() {
		AoEDTOS = append(AoEDTOS, *GetAoEDiff(AoE))
	}

	return GameStateDTO{
		Players:     playerDTOS,
		Projectiles: projectileDTOS,
		Npcs:        npcsDTOS,
		AreaEffects: AoEDTOS,
	}
}

func PositionDTOToEntity(positionDTO PositionDTO) *utils.Vector2 {
	return utils.NewVector2(positionDTO.X, positionDTO.Z)
}

func PositionToDTO(position utils.Vector2) *PositionDTO {
	x, z := position.GetPosition()
	return &PositionDTO{
		X: x,
		Z: z,
	}
}

func CharacterToDTO(character entities.Character) CharacterDTO {
	characterHealth := character.GetHealth()
	characterAbilities := make([]AbilityDTO, 0)
	for _, ability := range character.GetAbilities() {
		abilityDTO := AbilityToDTO(*ability)
		abilityDTO.RemainingCooldown = float64(character.RemainingCooldown(ability)) / 1000
		characterAbilities = append(characterAbilities, abilityDTO)
	}

	currentHealth := characterHealth.GetCurrentHealth()
	maxHealth := characterHealth.GetMaxHealth()
	radius := character.GetRadius()
	characterExecutingAction := character.GetExecutingAction()
	executingActionDirection := characterExecutingAction.Direction()
	actionType := characterExecutingAction.ActionType()
	characterInventory := InventoryDTO{Items: character.GetInventory()}

	return CharacterDTO{
		Id:            character.GetId(),
		Position:      PositionToDTO(character.GetPosition()),
		Radius:        &radius,
		MaxHealth:     &maxHealth,
		CurrentHealth: &currentHealth,
		Action:        &actionType,
		Direction:     PositionToDTO(executingActionDirection),
		Abilities:     characterAbilities,
		Inventory:     characterInventory,
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

func GetCharacterDiff(c *entities.Character) *CharacterDTO {
	diff := &CharacterDTO{Id: c.GetId()}
	if c.CharacterLastTickState.Position == nil ||
		!c.GetPosition().Equals(*c.CharacterLastTickState.Position) {
		diff.Position = PositionToDTO(c.GetPosition())
		c.CharacterLastTickState.Position = utils.NewVector2(c.GetPosition().GetPosition())
	}
	if c.CharacterLastTickState.Radius == nil ||
		c.GetRadius() != *c.CharacterLastTickState.Radius {
		radius := c.GetRadius()
		diff.Radius = &radius
		c.CharacterLastTickState.Radius = &radius
	}
	if c.CharacterLastTickState.MaxHealth == nil ||
		c.GetMaxHealth() != *c.CharacterLastTickState.MaxHealth {
		maxHealth := c.GetMaxHealth()
		diff.MaxHealth = &maxHealth
		c.CharacterLastTickState.MaxHealth = &maxHealth
	}
	if c.CharacterLastTickState.CurrentHealth == nil ||
		c.GetCurrentHealth() != *c.CharacterLastTickState.CurrentHealth {
		currentHealth := c.GetCurrentHealth()
		diff.CurrentHealth = &currentHealth
		c.CharacterLastTickState.CurrentHealth = &currentHealth
	}
	if c.CharacterLastTickState.Action == nil ||
		c.GetExecutingAction().
			ActionType() !=
			*c.CharacterLastTickState.Action {
		action := c.GetExecutingAction().ActionType()
		diff.Action = &action
		c.CharacterLastTickState.Action = &action
	}
	if c.CharacterLastTickState.Direction == nil ||
		c.GetExecutingAction().
			Direction().Equals(*c.CharacterLastTickState.Direction) {
		diff.Direction = PositionToDTO(c.GetExecutingAction().Direction())
		c.CharacterLastTickState.Direction = utils.NewVector2(
			c.GetExecutingAction().Direction().GetPosition(),
		)
	}
	// move this into a function and do diff check
	characterAbilities := make([]AbilityDTO, 0)
	for _, ability := range c.GetAbilities() {
		abilityDTO := AbilityToDTO(*ability)
		abilityDTO.RemainingCooldown = float64(c.RemainingCooldown(ability)) / 1000
		characterAbilities = append(characterAbilities, abilityDTO)
	}
	diff.Abilities = characterAbilities
	characterInventory := InventoryDTO{Items: c.ChangeLogs()}
	diff.Inventory = characterInventory
	return diff
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

func GameStateToDTO(gameState entities.GameState) *GameStateDTO {
	playerDTOS := make([]CharacterDTO, 0)
	projectileDTOS := make([]ProjectileDTO, 0)
	npcsDTOS := make([]CharacterDTO, 0)
	AoEDTOS := make([]AoEDTO, 0)

	for _, player := range gameState.GetPlayers() {
		playerDTOS = append(playerDTOS, CharacterToDTO(player))
	}

	for _, projectile := range gameState.GetProjectiles() {
		projectileDTOS = append(projectileDTOS, ProjectileTODTO(projectile))
	}

	for _, npcs := range gameState.GetNPCs() {
		npcsDTOS = append(npcsDTOS, CharacterToDTO(*npcs.Character))
	}

	areaEffects := gameState.AreaEffects()
	for _, AoE := range areaEffects {
		AoEDTOS = append(AoEDTOS, AoEToDTO(AoE))
	}

	return &GameStateDTO{
		Players:     playerDTOS,
		Projectiles: projectileDTOS,
		Npcs:        npcsDTOS,
		AreaEffects: AoEDTOS,
	}
}

func ProjectileTODTO(p entities.Projectile) ProjectileDTO {
	radius := p.GetRadius()
	projectile := ProjectileDTO{
		Id:       p.GetId(),
		Caster:   p.CasterId(),
		Position: PositionToDTO(p.GetPosition()),
		Radius:   &radius,
		State:    p.GetState(),
	}
	return projectile
}

func AoEToDTO(AoE *entities.AoE) AoEDTO {
	radius := AoE.Radius()
	AoEDTO := AoEDTO{
		Id:       AoE.Id(),
		Position: PositionToDTO(AoE.Position()),
		Radius:   &radius,
	}
	return AoEDTO
}
