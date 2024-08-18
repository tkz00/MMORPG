package character

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
	targeting	Targeting
	Mechanic  	func(caster Player, params AbilityParameters)
}

func NewAbility(id string, name string, rangeValue float64, cooldown int64, targeting Targeting, mechanic func(caster Player, params AbilityParameters)) *Ability {
	return &Ability{
		id: id,
		name: name,
		rangeValue: rangeValue,
		cooldown: cooldown,
		targeting: targeting,
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

func (ability Ability) Targeting() Targeting {
	return ability.targeting
}

func (ability Ability) Cast(caster Player, params AbilityParameters) {
	ability.Mechanic(caster, params)
}

func (ability Ability) CreateAction(abilityInfo AbilityInfo) CastAbilityAction {
	var abilityParams AbilityParameters
	switch ability.targeting {
	case Target:
		targetId, _ := abilityInfo.GetTargetId()
		abilityParams = TargetIdAbilityParams{
			TargetId: targetId,
		}
	case Coordinates:
		targetPosition, _ := abilityInfo.GetTargetPosition()
		abilityParams = CoordinateAbilityParams{
			Target: targetPosition,
		}
	}
	return CastAbilityAction{
		ability: ability,
		params: abilityParams,
	}
}
