package types

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
}

func NewAbility(id string, name string, rangeValue float64) *Ability {
	return &Ability{
		id: id,
		name: name,
		rangeValue: rangeValue,
	}
}

func (ability Ability) GetId() string {
	return ability.id
}

func (ability Ability) GetName() string {
	return ability.name
}

func (ability Ability) GetRange() float64 {
	return ability.rangeValue
}
