package types

type Position struct {
	x float32
	z float32
}

func (p Position) GetPosition() (float32, float32) {
	return p.x, p.z
}

func (p *Position) Move(to Position) {
	p.x += to.x
	p.z += to.z
}

func (p *Position) Teleport(to Position) {
	p.x = to.x
	p.z = to.z
}

func (p Position) ToDTO() PositionDTO {
	return PositionDTO{
		X: p.x,
		Z: p.z,
	}
}
