package utils

import "math"

func GetDiff(fromX, fromZ, toX, toZ float32) (float64, float64) {
	return float64(toX - fromX), float64(toZ - fromZ)
}

func GetDistance(a, b float64) float64 {
	return math.Hypot(a, b)
}
