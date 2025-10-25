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
			550,
			entities.Mechanic{
				MechanicType: "damage",
				Params: map[string]interface{}{
					"base_amount":            25,
					"damage_stat_multiplier": 0,
					"targeting_strategy":     "character_hit",
				},
			},
		),
		"1": entities.NewAbility(
			"1",
			"Projectile",
			12,
			2000,
			entities.Coordinates,
			entities.Attacking,
			700,
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
											"base_amount":            20,
											"damage_stat_multiplier": 1.5,
											"targeting_strategy":     "character_hit",
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
			1000,
			entities.Mechanic{
				MechanicType: "heal",
				Params: map[string]interface{}{
					"base_amount":            20,
					"damage_stat_multiplier": 0,
					"targeting_strategy":     "character_hit",
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
			1200,
			entities.Mechanic{
				MechanicType: "damage",
				Params: map[string]interface{}{
					"base_amount":        20,
					"targeting_strategy": "character_hit",
					"on_hit_mechanics": []entities.Mechanic{
						{
							MechanicType: "delay",
							Params: map[string]interface{}{
								"delay_ms": 1000,
								"execute_after_delay_mechanics": []entities.Mechanic{
									{
										MechanicType: "heal",
										Params: map[string]interface{}{
											"base_amount":            10,
											"damage_stat_multiplier": 0,
											"targeting_strategy":     "caster",
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
			800,
			entities.Mechanic{
				MechanicType: "create_AoE",
				Params: map[string]interface{}{
					"radius":      1.0,
					"duration_ms": 400,
					"on_hit_mechanics": []entities.Mechanic{
						{
							MechanicType: "damage",
							Params: map[string]interface{}{
								"base_amount":            20,
								"damage_stat_multiplier": 0,
								"targeting_strategy":     "character_hit",
							},
						},
					},
				},
			},
		),
		"5": entities.NewAbility(
			"5",
			"Buff Damage",
			7,
			3000,
			entities.Target,
			entities.CastingHeal,
			600,
			entities.Mechanic{
				MechanicType: "buff_stat",
				Params: map[string]interface{}{
					"target_stat":        "damage",
					"base_amount":        0,
					"multiplier":         0.5, // this is original value + (this * original value)
					"targeting_strategy": "character_hit",
					"duration_ms":        5000,
				},
			},
		),
	}
}
