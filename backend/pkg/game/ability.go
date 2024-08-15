package game

type Targeting int

const (
	Target Targeting = iota
	NoTarget
	Coordinates
)

type Ability struct {
	id			string
	name		string
	rangeValue	float64
	cooldown	int64
	Mechanic  	func(caster Player, params AbilityParameters)
}

func NewAbility(id string, name string, rangeValue float64, cooldown int64, mechanic func(caster Player, params AbilityParameters)) *Ability {
	return &Ability{
		id: id,
		name: name,
		rangeValue: rangeValue,
		cooldown: cooldown,
		Mechanic: mechanic,
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

func (ability Ability) Cast(caster Player, params AbilityParameters) {
	ability.Mechanic(caster, params)
}
