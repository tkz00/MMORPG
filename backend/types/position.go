package types

type Position struct {
	x float32
	z float32
}

func CreatePosition(data []byte) Position {
	positionDTO := CreatePositionDTO(data)

	return *GetMapper().PositionDTOToEntity(*positionDTO)
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

func (p Position) Equals(other Position) bool {
	return p.x == other.x && p.z == other.z
}