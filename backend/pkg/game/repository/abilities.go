package repository

import (
	"tkz00/backend/pkg/game/entities"
	"tkz00/backend/pkg/utils"
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
			12,
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
						{
							MechanicType: "create_projectile",
							Params: map[string]interface{}{
								"targeting_strategy": "arc",
								"number":             5,
								"radius":             utils.DegreesToRadians(360),
								"range":              5.0,
								"origin_position":    "target",
								"on_hit_mechanics": []entities.Mechanic{
									{
										MechanicType: "damage",
										Params: map[string]interface{}{
											"amount":             20,
											"targeting_strategy": "character_hit",
										},
									},
								},
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
					"amount":             0,
					"targeting_strategy": "character_hit",
					"on_hit_mechanics": []entities.Mechanic{
						{
							MechanicType: "delay",
							Params: map[string]interface{}{
								"delay_ms": 2000,
								"execute_after_delay_mechanics": []entities.Mechanic{
									{
										MechanicType: "damage",
										Params: map[string]interface{}{
											"amount":             20,
											"targeting_strategy": "character_hit",
										},
									},
								},
							},
						},
					},
				},
			},
		),
	}
}
