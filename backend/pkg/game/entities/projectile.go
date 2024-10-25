package entities

import (
	"tkz00/backend/pkg/utils"
)

const PROJECTILE_SPEED float64 = 15
const PROJECTILE_BOUNDS_RADIUS float64 = 0.25

type ProjectileState int

const (
	Active ProjectileState = iota
	Hit
)

type Projectile struct {
	id        string
	caster    string
	direction utils.Vector2
	position  utils.Vector2
	to        utils.Vector2
	damage    int
	state     ProjectileState
}

func CreateProjectile(id string, position utils.Vector2, targetDirection utils.Vector2, rangeValue float64, caster string) *Projectile {
	normalizedVector := utils.Normalize(position, targetDirection)
	to := normalizedVector.Scale(rangeValue).Add(position)
	direction := normalizedVector.Scale(PROJECTILE_SPEED)

	return &Projectile{
		id:        id,
		caster:    caster,
		direction: direction,
		position:  position,
		to:        to,
		damage:    40,
		state:     Active,
	}
}

func (p Projectile) GetId() string {
	return p.id
}

func (p Projectile) CasterId() string {
	return p.caster
}

func (p Projectile) GetDamage() int {
	return p.damage
}

func (p *Projectile) UpdatePosition(deltaTime float64) bool {
	distanceToTarget := p.position.Distance(p.to)
	if distanceToTarget < (PROJECTILE_SPEED * deltaTime) {
		p.position.Teleport(p.to)
		return true
	} else {
		p.position = p.position.Add(p.direction.Scale(deltaTime))
		return false
	}
}

func (p Projectile) GetPosition() utils.Vector2 {
	return p.position
}

func (p Projectile) GetRadius() float64 {
	return PROJECTILE_BOUNDS_RADIUS
}

func (p *Projectile) State() ProjectileState {
	return p.state
}

func (p *Projectile) SetState(state ProjectileState) {
	p.state = state
}

func (p Projectile) GetState() string {
	switch p.state {
	case Active:
		return "Active"
	case Hit:
		return "Hit"
	default:
		return "Unknown"
	}
}
