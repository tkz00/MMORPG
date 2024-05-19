package types

import (
	"math"
	"unnamed-mmo/backend/utils"
)

const RANGE float64 = 6
const ABILITY_SPEED float64 = 1

type Projectile struct {
	caster		string
	direction	Position
	position 	Position
	to 			Position
	damage		int
}

func CreateProjectile(position Position, input Position, caster string) *Projectile {
	diffX, diffZ := utils.GetDiff(position.x, position.z, input.x, input.z)
	distanceMagnitude := math.Hypot(diffX, diffZ)
	xNormalized := diffX / distanceMagnitude
	zNormalized := diffZ / distanceMagnitude

	to := Position{
		x: float32(xNormalized * RANGE),
		z: float32(zNormalized * RANGE),
	}

	direction := Position {
		x: float32(xNormalized * ABILITY_SPEED),
		z: float32(zNormalized * ABILITY_SPEED),
	}

	return &Projectile{
		caster: caster,
		direction: direction,
		position: position,
		to: to,
		damage: 30,
	}
}

func (p *Projectile) UpdatePosition() bool {
	diffX, diffZ := utils.GetDiff(p.position.x, p.position.z, p.to.x, p.to.z)
	distanceToTarget := utils.GetDistance(diffX, diffZ)
	if distanceToTarget < ABILITY_SPEED {
		p.position.Teleport(p.to)
	} else {
		p.position.Move(p.direction)
	}

	return p.position.x == p.to.x && p.position.z == p.to.z
}

func (p Projectile) GetPosition() Position {
	return p.position
}
