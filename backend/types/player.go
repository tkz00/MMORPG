package types
type Player struct{
	x float32
	y float32
}

func (p Player) CreatePlayer(x, y float32) Player{
	return Player{
		x: x,
		y: y,
	}
}