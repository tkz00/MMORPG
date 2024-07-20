package types

import (
	"fmt"
	"math"
	"unnamed-mmo/backend/utils"

	"github.com/google/uuid"
)

const BASE_MAX_HEALTH = 100
const PLAYER_SPEED float64 = 10
const PLAYER_BOUNDS_RADIUS float64 = 0.5

type Action int

const (
	Idle Action = iota
	Attacking
	CastingHeal
)

type Player struct {
	id 		  		string
	stats	  		PlayerStats
	executingAction Action
	position  		Position
	to        		Position
	direction 		Position
}

func CreatePlayer(x, z float64, id string) *Player {
	initPosition := Position{
		x: x,
		z: z,
	}

	return &Player{
		id: id,
		position: initPosition,
		to:       initPosition,
		stats: PlayerStats{
			currentHealth: BASE_MAX_HEALTH,
			maxHealth:     BASE_MAX_HEALTH,
		},
		executingAction: Idle,
	}
}

func (p Player) GetId() string {
	return p.id
}

func (p *Player) SetPosition(position Position) {
	p.position = position
}

func (p Player) GetPosition() Position {
	return p.position
}

func (p Player) GetStats() PlayerStats {
	return p.stats
}

func (p Player) GetExecutingAction() Action {
	return p.executingAction
}

func (p *Player) MoveTowards(to Position) {
	p.to = to

	diffX, diffZ := utils.GetDiff(p.position.x, p.position.z, p.to.x, p.to.z)
	distanceMagnitude := math.Hypot(diffX, diffZ)

	p.direction = Position{
		x: diffX * PLAYER_SPEED / distanceMagnitude,
		z: diffZ * PLAYER_SPEED / distanceMagnitude,
	}
}

func (p Player) IsMoving() bool {
	return p.position.x != p.to.x || p.position.z != p.to.z
}

func (p *Player) UpdatePosition(deltaTime float64) {
	diffX, diffZ := utils.GetDiff(p.position.x, p.position.z, p.to.x, p.to.z)
	distanceToTarget := utils.GetDistance(diffX, diffZ)
	if distanceToTarget < (PLAYER_SPEED * deltaTime){
		p.position.Teleport(p.to)
	} else {
		p.position.Move(p.direction.Multiply(deltaTime))
	}
}

func (p *Player) DealDamage(damagePoints int) {
	newHealth := p.stats.currentHealth - damagePoints
	if newHealth > 0 {
		if p.stats.maxHealth > newHealth {
			p.stats.currentHealth -= damagePoints
		} else {
			p.stats.currentHealth = p.stats.maxHealth
		}
	} else {
		p.stats.currentHealth = 0
	}
}

func (p Player) GetRadius() float64 {
	return PLAYER_BOUNDS_RADIUS
}

func (player *Player) CastAbility(gameState *GameState, abilityInfo AbilityInfo) {
	switch abilityInfo.GetName() {
	case "projectile":
		// is this ID necessary?
		projectileId := uuid.New().String()

		targetPosition, err := abilityInfo.GetTargetPosition()
        if err != nil {
            fmt.Println("Error:", err)
            return
        }
	
		gameState.projectiles[projectileId] = CreateProjectile(projectileId, player.position, targetPosition, player.id)
		player.executingAction = Attacking
	case "heal":
		targetId, err := abilityInfo.GetTargetId()
		if err != nil {
            fmt.Println("Error:", err)
            return
        }

		gameState.players[targetId].DealDamage(-10)
		player.executingAction = CastingHeal
	default:
	}
}

type AbilityInfo interface {
	GetName() string
	GetTargetPosition() (Position, error)
	GetTargetId() (string, error)
}
