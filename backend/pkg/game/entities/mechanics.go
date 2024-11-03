package entities

import "fmt"

type Mechanic struct {
	MechanicType string                 // Type of effect (e.g., "heal", "unlock")
	Params       map[string]interface{} // Dynamic parameters for the effect (e.g., healing amount, target)
}

type MechanicHandler func(caster Character, gs GameState, params map[string]interface{}) (GameState, error)

var mechanicHandlers = map[string]MechanicHandler{}

func RegisterMechanicHandler(mechanicType string, handler MechanicHandler) {
	mechanicHandlers[mechanicType] = handler
}

func HealMechanic(caster Character, gs GameState, params map[string]interface{}) (GameState, error) {
	if amount, ok := params["amount"].(int); ok {
		gs.players[caster.id].HealthVariation(amount)
		// caster.HealthVariation(amount)
		return gs, nil
	}
	return gs, fmt.Errorf("missing or invalid 'amount' parameter")
}

func DamageMechanic(caster Character, gs GameState, params map[string]interface{}) (GameState, error) {
	if amount, ok := params["amount"].(int); ok {
		caster.HealthVariation(-amount)
		return gs, nil
	}
	return gs, fmt.Errorf("missing or invalid 'amount' parameter")
}
