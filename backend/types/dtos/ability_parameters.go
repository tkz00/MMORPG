package dtos

import (
	"encoding/json"
	"fmt"
	"unnamed-mmo/backend/types"
)

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

func (abilityCast AbilityCastDTO) GetTargetPosition() (types.Position, error) {
    targetPositionMap, ok := abilityCast.AbilityParameters[TargetPosition].(map[string]interface{})
    if !ok {
        return types.Position{}, fmt.Errorf("unable to cast TargetPosition to map[string]interface{}")
    }

    xValue, xOk := targetPositionMap["x"]
    zValue, zOk := targetPositionMap["z"]

    if !xOk || !zOk {
        return types.Position{}, fmt.Errorf("x or z value not found in the target position map")
    }

    x, xConvOk := xValue.(float64)
    z, zConvOk := zValue.(float64)

    if !xConvOk || !zConvOk {
        return types.Position{}, fmt.Errorf("x or z value could not be converted to float64")
    }

    return *types.NewPosition(x, z), nil
}

func (abilityCast AbilityCastDTO) GetTargetId() (string, error) {
	targetId, ok := abilityCast.AbilityParameters[TargetId].(string)
    if !ok {
        return "", fmt.Errorf("unable to cast TargetId to string")
    }

    return targetId, nil
}

func (abilityCast AbilityCastDTO) GetName() string {
	return abilityCast.Name
}
