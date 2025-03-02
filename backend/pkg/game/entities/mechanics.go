package entities

import (
	"fmt"
	"tkz00/backend/pkg/utils"
)

type Mechanic struct {
	MechanicType string                 `json:"mechanic_type"` // Type of effect (e.g., "heal", "unlock")
	Params       map[string]interface{} `json:"params"`        // Dynamic parameters for the effect (e.g., healing amount, target)
}

type MechanicHandler func(caster *Character, gs *GameState, params map[string]interface{}) error

var mechanicHandlers = map[string]MechanicHandler{}

func RegisterMechanicHandler(mechanicType string, handler MechanicHandler) {
	mechanicHandlers[mechanicType] = handler
}

func HealMechanic(caster *Character, gs *GameState, params map[string]interface{}) error {
	if amount, ok := params["amount"].(float64); ok {
		target, err := gs.GetCharacterById(params["target_id"].(string))
		if err != nil {
			fmt.Println(err)
			return err
		}
		target.HealthVariation(int(amount))

		if onHitMechanics, ok := params["on_hit_mechanics"]; ok {
			for _, mechanic := range onHitMechanics.([]Mechanic) {
				mechanic.Params["target_id"] = target.id
			}
			resolveMechanics(caster.id, gs, onHitMechanics.([]Mechanic))
		}

		return nil
	}
	return fmt.Errorf("missing or invalid 'amount' parameter")
}

func DamageMechanic(caster *Character, gs *GameState, params map[string]interface{}) error {
	if amount, ok := params["amount"].(float64); ok {
		target, err := gs.GetCharacterById(params["target_id"].(string))
		if err != nil {
			fmt.Println(err)
			return err
		}
		target.HealthVariation(-int(amount))

		if npc, ok := gs.npcs[params["target_id"].(string)]; ok {
			npc.BecomeAggressive(caster)
		}

		if !target.IsAlive() {
			gs.Players()[caster.id].Loot(target.Inventory)
		}

		if onHitMechanics, ok := params["on_hit_mechanics"]; ok {
			for _, mechanic := range onHitMechanics.([]Mechanic) {
				mechanic.Params["target_id"] = target.id
			}
			resolveMechanics(caster.id, gs, onHitMechanics.([]Mechanic))
		}

		return nil
	}
	return fmt.Errorf("missing or invalid 'amount' parameter")
}

func DelayMechanic(caster *Character, gs *GameState, params map[string]interface{}) error {
	if delayMechanics, ok := params["execute_after_delay_mechanics"]; ok {
		if delayMs, ok := params["delay_ms"]; ok {
			gs.DelayMechanics(delayMechanics.([]Mechanic), int(delayMs.(float64)), caster.id)
			return nil
		}
		return fmt.Errorf("missing 'delay_ms' parameter")
	}
	return fmt.Errorf("missing 'execute_after_delay_mechanics' parameter")
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

func AoEMechanic(
	caster *Character,
	gs *GameState,
	params map[string]interface{},
) error {
	AoE := InstantiateAoE(
		params["target_coordinates"].(utils.Vector2),
		params["radius"].(float64),
		int(params["duration_ms"].(float64)),
		caster.id,
		params["on_hit_mechanics"].([]Mechanic),
	)

	gs.AddAreaEffect(AoE)

	return nil
}

func resolveMechanics(
	casterId string,
	gs *GameState,
	mechanics []Mechanic,
) {
	for _, mechanic := range mechanics {
		if handler, exists := mechanicHandlers[mechanic.MechanicType]; exists {
			resolveParameters(
				mechanic,
				casterId,
				gs,
			)
			caster, err := gs.GetCharacterById(casterId)
			if err != nil {
				fmt.Println(err)
			}
			if err := handler(caster, gs, mechanic.Params); err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Printf("no handler found for effect type: %s\n", mechanic.MechanicType)
		}
	}
}

// For now this designates the target of a mechanic depending on it's targeting strategy
func resolveParameters(
	mechanic Mechanic,
	casterId string,
	gs *GameState,
) {
	params := mechanic.Params

	switch mechanic.MechanicType {
	case "damage":
		switch params["targeting_strategy"] {
		case "character_hit":
		case "caster":
			params["target_id"] = casterId
		default:
			panic("no targeting_strategy for damage mechanic found")
		}
	case "heal":
		switch params["targeting_strategy"] {
		case "character_hit":
		case "caster":
			params["target_id"] = casterId
		default:
			panic("no targeting_strategy for heal mechanic found")
		}
	case "delay":
		for _, delayedMechanic := range params["execute_after_delay_mechanics"].([]Mechanic) {
			delayedMechanic.Params["target_id"] = params["target_id"]
		}
	case "create_projectile":
		switch params["targeting_strategy"] { // targeting_strategy is a very bad name for this, it should be the shape the projectile(s) spawn, arc or line or whatever
		case "arc":
			baseCharacter := &Character{}
			if params["origin_position"] == "target" {
				baseCharacter, _ = gs.GetCharacterById(params["target_id"].(string))
			} else {
				baseCharacter, _ = gs.GetCharacterById(casterId)
			}

			params["initial_coordinates"] = baseCharacter.position

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
					(baseCharacter.GetRadius() + 0.5), // fix this to be the radius of the character + radius of projectile + epsilon)
				)
			}
		default:
			panic("no targeting_strategy for create_projectile hit mechanic found")
		}
	case "create_AoE":
		if params["target_coordinates"] == nil && params["target_id"] != nil {
			targetCharacter, _ := gs.GetCharacterById(params["target_id"].(string))
			params["target_coordinates"] = targetCharacter.position
		}
	default:
	}
}

func deepCopyMap(original map[string]interface{}) map[string]interface{} {
	newMap := make(map[string]interface{}, len(original))
	for k, v := range original {
		// Handle nested maps
		if nestedMap, ok := v.(map[string]interface{}); ok {
			newMap[k] = deepCopyMap(nestedMap) // Recursively copy
		} else {
			newMap[k] = v // Copy other values as-is
		}
	}
	return newMap
}

func (m *Mechanic) Clone() Mechanic {
	newMechanic := *m
	if m.Params != nil {
		newMechanic.Params = deepCopyMap(m.Params) // Use deep copy for maps
	}
	return newMechanic
}
