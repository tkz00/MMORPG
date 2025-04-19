package configurator

import (
	"backend/pkg/game/entities"
)

func GetSeedsAbilities() map[string]*entities.Ability {
	return map[string]*entities.Ability{
		"0": entities.NewAbility(
			"0",
			"Sword Slash",
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
		// "0": entities.NewAbility(
		// 	"0",
		// 	"projectile",
		// 	12,
		// 	2000,
		// 	entities.Coordinates,
		// 	entities.Attacking,
		// 	entities.Mechanic{
		// 		MechanicType: "create_projectile",
		// 		Params: map[string]interface{}{
		// 			"on_hit_mechanics": []entities.Mechanic{
		// 				{
		// 					MechanicType: "damage",
		// 					Params: map[string]interface{}{
		// 						"amount":             40,
		// 						"targeting_strategy": "character_hit",
		// 					},
		// 				},
		// 				{
		// 					MechanicType: "heal",
		// 					Params: map[string]interface{}{
		// 						"amount":             20,
		// 						"targeting_strategy": "caster",
		// 					},
		// 				},
		// 				{
		// 					MechanicType: "create_projectile",
		// 					Params: map[string]interface{}{
		// 						"targeting_strategy": "arc",
		// 						"number":             5,
		// 						"radius":             utils.DegreesToRadians(360),
		// 						"range":              5.0,
		// 						"origin_position":    "target",
		// 						"on_hit_mechanics": []entities.Mechanic{
		// 							{
		// 								MechanicType: "damage",
		// 								Params: map[string]interface{}{
		// 									"amount":             20,
		// 									"targeting_strategy": "character_hit",
		// 								},
		// 							},
		// 						},
		// 					},
		// 				},
		// 			},
		// 		},
		// 	},
		// ),
		"1": entities.NewAbility(
			"1",
			"Projectile",
			12,
			2000,
			entities.Coordinates,
			entities.Attacking,
			entities.Mechanic{
				MechanicType: "create_projectile",
				Params: map[string]interface{}{
					"on_hit_mechanics": []entities.Mechanic{
						{
							MechanicType: "create_AoE",
							Params: map[string]interface{}{
								"radius":      3.0,
								"duration_ms": 400,
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
		"2": entities.NewAbility(
			"2",
			"Heal",
			7,
			3000,
			entities.Target,
			entities.CastingHeal,
			entities.Mechanic{
				MechanicType: "damage",
				Params: map[string]interface{}{
					"amount":             -20,
					"targeting_strategy": "character_hit",
				},
			},
		),
		"3": entities.NewAbility(
			"3",
			"Life Drain",
			7,
			3000,
			entities.Target,
			entities.CastingHeal,
			entities.Mechanic{
				MechanicType: "damage",
				Params: map[string]interface{}{
					"amount":             20,
					"targeting_strategy": "character_hit",
					"on_hit_mechanics": []entities.Mechanic{
						{
							MechanicType: "delay",
							Params: map[string]interface{}{
								"delay_ms": 1000,
								"execute_after_delay_mechanics": []entities.Mechanic{
									{
										MechanicType: "damage",
										Params: map[string]interface{}{
											"amount":             -10,
											"targeting_strategy": "caster",
										},
									},
								},
							},
						},
					},
				},
			},
		),
		"4": entities.NewAbility(
			"4",
			"Ground Slam",
			2,
			2000,
			entities.Coordinates,
			entities.Attacking,
			entities.Mechanic{
				MechanicType: "create_AoE",
				Params: map[string]interface{}{
					"radius":      1.0,
					"duration_ms": 400,
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
		),
	}
}
