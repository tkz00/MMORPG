package entities

type EffectTrigger string

const (
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
