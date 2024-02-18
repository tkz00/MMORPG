package types

type Position struct {
	x float32
	z float32
}

func (p Position) ToDTO() PositionDTO {
	return PositionDTO{
		X: p.x,
		Z: p.z,
	}
}