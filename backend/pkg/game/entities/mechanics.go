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
	if numberOfProjectiles, ok := params["number"].(int); ok {
		for i := 0; i < numberOfProjectiles; i++ {
			if targetPosition, ok := params[fmt.Sprint("target_coordinates_", i)].(utils.Vector2); ok {
				onHitMechanics, _ := params["on_hit_mechanics"].([]Mechanic)
				newProjectile := CreateProjectile(
					params[fmt.Sprint("initial_coordinates_", i)].(utils.Vector2),
					targetPosition,
					params["range"].(float64),
					caster.id,
					onHitMechanics,
				)
				gs.projectiles[newProjectile.id] = newProjectile
			} else {
				return fmt.Errorf("missing or invalid 'targetCoordinates' parameter")
			}
		}
		return nil
	} else { // no number of projectiles specified, so 1 projectile is assumed
		if targetPosition, ok := params["target_coordinates"].(utils.Vector2); ok {
			onHitMechanics, _ := params["on_hit_mechanics"].([]Mechanic)
			newProjectile := CreateProjectile(
				params["initial_coordinates"].(utils.Vector2),
				targetPosition,
				params["range"].(float64),
				caster.id,
				onHitMechanics,
			)
			gs.projectiles[newProjectile.id] = newProjectile
			return nil
		}
	}
	return fmt.Errorf("missing or invalid 'targetCoordinates' parameter")
}

// For now this designates the target of a mechanic depending on it's targeting strategy
func resolveParameters(
	params map[string]interface{},
	casterId string,
	targetId string,
	gs *GameState, // currently unused but will be in the future
) {
	if params["origin_position"] == "target" {
		target, _ := gs.GetCharacterById(targetId)
		params["initial_coordinates"] = target.position
	} else {
		caster, _ := gs.GetCharacterById(casterId)
		params["initial_coordinates"] = caster.position
	}

	switch params["targeting_strategy"] {
	case "caster":
		params["targetId"] = casterId
	case "character_hit":
		params["targetId"] = targetId
	case "arc":
		baseCharacter := &Character{}
		if params["origin_position"] == "target" {
			target, _ := gs.GetCharacterById(targetId)
			baseCharacter = target
		} else {
			caster, _ := gs.GetCharacterById(casterId)
			baseCharacter = caster
		}
		// get's the target coordinates for projectiles when they're shot in an arc, depending on the radius of the arc and number of projectiles to spawn, these are equally distributed around the radius
		for i := 0; i < params["number"].(int); i++ {
			params[fmt.Sprint("target_coordinates_", i)] = utils.CalculateNewPosition(
				params["projectile_last_position"].(utils.Vector2),
				params["range"].(float64),
				params["radius"].(float64)*float64(i)/float64(params["number"].(int)),
			)
			params[fmt.Sprint("initial_coordinates_", i)] = utils.ClosestPositionInRange(
				params[fmt.Sprint("target_coordinates_", i)].(utils.Vector2),
				baseCharacter.position,
				(baseCharacter.GetRadius() + 0.5),
			)
		}
	default:
		panic("no targeting_strategy for projectile hit mechanic found")
	}
}
