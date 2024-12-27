package entities

import (
	"fmt"
	"tkz00/backend/pkg/utils"
)

type Mechanic struct {
	MechanicType string                 // Type of effect (e.g., "heal", "unlock")
	Params       map[string]interface{} // Dynamic parameters for the effect (e.g., healing amount, target)
}

type MechanicHandler func(caster *Character, gs *GameState, params map[string]interface{}) error

var mechanicHandlers = map[string]MechanicHandler{}

func RegisterMechanicHandler(mechanicType string, handler MechanicHandler) {
	mechanicHandlers[mechanicType] = handler
}

func HealMechanic(caster *Character, gs *GameState, params map[string]interface{}) error {
	if amount, ok := params["amount"].(int); ok {
		target, err := gs.GetCharacterById(params["targetId"].(string))
		if err != nil {
			fmt.Println(err)
			return err
		}
		target.HealthVariation(amount)

		return nil
	}
	return fmt.Errorf("missing or invalid 'amount' parameter")
}

func DamageMechanic(caster *Character, gs *GameState, params map[string]interface{}) error {
	if amount, ok := params["amount"].(int); ok {
		target, err := gs.GetCharacterById(params["targetId"].(string))
		if err != nil {
			fmt.Println(err)
			return err
		}
		target.HealthVariation(-amount)

		return nil
	}
	return fmt.Errorf("missing or invalid 'amount' parameter")
}

func CreateProjectileMechanic(
	caster *Character,
	gs *GameState,
	params map[string]interface{},
) error {
	if targetPosition, ok := params["target_coordinates"].(utils.Vector2); ok {
		onHitMechanics, _ := params["on_hit_mechanics"].([]Mechanic)
		newProjectile := CreateProjectile(
			caster.position,
			targetPosition,
			params["range"].(float64),
			caster.id,
			onHitMechanics,
		)
		gs.projectiles[newProjectile.id] = newProjectile
		return nil
	}
	return fmt.Errorf("missing or invalid 'targetCoordinates' parameter")
}
