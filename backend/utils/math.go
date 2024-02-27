package utils

func GetDiff(fromX, fromZ, toX, toZ float32) (float64, float64) {
	return float64(toX - fromX), float64(toZ - fromZ)
}
