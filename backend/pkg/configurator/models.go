package configurator

import "backend/pkg/game/entities"

type ConfiguratorAbility struct {
	ID             string              `json:"id"`
	Name           string              `json:"name"`
	RangeValue     *float64            `json:"range"`
	Cooldown       *int64              `json:"cooldown"`
	Targeting      *entities.Targeting `json:"targeting"`
	CharacterState *entities.Action    `json:"character_state"`
	Mechanics      []entities.Mechanic `json:"mechanics"`
}

func ConvertToConfiguratorAbility(ability entities.Ability) ConfiguratorAbility {
	return ConfiguratorAbility{
		ID:             ability.Id(),
		Name:           ability.Name(),
		RangeValue:     ptr(ability.Range()),
		Cooldown:       ptr(ability.Cooldown()),
		Targeting:      ptr(ability.Targeting()),
		CharacterState: ptr(ability.CharacterState()),
		Mechanics:      ability.Mechanics(),
	}
}

func ptr[T any](v T) *T { return &v }
