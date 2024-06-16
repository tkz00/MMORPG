package types

import (
	"encoding/json"
	"fmt"
	"unnamed-mmo/backend/utils"
)

type DTO interface {
	GetType() string
}

type GameStateDTO struct {
	Players 	[]PlayerDTO `json:"players"`
	Projectiles []ProjectileDTO `json:"projectiles"`
}

func (g GameStateDTO) GetType() string {
	return "GameState"
}

type PlayerDTO struct {
	Id       		string      `json:"id"`
	MaxHealth 		int			`json:"maxHealth"`	
	CurrentHealth 	int			`json:"currentHealth"`	
	Position 		PositionDTO `json:"position"`
}

func (p PlayerDTO) GetType() string {
	return "Player"
}

type PositionDTO struct {
	X float32 `json:"x"`
	Z float32 `json:"z"`
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
	return "Position"
}

type AbilityParameters int

const (
	TargetPosition AbilityParameters = iota
	TargetId
)

type AbilityCastDTO struct {
	Name             	string                           	`json:"name"`
	AbilityParameters 	map[AbilityParameters]interface{} 	`json:"abilityParameters"`
}

func (p AbilityCastDTO) GetType() string {
	return "AbilityCast"
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

type ProjectileDTO struct {
	Id			string 		`json:"id"`
	Caster 		string 		`json:"caster"`
	Position 	PositionDTO `json:"position"`
	Damage		int 		`json:"damage"`
}

func (p ProjectileDTO) GetType() string {
	return "Projectile"
}
