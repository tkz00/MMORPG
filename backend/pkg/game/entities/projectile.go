package entities

import (
	"backend/pkg/utils"

	"github.com/google/uuid"
)

const PROJECTILE_SPEED float64 = 15
const PROJECTILE_BOUNDS_RADIUS float64 = 0.25

type ProjectileState int

const (
	Active ProjectileState = iota
	Hit
)

type Projectile struct {
	id             string
	caster         string
	direction      utils.Vector2
	position       utils.Vector2
	to             utils.Vector2
	state          ProjectileState
	onHitMechanics []Mechanic

	ProjectileLastSentState *ProjectileLastTickState
}

func CreateProjectile(
	initialPosition utils.Vector2,
	targetDirection utils.Vector2,
	rangeValue float64,
	caster string,
	onHitMechanics []Mechanic,
) *Projectile {
	normalizedVector := utils.Normalize(initialPosition, targetDirection)
	to := normalizedVector.Scale(rangeValue).Add(initialPosition)
	direction := normalizedVector.Scale(PROJECTILE_SPEED)

	return &Projectile{
		id:                      uuid.New().String(),
		caster:                  caster,
		direction:               direction,
		position:                initialPosition,
		to:                      to,
		state:                   Active,
		onHitMechanics:          onHitMechanics,
		ProjectileLastSentState: &ProjectileLastTickState{Position: nil, State: nil},
	}
}

func (p Projectile) GetId() string {
	return p.id
}

func (p Projectile) CasterId() string {
	return p.caster
}

func (p *Projectile) UpdatePosition(deltaTime float64) bool {
	distanceToTarget := p.position.Distance(p.to)
	if distanceToTarget < (PROJECTILE_SPEED * deltaTime) {
		p.position.Teleport(p.to)
		return true
	} else {
		p.position = p.position.Add(p.direction.Scale(deltaTime))
		return false
	}
}

func (p Projectile) GetPosition() utils.Vector2 {
	return p.position
}

func (p Projectile) GetRadius() float64 {
	return PROJECTILE_BOUNDS_RADIUS
}

func (p *Projectile) State() ProjectileState {
	return p.state
}

func (p *Projectile) SetState(state ProjectileState) {
	p.state = state
}

func (p Projectile) GetState() string {
	switch p.state {
	case Active:
		return "Active"
	case Hit:
		return "Hit"
	default:
		return "Unknown"
	}
}

func (projectile Projectile) Hit(target *Character, gs *GameState) {
	for _, mechanic := range projectile.onHitMechanics {
		mechanic.Params["target_id"] = target.id

		switch mechanic.MechanicType {
		case "create_projectile":
			mechanic.Params["projectile_last_position"] = projectile.position
		case "create_AoE":
			mechanic.Params["target_coordinates"] = projectile.position
		}
	}
	resolveMechanics(projectile.caster, gs, projectile.onHitMechanics)
}

type ProjectileLastTickState struct {
	Position *utils.Vector2
	State    *ProjectileState
	Radius   *float64
}
