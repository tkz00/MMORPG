package types

import (
	"math"
	"math/rand"
	"unnamed-mmo/backend/utils"
)

const BASE_MAX_HEALTH = 100
const SPEED float64 = 1
const BOUNDS_RADIUS float64 = 0.5

type Collisionable interface {
	GetBounds() Position
}

type Player struct {
	stats     PlayerStats
	position  Position
	to        Position
	direccion Position
}

func CreatePlayer(x, z float32) *Player {
	initPosition := Position{
		x: x,
		z: z,
	}

	return &Player{
		position: initPosition,
		to:       initPosition,
		stats: PlayerStats{
			currentHealth: rand.Intn(BASE_MAX_HEALTH) + 1,
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

	p.direccion = Position{
		x: float32(diffX * SPEED / distanceMagnitude),
		z: float32(diffZ * SPEED / distanceMagnitude),
	}
}

func (p Player) IsMoving() bool {
	return p.position.x != p.to.x || p.position.z != p.to.z
}

func (p *Player) UpdatePosition() {
	diffX, diffZ := utils.GetDiff(p.position.x, p.position.z, p.to.x, p.to.z)
	distanceToTarget := utils.GetDistance(diffX, diffZ)
	if distanceToTarget < SPEED {
		p.position.Teleport(p.to)
	} else {
		p.position.Move(p.direccion)
	}
}

func (p *Player) RandomizeHealth() {
	p.stats.currentHealth = rand.Intn(BASE_MAX_HEALTH) + 1
}

func (p Player) ToDTO(id string) PlayerDTO {
	return PlayerDTO{
		Id:       id,
		Position: p.position.ToDTO(),

		MaxHealth:     p.stats.maxHealth,
		CurrentHealth: p.stats.currentHealth,
	}
}

func (p Player) GetRadius() float64 {
	return BOUNDS_RADIUS
}
