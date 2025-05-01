package dtos

import (
	"backend/pkg/game/entities"
	"backend/pkg/utils"
	"encoding/json"
	"fmt"
)

type DTO interface {
	GetType() string
}

type GameStateDTO struct {
	Players     []CharacterDTO  `json:"players,omitempty"`
	Projectiles []ProjectileDTO `json:"projectiles,omitempty"`
	AreaEffects []AoEDTO        `json:"area_effects,omitempty"`
	Npcs        []CharacterDTO  `json:"npcs,omitempty"`
}

func (g GameStateDTO) GetType() string {
	return "GameState"
}

type CharacterDTO struct {
	Id            string           `json:"id"`
	MaxHealth     *int             `json:"maxHealth,omitempty"`
	CurrentHealth *int             `json:"currentHealth,omitempty"`
	Radius        *float64         `json:"radius,omitempty"`
	Position      *PositionDTO     `json:"position,omitempty"`
	Action        *entities.Action `json:"action,omitempty"`
	Direction     *PositionDTO     `json:"direction,omitempty"`
	Abilities     []AbilityDTO     `json:"abilities,omitempty"`
	Inventory     InventoryDTO     `json:"inventory"` // Refactor this
}

func (p CharacterDTO) GetType() string {
	return "Player"
}

type PositionDTO struct {
	X float64 `json:"x"`
	Z float64 `json:"z"`
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
	return "position"
}

type ProjectileDTO struct {
	Id       string       `json:"id"`
	Caster   string       `json:"caster,omitempty"`
	Position *PositionDTO `json:"position,omitempty"`
	Radius   *float64     `json:"radius,omitempty"`
	State    string       `json:"state,omitempty"`
}

func (p ProjectileDTO) GetType() string {
	return "Projectile"
}

type AoEDTO struct {
	Id       string       `json:"id"`
	Position *PositionDTO `json:"position,omitempty"`
	Radius   *float64     `json:"radius,omitempty"`
}

func (AoEDTO) GetType() string {
	return "AoE"
}

type InventoryDTO struct {
	Items []entities.ItemChange `json:"items"`
}

type UseItemDTO struct {
	ItemId   string `json:"item_id"`
	TargetId string `json:"target_id"`
}

func (u UseItemDTO) GetType() string {
	return "use_item"
}

type AbilityParameters int

const (
	TargetPosition AbilityParameters = iota
	TargetId
)

type AbilityCastDTO struct {
	Id                string                            `json:"id"`
	AbilityParameters map[AbilityParameters]interface{} `json:"abilityParameters"`
}

func (p AbilityCastDTO) GetType() string {
	return "ability_cast"
}

var stringToAbilityParameters = map[string]AbilityParameters{
	"TargetPosition": TargetPosition,
	"TargetId":       TargetId,
}

func (a *AbilityCastDTO) UnmarshalJSON(data []byte) error {
	type Alias AbilityCastDTO
	aux := &struct {
		AbilityParameters map[string]interface{} `json:"abilityParameters"`
		*Alias
	}{
		Alias: (*Alias)(a),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	a.AbilityParameters = make(map[AbilityParameters]interface{})
	for k, v := range aux.AbilityParameters {
		if param, found := stringToAbilityParameters[k]; found {
			a.AbilityParameters[param] = v
		} else {
			return fmt.Errorf("unknown ability parameter: %s", k)
		}
	}

	return nil
}
