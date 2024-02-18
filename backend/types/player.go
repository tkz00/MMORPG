package types

const SPEED float32 = 1.0

type Player struct {
	position Position
	to       Position
}

func CreatePlayer(x, z float32) *Player {
	initPosition := Position{
		x: x,
		z: z,
	}

	return &Player{
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
}

func (p Player) IsMoving() bool {
	return p.position.x != p.to.x || p.position.z != p.to.z
}

func (p *Player) UpdatePosition() {
	direction := Position{}

	if p.position.x < p.to.x {
		direction.x += SPEED
	} else if p.position.x > p.to.x {
		direction.x -= SPEED
	}

	if p.position.z < p.to.z {
		direction.z += SPEED
	} else if p.position.z > p.to.z {
		direction.z -= SPEED
	}

	p.position.Move(direction)
}

func (p Player) ToDTO(id string) PlayerDTO {
	return PlayerDTO{
		Id:       id,
		Position: p.position.ToDTO(),
	}
}
