package types

import (
	"encoding/json"
	"fmt"
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

func extractTargetPosition(abilityParameters map[AbilityParameters]interface{}) (Position, error) {
    targetPositionMap, ok := abilityParameters[TargetPosition].(map[string]interface{})
    if !ok {
        return Position{}, fmt.Errorf("unable to cast TargetPosition to map[string]interface{}")
    }

    xValue, xOk := targetPositionMap["x"]
    zValue, zOk := targetPositionMap["z"]

    if !xOk || !zOk {
        return Position{}, fmt.Errorf("x or z value not found in the target position map")
    }

    x, xConvOk := xValue.(float64)
    z, zConvOk := zValue.(float64)

    if !xConvOk || !zConvOk {
        return Position{}, fmt.Errorf("x or z value could not be converted to float64")
    }

    return Position{
        x: float32(x),
        z: float32(z),
    }, nil
}
