package repository

import (
	"tkz00/backend/pkg/game/entities"
)

func GetSkeletonEnemyAbilities(gs *entities.GameState) map[string]*entities.Ability {
	skeletonEnemiesAbilities := map[string]*entities.Ability{
		"0": entities.NewAbility(
			"0",
			"sword slash",
			2,
			4000,
			entities.Target,
			entities.Attacking,
			entities.Mechanic{MechanicType: "damage", Params: map[string]interface{}{"amount": 19}},
		),
	}
	return skeletonEnemiesAbilities
}

func GetPlayerAbilities(gs *entities.GameState) map[string]*entities.Ability {
	return map[string]*entities.Ability{
		"0": entities.NewAbility(
			"0",
			"projectile",
			5,
			2000,
			entities.Coordinates,
			entities.Attacking,
			entities.Mechanic{MechanicType: "create_projectile", Params: map[string]interface{}{}},
		),
		"1": entities.NewAbility(
			"1",
			"heal",
			7,
			3000,
			entities.Target,
			entities.CastingHeal,
			entities.Mechanic{MechanicType: "heal", Params: map[string]interface{}{"amount": 40}},
		),
	}
}
