package types

type Player struct {
	position Position
}

func (p Player) CreatePlayer(x, z float32) Player {
	return Player{
		position: Position{X: x, Z: z},
	}
}
