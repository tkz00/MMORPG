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
}

func NewAbility(id string, name string, rangeValue float64, cooldown int64) *Ability {
	return &Ability{
		id: id,
		name: name,
		rangeValue: rangeValue,
		cooldown: cooldown,
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
