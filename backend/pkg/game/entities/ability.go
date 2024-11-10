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
}

func NewAbility(
	id string,
	name string,
	rangeValue float64,
	cooldown int64,
	targeting Targeting,
	characterState Action) *Ability {
	return &Ability{
		id:             id,
		name:           name,
		rangeValue:     rangeValue,
		cooldown:       cooldown,
		targeting:      targeting,
		characterState: characterState,
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
