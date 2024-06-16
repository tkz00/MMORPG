package types

import (
	"math"
	"unnamed-mmo/backend/utils"
)

const RANGE float64 = 12
const ABILITY_SPEED float64 = 1

type Projectile struct {
	id			string
	caster		string
	direction	Position
	position 	Position
	to 			Position
	damage		int
}

func CreateProjectile(id string, position Position, targetDirection Position, caster string) *Projectile {
	diffX, diffZ := utils.GetDiff(position.x, position.z, targetDirection.x, targetDirection.z)
	distanceMagnitude := math.Hypot(diffX, diffZ)
	xNormalized := diffX / distanceMagnitude
	zNormalized := diffZ / distanceMagnitude

	to := Position{
		x: float32(xNormalized * RANGE) + position.x,
		z: float32(zNormalized * RANGE) + position.z,
	}

	direction := Position{
		x: float32(xNormalized * ABILITY_SPEED),
		z: float32(zNormalized * ABILITY_SPEED),
	}

	return &Projectile{
		id: id,
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
		return true
	} else {
		p.position.Move(p.direction)
		return false
	}
}

func (p Projectile) GetPosition() Position {
	return p.position
}
