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

func (v Vector2) Subtract(subtractive Vector2) Vector2 {
	return Vector2{(v.x - subtractive.x), (v.z - subtractive.z)}
}

func (v *Vector2) Teleport(to Vector2) {
	v.x = to.x
	v.z = to.z
}

func (v Vector2) Equals(other Vector2) bool {
	return v.x == other.x && v.z == other.z
}

func (v Vector2) Dot(other Vector2) float64 {
	return v.x*other.x + v.z*other.z
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

// CalculateNewPosition calculates the new position after traveling a certain distance at a given angle.
func CalculateNewPosition(initial Vector2, distance float64, angleRadians float64) Vector2 {
	newX := initial.x + distance*math.Cos(angleRadians)
	newZ := initial.z + distance*math.Sin(angleRadians)
	return Vector2{x: newX, z: newZ}
}

func RandomCoordinatesInRadius(v Vector2, radius float64) Vector2 {
	angle := rand.Float64() * 2 * math.Pi
	r := math.Sqrt(rand.Float64()) * radius
	return CalculateNewPosition(v, r, angle)
}

func ClosestPositionInRange(startingPosition Vector2, target Vector2, rangeValue float64) Vector2 {
	totalDistance := startingPosition.Distance(target)
	normalizedMovementVector := Normalize(startingPosition, target)
	movementVector := normalizedMovementVector.Scale(totalDistance - rangeValue)
	targetPosition := startingPosition.Add(movementVector)
	return targetPosition
}

// CircleSegmentIntersect checks if a circle intersects or touches a line segment
func CircleSegmentIntersect(center Vector2, radius float64, segStart Vector2, segEnd Vector2) bool {
	// Segment vector (from segStart to segEnd)
	segment := segEnd.Subtract(segStart)

	// Vector from segment start to circle center
	toCenter := center.Subtract(segStart)

	// Project circle center onto the line segment
	segmentLengthSquared := segment.Dot(segment) // Length of the segment squared
	projection := toCenter.Dot(segment) / segmentLengthSquared

	// Clamp projection to the range [0, 1]
	t := math.Max(0, math.Min(1, projection))

	// Find the closest point on the segment to the circle's center
	closestPoint := segStart.Add(segment.Scale(t))

	// Calculate the distance from the closest point to the circle's center
	distanceToCenter := closestPoint.Distance(center)

	// Return true if the distance is less than or equal to the radius
	return distanceToCenter <= radius
}

// returns if circle and polygon are intersecting
func CirclePolygonIntersect(circleCenter Vector2, circleRadius float64, polygon []Vector2) bool {
	for index := range polygon {
		// Get the current and next vertices of the obstacle segment
		vertice := polygon[index]
		var nextVertice Vector2
		if index+1 < len(polygon) {
			nextVertice = polygon[index+1]
		} else {
			// If it's the last vertex, connect it to the first one (assuming closed shape)
			nextVertice = polygon[0]
		}

		// Check if the circle (character) intersects or touches the segment
		if CircleSegmentIntersect(circleCenter, circleRadius, vertice, nextVertice) {
			// Handle the collision, e.g., stop movement or adjust direction
			return true // Early exit if a collision is detected
		}
	}
	return false
}
