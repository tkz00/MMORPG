package repository

import (
	"tkz00/backend/pkg/game/entities"
)

func GetSkeletonEnemyAbilities(gs *entities.GameState) map[string]*entities.Ability {
	skeletonEnemiesAbilities := map[string]*entities.Ability{}
	return skeletonEnemiesAbilities
}

func GetPlayerAbilities(gs *entities.GameState) map[string]*entities.Ability {
	return map[string]*entities.Ability{}
}
