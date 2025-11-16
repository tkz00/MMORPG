package entities

type EffectTrigger string

const (
	EffectPassive       EffectTrigger = "passive"
	EffectOnDamageDealt EffectTrigger = "on_damage_dealt"
)

// Effects are passive modifications that characters can have, affecting the character, their stats, or their abilities.
type Effect struct {
	id         string
	name       string
	trigger    EffectTrigger
	mechanics  []Mechanic
	parameters map[string]interface{}
}

func (c *Character) AddEffect(effect Effect, gs *GameState) {
	c.effects = append(c.effects, effect)

	// Apply passive effects immediately through the normal mechanic flow
	if effect.trigger == EffectPassive {
		resolveMechanics(c.id, gs, effect.mechanics)
	}
}

func (c *Character) TriggerEffects(trigger EffectTrigger, eventParams map[string]interface{}, gs *GameState) {
	for _, effect := range c.effects {
		if effect.trigger == trigger {
			// Inject event parameters into each mechanic before resolving
			for i := range effect.mechanics {
				effect.mechanics[i].Params["event_params"] = eventParams
			}
			// Reuse the existing mechanic resolution logic
			resolveMechanics(c.id, gs, effect.mechanics)
		}
	}
}

// NewSpellVampirismEffect creates a lifesteal effect that heals the caster for a percentage of damage dealt.
func NewSpellVampirismEffect(healPercentage float64) Effect {
	return Effect{
		id:      "spell_vampirism",
		name:    "Spell Vampirism",
		trigger: EffectOnDamageDealt,
		mechanics: []Mechanic{
			{
				MechanicType: "heal",
				Params: map[string]interface{}{
					"targeting_strategy": "caster",
					"scaling_from_event": map[string]interface{}{
						"event_field": "damage",
						"factor":      healPercentage,
					},
				},
			},
		},
	}
}

// NewPassiveStatBoostEffect creates a passive effect that permanently increases a character's stat.
// Can provide either flatBonus (additive) or percentBonus (multiplicative) or both.
func NewPassiveStatBoostEffect(id, name, statName string, flatBonus int64, percentBonus float64) Effect {
	return Effect{
		id:      id,
		name:    name,
		trigger: EffectPassive,
		mechanics: []Mechanic{
			{
				MechanicType: "buff_stat",
				Params: map[string]interface{}{
					"targeting_strategy": "caster",
					"target_stat":        statName,
					"base_amount":        flatBonus,
					"multiplier":         percentBonus,
				},
			},
		},
	}
}
