package types

import (
	"math"
	"unnamed-mmo/backend/utils"
)

const SPEED float64 = 1

type Player struct {
	id 		  string
	position  Position
	to        Position
	direccion Position
}

func CreatePlayer(x, z float32, id string) *Player {
	initPosition := Position{
		x: x,
		z: z,
	}

	return &Player{
		id: id,
		position: initPosition,
		to:       initPosition,
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
	distanceToTarget := math.Hypot(diffX, diffZ)
	if distanceToTarget < SPEED {
		p.position.Teleport(p.to)
	} else {
		p.position.Move(p.direccion)
	}
}