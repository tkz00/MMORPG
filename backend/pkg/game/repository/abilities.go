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
			entities.Mechanic{
				MechanicType: "damage",
				Params: map[string]interface{}{
					"amount":             19,
					"targeting_strategy": "character_hit",
				},
			},
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
			entities.Mechanic{
				MechanicType: "create_projectile",
				Params: map[string]interface{}{
					"on_hit_mechanics": []entities.Mechanic{
						{
							MechanicType: "damage",
							Params: map[string]interface{}{
								"amount":             40,
								"targeting_strategy": "character_hit",
							},
						},
						{
							MechanicType: "heal",
							Params: map[string]interface{}{
								"amount":             20,
								"targeting_strategy": "caster",
							},
						},
					},
				},
			},
		),
		"1": entities.NewAbility(
			"1",
			"heal",
			7,
			3000,
			entities.Target,
			entities.CastingHeal,
			entities.Mechanic{
				MechanicType: "heal",
				Params: map[string]interface{}{
					"amount":             40,
					"targeting_strategy": "character_hit",
				},
			},
		),
	}
}
