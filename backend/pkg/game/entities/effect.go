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
