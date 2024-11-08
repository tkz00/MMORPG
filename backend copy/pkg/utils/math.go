package utils

import (
	"math"
)

func GetDiff(fromX, fromZ, toX, toZ float64) (float64, float64) {
	return toX - fromX, toZ - fromZ
}

func GetDistance(a, b float64) float64 {
	return math.Hypot(a, b)
}
