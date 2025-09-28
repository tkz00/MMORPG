package dtos

import (
	"backend/config"
	"backend/pkg/game/entities"
	"backend/pkg/utils"
	"slices"
	"strings"

	"github.com/samber/lo"
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
		if c.GetExecutingAction().DurationInRemainingTicks() != nil {
			actionDurationMs := *c.GetExecutingAction().
				DurationInRemainingTicks() *
				config.TICKER_TIME.Milliseconds()
			diff.ActionDurationMs = &actionDurationMs
		}
	}
	if c.CharacterLastTickState.Direction == nil ||
		!c.GetExecutingAction().Direction().Equals(*c.CharacterLastTickState.Direction) {
		diff.Direction = PositionToDTO(c.GetExecutingAction().Direction())
		c.CharacterLastTickState.Direction = utils.NewVector2(
			c.GetExecutingAction().Direction().GetPosition(),
		)
	}
	diff.Stats = getStatsDiff(c)
	diff.Abilities = getAbilitiesDiff(c)
	diff.Inventory = getInventoryDiff(c)
	return diff
}

func getAbilitiesDiff(c *entities.Character) []AbilityDTO {
	abilitiesDTOs := make([]AbilityDTO, 0)
	characterAbilities := c.GetAbilities()

	newAbilityIds, _ := lo.Difference(
		lo.Keys(characterAbilities),
		c.CharacterLastTickState.Abilities,
	)

	if len(newAbilityIds) > 0 {
		for abilityId, ability := range characterAbilities {
			abilitiesDTOs = append(abilitiesDTOs, buildAbilityDTO(c, abilityId, ability))
		}
		c.CharacterLastTickState.Abilities = lo.Keys(characterAbilities)
	} else {
		for abilityId, ability := range characterAbilities {
			remainingCooldown := c.RemainingCooldown(ability)
			if remainingCooldown != c.CharacterLastTickState.AbilitiesRemainingCooldows[abilityId] {
				abilitiesDTOs = append(abilitiesDTOs, buildAbilityDTO(c, abilityId, ability))
			}
		}
	}
	slices.SortFunc(abilitiesDTOs, func(a, b AbilityDTO) int {
		return strings.Compare(a.Id, b.Id)
	})

	return abilitiesDTOs
}

// getInventoryDiff compares the current character inventory against the last tick
// and returns an InventoryDTO diff if any relevant changes are detected.
func getInventoryDiff(c *entities.Character) *InventoryDTO {
	currentInventory := c.GetInventory()

	if len(lo.Keys(currentInventory)) == 0 {
		return nil
	}

	currentlyEquipped := getEquippedMap(c)

	newItemIDs, removedItemIDs := lo.Difference(
		lo.Keys(currentInventory),
		lo.Keys(c.CharacterLastTickState.Items),
	)

	inventoryDTO := &InventoryDTO{}
	if len(newItemIDs) > 0 || len(removedItemIDs) > 0 || hasEquipmentChanged(c, currentlyEquipped) {
		inventoryDTO = buildFullInventoryDiff(c, currentInventory, currentlyEquipped, removedItemIDs)
	} else {
		inventoryDTO = buildIncrementalInventoryDiff(c, currentlyEquipped)
	}

	return inventoryDTO
}

// getEquippedMap returns a map of equipped item IDs for fast lookup.
func getEquippedMap(c *entities.Character) map[string]bool {
	m := make(map[string]bool, len(c.GetEquipped()))
	for _, equip := range c.GetEquipped() {
		m[equip.Id()] = true
	}
	return m
}

// hasEquipmentChanged checks if the equipped items differ from the last tick state.
func hasEquipmentChanged(c *entities.Character, currentlyEquipped map[string]bool) bool {
	for itemID := range c.GetInventory() {
		if currentlyEquipped[itemID] != c.CharacterLastTickState.EquippedItems[itemID] {
			return true
		}
	}
	return false
}

