package entities

type Targeting int

const (
	Target Targeting = iota
	NoTarget
	Coordinates
)

type Ability struct {
	id             string
	name           string
	rangeValue     float64
	cooldown       int64
	targeting      Targeting
	characterState Action
	mechanics      []Mechanic
}

func NewAbility(
	id string,
	name string,
	rangeValue float64,
	cooldown int64,
	targeting Targeting,
	characterState Action,
	mechanics ...Mechanic) *Ability {
	return &Ability{
		id:             id,
		name:           name,
		rangeValue:     rangeValue,
		cooldown:       cooldown,
		targeting:      targeting,
		characterState: characterState,
		mechanics:      mechanics,
	}
}

func (ability Ability) Id() string {
	return ability.id
}

func (ability Ability) Name() string {
	return ability.name
}

func (ability Ability) Range() float64 {
	return ability.rangeValue
}

func (ability Ability) Cooldown() int64 {
	return ability.cooldown
}

func (ability Ability) Targeting() Targeting {
	return ability.targeting
}

func (ability Ability) CharacterState() Action {
	return ability.characterState
}

func (ability Ability) Mechanics() []Mechanic {
	return ability.mechanics
}

func (a *Ability) Clone() *Ability {
	newAbility := *a // Start with a shallow copy

	// Deep copy the mechanics slice
	if a.mechanics != nil {
		newAbility.mechanics = make([]Mechanic, len(a.mechanics))
		for i, mechanic := range a.mechanics {
			newAbility.mechanics[i] = mechanic.Clone() // Use the deep copy method
		}
	}

	return &newAbility
}
