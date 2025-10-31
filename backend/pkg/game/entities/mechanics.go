package entities

import (
	"backend/pkg/utils"
	"fmt"
	"time"
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
	targetID, ok := params["target_id"].(string)
	if !ok {
		fmt.Println("Warning: 'target_id' not set or invalid.")
		return nil
	}

	target, err := gs.GetCharacterById(targetID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if !target.IsAlive() {
		return nil
	}

	healAmount := 0.0
	hasHealSource := false

	if _, hasBaseAmount := params["base_amount"]; hasBaseAmount {
		fmt.Println("parse error")
	}

	if base_amount, ok := params["base_amount"].(float64); ok {
		healAmount += base_amount
		hasHealSource = true
	}

	if damageScaling, ok := params["damage_stat_multiplier"].(float64); ok {
		healAmount += float64(caster.GetStat("damage")) * damageScaling
		hasHealSource = true
	}

	if !hasHealSource {
		fmt.Println(
			"Warning: No valid heal source ('base_amount' or 'damage_stat_multiplier') provided. No heal will be done.",
		)
	}

	target.Heal(int(healAmount))

	if npc, ok := gs.npcs[targetID]; ok {
		npc.BecomeAggressive(caster)
	}

	if !target.IsAlive() {
		caster, _ := gs.GetCharacterById(caster.id)
		caster.Loot(target.Inventory)
	}

	if onHitMechanics, ok := params["on_hit_mechanics"]; ok {
		for _, mechanic := range onHitMechanics.([]Mechanic) {
			mechanic.Params["target_id"] = target.id
		}
		resolveMechanics(caster.id, gs, onHitMechanics.([]Mechanic))
	}

	return nil
}

func DamageMechanic(caster *Character, gs *GameState, params map[string]interface{}) error {
	targetID, ok := params["target_id"].(string)
	if !ok {
		fmt.Println("Warning: 'target_id' not set or invalid.")
		return nil
	}

	target, err := gs.GetCharacterById(targetID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if !target.IsAlive() {
		return nil
	}

	damageAmount := 0.0
	hasDamageSource := false

	if base_amount, ok := params["base_amount"].(float64); ok {
		damageAmount += base_amount
		hasDamageSource = true
	}

	if damageScaling, ok := params["damage_stat_multiplier"].(float64); ok {
		damageAmount += float64(caster.GetStat("damage")) * damageScaling
		hasDamageSource = true
	}

	if !hasDamageSource {
		fmt.Println(
			"Warning: No valid damage source ('base_amount' or 'damage_stat_multiplier') provided. No damage will be dealt.",
		)
	}

	target.TakeDamage(int(damageAmount))

	caster.TriggerEffects(EffectOnDamageDealt, map[string]interface{}{
		"damage":      damageAmount,
		"damage_type": "TODO", // Add this later
		"target":      target,
	}, gs)

	if npc, ok := gs.npcs[targetID]; ok {
		npc.BecomeAggressive(caster)
	}

	if !target.IsAlive() {
		caster, _ := gs.GetCharacterById(caster.id)
		caster.Loot(target.Inventory)
	}

	if onHitMechanics, ok := params["on_hit_mechanics"]; ok {
		for _, mechanic := range onHitMechanics.([]Mechanic) {
			mechanic.Params["target_id"] = target.id
		}
		resolveMechanics(caster.id, gs, onHitMechanics.([]Mechanic))
	}

	return nil
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
	if numberOfProjectiles, ok := params["number"].(float64); ok {
		for i := 0; i < int(numberOfProjectiles); i++ {
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

func BuffStatMechanic(caster *Character, gs *GameState, params map[string]interface{}) error {
	targetID, ok := params["target_id"].(string)
	if !ok {
		fmt.Println("Warning: 'target_id' not set or invalid.")
		return nil
	}

	target, err := gs.GetCharacterById(targetID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if !target.IsAlive() {
		return nil
	}

	statName := params["target_stat"].(string)
	flatValue := int64(0)
	percentValue := float64(0)
	modifiedStat := false

	if statMultiplier, ok := params["multiplier"].(float64); ok {
		percentValue = statMultiplier
		modifiedStat = true
	}

	if base_amount, ok := params["base_amount"].(int64); ok {
		flatValue = base_amount
		modifiedStat = true
	}

	if !modifiedStat {
		fmt.Println(
			"Warning: No valid modification value ('base_amount' or 'multiplier') provided. No change to stat will be made.",
		)
		return nil
	}

	// Get duration in milliseconds
	durationMs, ok := params["duration_ms"].(float64)
	if !ok {
		fmt.Println("Warning: 'duration_ms' not provided for temporary buff. Buff will be permanent.")
		durationMs = 0
	}

	// Generate a unique ID for this modification
	modID := fmt.Sprintf("buff_%s_%s_%d", caster.id, statName, time.Now().UnixNano())

	var expiresAt *time.Time
	if durationMs > 0 {
		expirationTime := time.Now().Add(time.Duration(durationMs) * time.Millisecond)
		expiresAt = &expirationTime
	}

	mod := StatModification{
		ID:           modID,
		StatName:     statName,
		FlatValue:    flatValue,
		PercentValue: percentValue,
		Source:       "buff",
		ExpiresAt:    expiresAt,
	}

	target.AddStatModification(mod)

	if onHitMechanics, ok := params["on_hit_mechanics"]; ok {
		for _, mechanic := range onHitMechanics.([]Mechanic) {
			mechanic.Params["target_id"] = target.id
		}
		resolveMechanics(caster.id, gs, onHitMechanics.([]Mechanic))
	}

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
	case "buff_stat":
		switch params["targeting_strategy"] {
		case "character_hit":
		case "caster":
			params["target_id"] = casterId
		default:
			panic("no targeting_strategy for buff_stat mechanic found")
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

			// gets the target coordinates for projectiles when they're shot in an arc, depending on the radius of the arc and number of projectiles to spawn, these are equally distributed around the radius
			for i := range int(params["number"].(float64)) {
				params[fmt.Sprint("target_coordinates_", i)] = utils.CalculateNewPosition(
					// params["projectile_last_position"].(utils.Vector2),
					params["initial_coordinates"].(utils.Vector2),
					params["range"].(float64),
					params["radius"].(float64)*float64(i)/params["number"].(float64),
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

func (m Mechanic) Clone() Mechanic {
	cloned := m
	cloned.Params = deepCloneParams(m.Params)
	return cloned
}

func deepCloneParams(params map[string]interface{}) map[string]interface{} {
	if params == nil {
		return nil
	}
	cloned := make(map[string]interface{}, len(params))
	for k, v := range params {
		switch val := v.(type) {
		case []Mechanic:
			// Clone []Mechanic recursively
			newSlice := make([]Mechanic, len(val))
			for i, mech := range val {
				newSlice[i] = mech.Clone()
			}
			cloned[k] = newSlice
		case map[string]interface{}:
			// Recursively deep copy nested maps
			cloned[k] = deepCloneParams(val)
		default:
			cloned[k] = val // Primitives and unhandled types
		}
	}
	return cloned
}
