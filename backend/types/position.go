package types

type Position struct {
	x float64
	z float64
}

func CreatePosition(data []byte) Position {
	positionDTO := CreatePositionDTO(data)

	return *GetMapper().PositionDTOToEntity(*positionDTO)
}

func (p Position) GetPosition() (float64, float64) {
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

func (p Position) Equals(other Position) bool {
	return p.x == other.x && p.z == other.z
}
