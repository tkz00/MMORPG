package types

type Player struct {
	position Position
}

func CreatePlayer(x, z float32) *Player {
	return &Player{
		position: Position{x: x, z: z},
	}
}

func (p *Player) SetPosition(position Position) {
	p.position = position
}

func (p Player) GetPosition() Position {
	return p.position
}

func (p Player) ToDTO(id string) PlayerDTO {
	return PlayerDTO{
		Id:       id,
		Position: p.position.ToDTO(),
	}
}
