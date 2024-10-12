package repository

import (
	"tkz00/backend/pkg/utils"
)

func GetObstacleColliders() [][]utils.Vector2 {
	tomb := []utils.Vector2{
		*utils.NewVector2(-12.1, 2.25),
		*utils.NewVector2(-7.9, 2.25),
		*utils.NewVector2(-7.9, -2.2),
		*utils.NewVector2(-12.1, -2.2),
	}

	return [][]utils.Vector2{tomb}
}
