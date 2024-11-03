package repository

import (
	"fmt"
	"tkz00/backend/pkg/game/entities"

	"github.com/google/uuid"
)

func GetSkeletonEnemyAbilities(gs entities.GameState) map[string]*entities.Ability {
	skeletonEnemiesAbilities := map[string]*entities.Ability{
		"0": entities.NewAbility("0", "sword slash", 2, 1500, entities.Target, entities.Attacking,
			func(caster entities.Character, params entities.AbilityParameters) {
				targetId := params.(entities.TargetIdAbilityParams).TargetId

				target, err := gs.GetCharacterById(targetId)
				if err != nil {
					fmt.Println(err)
					return
				}

				target.HealthVariation(-10)
			}),
	}
	return skeletonEnemiesAbilities
}

func GetPlayerAbilities(gs *entities.GameState) map[string]*entities.Ability {
	return map[string]*entities.Ability{
		"0": entities.NewAbility("0", "projectile", 5, 2000, entities.Coordinates, entities.Attacking,
			func(caster entities.Character, params entities.AbilityParameters) {
				projectile := entities.CreateProjectile(uuid.NewString(), caster.GetPosition(), params.(entities.CoordinateAbilityParams).Target, 5, caster.GetId())
				gs.AddProjectile(projectile)
			}),
		"1": entities.NewAbility("1", "heal", 7, 3000, entities.Target, entities.CastingHeal,
			func(caster entities.Character, params entities.AbilityParameters) {
				targetId := params.(entities.TargetIdAbilityParams).TargetId

				target, err := gs.GetCharacterById(targetId)
				if err != nil {
					fmt.Println(err)
					return
				}

				target.HealthVariation(10)
			}),
	}
}
