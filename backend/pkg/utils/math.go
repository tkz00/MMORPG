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

func DegreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func RadiansToDegrees(radians float64) float64 {
	return radians * 180 / math.Pi
}
