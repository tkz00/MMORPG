package entities

import "fmt"

type EffectTrigger string

const (
	EffectPassive       EffectTrigger = "passive"
	EffectOnKill        EffectTrigger = "on_kill"
	EffectOnDamageDealt EffectTrigger = "on_damage_dealt"
)

// Effects are passive modifications that characters can have, affecting the character, their stats, or their abilities.
type Effect struct {
	id        string
	name      string
	trigger   EffectTrigger
	mechanics []Mechanic
}

func (e *Effect) Id() string {
	return e.id
}

func (e *Effect) Name() string {
	return e.name
}

func (e *Effect) Trigger() EffectTrigger {
	return e.trigger
}

func (e *Effect) Mechanics() []Mechanic {
	return e.mechanics
}

func CreateEffect(
	id string,
	name string,
	trigger EffectTrigger,
	mechanics []Mechanic,
) *Effect {
	return &Effect{
		id:        id,
		name:      name,
		trigger:   trigger,
		mechanics: mechanics,
	}
}

func (c *Character) AddEffect(effectId string, gs *GameState) {
	c.effects = append(c.effects, effectId)

	// Apply passive effects immediately through the normal mechanic flow
	if effect, ok := ExistingEffects[effectId]; ok && effect.trigger == EffectPassive {
		for i := range effect.mechanics {
			effect.mechanics[i].Params["source_effect_id"] = effectId
		}
		resolveMechanics(c.id, gs, effect.mechanics)
		return
	}
	fmt.Printf("Error adding effect, not found in existing effects, effect_id: %s\n", effectId)
}

func (c *Character) TriggerEffects(trigger EffectTrigger, eventParams map[string]interface{}, gs *GameState) {
	for _, effectId := range c.effects {
		if _, ok := ExistingEffects[effectId]; !ok {
			fmt.Printf("Error, tried to execute effect %s, not found in existing effects\n", effectId)
		}
		if ExistingEffects[effectId].trigger == trigger {
			// Inject event parameters and source effect ID into each mechanic before resolving
			for i := range ExistingEffects[effectId].mechanics {
				ExistingEffects[effectId].mechanics[i].Params["event_params"] = eventParams
				ExistingEffects[effectId].mechanics[i].Params["source_effect_id"] = effectId
			}
			// Reuse the existing mechanic resolution logic
			resolveMechanics(c.id, gs, ExistingEffects[effectId].mechanics)
		}
	}
}

// move outside of character_seeds
// this should be loaded from DB and cached
var ExistingEffects = map[string]Effect{
	"spell_vampirism": NewSpellVampirismEffect(0.1),
	"iron_will":       NewStatBoostOnKillEffect("iron_will", "Iron Will", "defense", 0, 0.1),
	"cursed_blade":    NewPassiveStatBoostEffect("cursed_blade", "Cursed Blade", "damage", 20, 0),
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

// NewStatBoostOnKillEffect is an effect that permanently increases a character's stat when it kills a character.
// Can provide either flatBonus (additive) or percentBonus (multiplicative) or both.
func NewStatBoostOnKillEffect(id, name, statName string, flatBonus int64, percentBonus float64) Effect {
	return Effect{
		id:      id,
		name:    name,
		trigger: EffectOnKill,
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
