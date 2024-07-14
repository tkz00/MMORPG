package types

import (
	"math"
	"unnamed-mmo/backend/utils"
)

const RANGE float64 = 12
const PROJECTILE_SPEED float64 = 15
const PROJECTILE_BOUNDS_RADIUS float64 = 0.25

type ProjectileState int

const (
	Active ProjectileState = iota
	Hit
)

type Projectile struct {
	id			string
	caster		string
	direction	Position
	position 	Position
	to 			Position
	damage		int
	state 		ProjectileState
}

func CreateProjectile(id string, position Position, targetDirection Position, caster string) *Projectile {
	diffX, diffZ := utils.GetDiff(position.x, position.z, targetDirection.x, targetDirection.z)
	distanceMagnitude := math.Hypot(diffX, diffZ)
	xNormalized := diffX / distanceMagnitude
	zNormalized := diffZ / distanceMagnitude

	to := Position{
		x: xNormalized * RANGE + position.x,
		z: zNormalized * RANGE + position.z,
	}

	direction := Position{
		x: xNormalized * PROJECTILE_SPEED,
		z: zNormalized * PROJECTILE_SPEED,
	}

	return &Projectile{
		id: id,
		caster: caster,
		direction: direction,
		position: position,
		to: to,
		damage: 30,
		state: Active,
	}
}

func (p *Projectile) UpdatePosition(deltaTime float64) bool {
	diffX, diffZ := utils.GetDiff(p.position.x, p.position.z, p.to.x, p.to.z)
	distanceToTarget := utils.GetDistance(diffX, diffZ)
	if distanceToTarget < (PROJECTILE_SPEED * deltaTime) {
		p.position.Teleport(p.to)
		return true
	} else {
		p.position.Move(p.direction.Multiply(deltaTime))
		return false
	}
}

func (p Projectile) GetPosition() Position {
	return p.position
}

func (p Projectile) GetRadius() float64 {
	return PROJECTILE_BOUNDS_RADIUS
}

func (s ProjectileState) String() string {
    switch s {
    case Active:
        return "Active"
    case Hit:
        return "Hit"
    default:
        return "Unknown"
    }
}
