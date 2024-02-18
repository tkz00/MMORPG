package types

type Player struct {
	Position Position
}

func CreatePlayer(x, z float32) *Player {
	return &Player{
		Position: Position{X: x, Z: z},
	}
}

func (p *Player) SetPosition(position Position) {
	p.Position = position
}
