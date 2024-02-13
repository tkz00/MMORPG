package types

type Player struct {
	position Position
}

func CreatePlayer(x, z float32) Player {
	return Player{
		position: Position{X: x, Z: z},
	}
}

func (p Player) SetPosition(position Position) {
	p.position = position
}