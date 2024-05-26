package types

import (
	"math"
	"unnamed-mmo/backend/utils"
)

const BASE_MAX_HEALTH = 100
const PLAYER_SPEED float64 = 1
const PLAYER_BOUNDS_RADIUS float64 = 0.5

type Collisionable interface {
	GetBounds() Position
}

type Player struct {
	id 		  string
	stats	  PlayerStats
	position  Position
	to        Position
	direction Position
}

func CreatePlayer(x, z float64, id string) *Player {
	initPosition := Position{
		x: x,
		z: z,
	}

	return &Player{
		id: id,
		position: initPosition,
		to:       initPosition,
		stats: PlayerStats{
			currentHealth: BASE_MAX_HEALTH,
			maxHealth:     BASE_MAX_HEALTH,
		},
	}
}


func (p *Player) SetPosition(position Position) {
	p.position = position
}

func (p Player) GetPosition() Position {
	return p.position
}

func (p *Player) MoveTowards(to Position) {
	p.to = to

	diffX, diffZ := utils.GetDiff(p.position.x, p.position.z, p.to.x, p.to.z)
	distanceMagnitude := math.Hypot(diffX, diffZ)

	p.direction = Position{
		x: diffX * PLAYER_SPEED / distanceMagnitude,
		z: diffZ * PLAYER_SPEED / distanceMagnitude,
	}
}

func (p Player) IsMoving() bool {
	return p.position.x != p.to.x || p.position.z != p.to.z
}

func (p *Player) UpdatePosition() {
	diffX, diffZ := utils.GetDiff(p.position.x, p.position.z, p.to.x, p.to.z)
	distanceToTarget := utils.GetDistance(diffX, diffZ)
	if distanceToTarget < PLAYER_SPEED {
		p.position.Teleport(p.to)
	} else {
		p.position.Move(p.direction)
	}
}

func (p *Player) DealDamage(damagePoints int) {
	p.stats.currentHealth -= damagePoints
}

func (p Player) GetRadius() float64 {
	return PLAYER_BOUNDS_RADIUS
}
