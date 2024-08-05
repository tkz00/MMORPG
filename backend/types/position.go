package types

import "unnamed-mmo/backend/utils"

type Position struct {
	x float64
	z float64
}

func NewPosition(x float64, z float64) *Position {
	return &Position{x: x, z: z}
}

func (p Position) GetPosition() (float64, float64) {
	return p.x, p.z
}

func (p *Position) Move(to Position) {
	p.x += to.x
	p.z += to.z
}

func (p *Position) Teleport(to Position) {
	p.x = to.x
	p.z = to.z
}

func (p Position) Equals(other Position) bool {
	return p.x == other.x && p.z == other.z
}

func (p Position) Multiply(multiplier float64) Position {
	return Position{p.x * multiplier, p.z * multiplier}
}

func (positionA Position) Distance(positionB Position) float64 {
	diffX, diffZ := utils.GetDiff(positionA.x, positionA.z, positionB.x, positionB.z)
	return utils.GetDistance(diffX, diffZ)
}
