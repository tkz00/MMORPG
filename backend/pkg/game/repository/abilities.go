package repository

import (
	"fmt"
	"tkz00/backend/pkg/game/entities"
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

				target.HealthVariation(-2)
			}),
	}
	return skeletonEnemiesAbilities
}
