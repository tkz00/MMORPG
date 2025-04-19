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
	Players           []CharacterDTO  `json:"players"`
	Projectiles       []ProjectileDTO `json:"projectiles"`
	AreaEffects       []AoEDTO        `json:"area_effects"`
	Npcs              []CharacterDTO  `json:"npcs"`
	EntitiesToDestroy []string        `json:"entities_to_destroy"`
}

func (g GameStateDTO) GetType() string {
	return "GameState"
}

type CharacterDTO struct {
	Id              string             `json:"id"`
	MaxHealth       int                `json:"maxHealth"`
	CurrentHealth   int                `json:"currentHealth"`
	Radius          float64            `json:"radius"`
	Position        PositionDTO        `json:"position"`
	ExecutingAction ExecutingActionDTO `json:"executingAction"`
	Abilities       []AbilityDTO       `json:"abilities"`
	Inventory       InventoryDTO       `json:"inventory"`
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
	Id       string      `json:"id"`
	Caster   string      `json:"caster"`
	Position PositionDTO `json:"position"`
	Radius   float64     `json:"radius"`
	State    string      `json:"state"`
}

func (p ProjectileDTO) GetType() string {
	return "Projectile"
}

type AoEDTO struct {
	Id       string      `json:"id"`
	Caster   string      `json:"caster"`
	Position PositionDTO `json:"position"`
	Radius   float64     `json:"radius"`
}

func (AoEDTO) GetType() string {
	return "AoE"
}

type ExecutingActionDTO struct {
	Action    entities.Action `json:"action"`
	Direction PositionDTO     `json:"direction"`
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
