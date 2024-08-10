package utils

import (
	"math"
	"math/rand"
)

type Vector2 struct {
	x float64
	z float64
}

func NewVector2(x float64, z float64) *Vector2 {
	return &Vector2{x: x, z: z}
}

func (v Vector2) GetPosition() (float64, float64) {
	return v.x, v.z
}

func (v Vector2) Add(additive Vector2) Vector2 {
	return Vector2{(v.x + additive.x), (v.z + additive.z)}
}

func (v *Vector2) Teleport(to Vector2) {
	v.x = to.x
	v.z = to.z
}

func (v Vector2) Equals(other Vector2) bool {
	return v.x == other.x && v.z == other.z
}

func (v Vector2) Scale(scalar float64) Vector2 {
	return Vector2{
		x: v.x * scalar,
		z: v.z * scalar,
	}
}

func (v1 Vector2) Distance(v2 Vector2) float64 {
	diffX, diffZ := GetDiff(v1.x, v1.z, v2.x, v2.z)
	return GetDistance(diffX, diffZ)
}

func Normalize(a, b Vector2) Vector2 {
	diffX, diffZ := GetDiff(a.x, a.z, b.x, b.z)
	distanceMagnitude := math.Hypot(diffX, diffZ)
	return Vector2{diffX / distanceMagnitude, diffZ / distanceMagnitude}
}

func RandomCoordinatesInRadius(v Vector2, radius float64) Vector2 {
	angle := rand.Float64() * 2 * math.Pi
	r := math.Sqrt(rand.Float64()) * radius
	x := v.x + r * math.Cos(angle)
	z := v.z + r * math.Sin(angle)
	return Vector2{x: x, z: z}
}