// buildFullInventoryDiff rebuilds the entire inventory state when items were added,
// removed, or equipment changed.
func buildFullInventoryDiff(
	c *entities.Character,
	currentInventory map[string]int64,
	currentlyEquipped map[string]bool,
	removedItemIDs []string,
) *InventoryDTO {
	dto := &InventoryDTO{}

	// Add all current items
	for itemID, quantity := range currentInventory {
		dto.Items = append(dto.Items, ItemDTO{
			Id:         itemID,
			Quantity:   quantity,
			IsEquipped: currentlyEquipped[itemID],
		})
	}

	// Update last tick state
	c.CharacterLastTickState.Items = lo.Associate(dto.Items, func(i ItemDTO) (string, int64) {
		return i.Id, i.Quantity
	})
	c.CharacterLastTickState.EquippedItems = currentlyEquipped

	// Add removed items (with quantity 0)
	for _, itemID := range removedItemIDs {
		dto.Items = append(dto.Items, ItemDTO{
			Id:         itemID,
			Quantity:   0,
			IsEquipped: false,
		})
	}

	return dto
}

// buildIncrementalInventoryDiff adds only items whose quantity has changed.
// Equipment changes are excluded here, since they’re caught earlier.
func buildIncrementalInventoryDiff(
	c *entities.Character,
	currentlyEquipped map[string]bool,
) *InventoryDTO {
	dto := &InventoryDTO{}

	for itemID, quantity := range c.GetInventory() {
		if quantity != c.CharacterLastTickState.Items[itemID] {
			dto.Items = append(dto.Items, ItemDTO{
				Id:         itemID,
				Quantity:   quantity,
				IsEquipped: currentlyEquipped[itemID],
			})
		}
	}

	if len(dto.Items) == 0 {
		return nil
	}

	// Update last tick state
	c.CharacterLastTickState.Items = lo.Associate(dto.Items, func(i ItemDTO) (string, int64) {
		return i.Id, i.Quantity
	})
	c.CharacterLastTickState.EquippedItems = currentlyEquipped

	return dto
}

func getStatsDiff(c *entities.Character) *map[string]int64 {
	diff := make(map[string]int64)
	for stat, value := range c.GetStats() {
		if c.CharacterLastTickState.Stats[stat] != value {
			diff[stat] = value
			c.CharacterLastTickState.Stats[stat] = value
		}
	}
	if len(diff) > 0 {
		return &diff
	}
	return nil
}

func buildAbilityDTO(
	c *entities.Character,
	abilityId string,
	ability *entities.Ability,
) AbilityDTO {
	remainingCooldown := c.RemainingCooldown(ability)
	abilityDTO := AbilityToDTO(*ability)
	abilityDTO.RemainingCooldown = float64(remainingCooldown) / 1000
	c.CharacterLastTickState.AbilitiesRemainingCooldows[abilityId] = remainingCooldown
	return abilityDTO
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

func CharacterToDTO(character entities.Character) CharacterDTO {
	characterHealth := character.GetHealth()
	characterAbilities := make([]AbilityDTO, 0)
	for _, ability := range character.GetAbilities() {
		abilityDTO := AbilityToDTO(*ability)
		abilityDTO.RemainingCooldown = float64(character.RemainingCooldown(ability)) / 1000
		characterAbilities = append(characterAbilities, abilityDTO)
	}
	slices.SortFunc(characterAbilities, func(a, b AbilityDTO) int {
		return strings.Compare(a.Id, b.Id)
	})

	currentHealth := characterHealth.GetCurrentHealth()
	maxHealth := characterHealth.GetMaxHealth()
	radius := character.GetRadius()
	characterExecutingAction := character.GetExecutingAction()
	executingActionDirection := characterExecutingAction.Direction()
	actionType := characterExecutingAction.ActionType()
	characterStats := character.GetStats()

	characterInventory := InventoryDTO{}
	equipped := lo.Map(lo.Values(character.GetEquipped()), func(item *entities.Equipment, index int) string {
		return item.Id()
	})
	for item, quantity := range character.GetInventory() {
		characterInventory.Items = append(
			characterInventory.Items,
			ItemDTO{Id: item, Quantity: quantity, IsEquipped: lo.Contains(equipped, item)},
		)
	}

	return CharacterDTO{
		Id:            character.GetId(),
		Position:      PositionToDTO(character.GetPosition()),
		Radius:        &radius,
		MaxHealth:     &maxHealth,
		CurrentHealth: &currentHealth,
		Stats:         &characterStats,
		Action:        &actionType,
		Direction:     PositionToDTO(executingActionDirection),
		Abilities:     characterAbilities,
		Inventory:     &characterInventory,
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
