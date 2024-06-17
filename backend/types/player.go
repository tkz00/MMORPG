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

type Player struct {
	id 		  string
	stats	  PlayerStats
	position  Position
	to        Position
	direction Position
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
	}
}


func (p *Player) SetPosition(position Position) {
	p.position = position
}

func (p Player) GetPosition() Position {
	return p.position
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
	p.stats.currentHealth -= damagePoints
}

func (p Player) GetRadius() float64 {
	return PLAYER_BOUNDS_RADIUS
}

func (player *Player) CastAbility(gameState *GameState, abilityInfo AbilityCastDTO) {
	switch abilityInfo.Name {
	case "projectile":
		// is this ID necessary?
		projectileId := uuid.New().String()

		targetPosition, err := extractTargetPosition(abilityInfo.AbilityParameters)
        if err != nil {
            fmt.Println("Error:", err)
            return
        }
	
		gameState.projectiles[projectileId] = CreateProjectile(projectileId, player.position, targetPosition, player.id)
	case "heal":
		targetId, err := extractTargetId(abilityInfo.AbilityParameters)
		if err != nil {
            fmt.Println("Error:", err)
            return
        }
        fmt.Println(targetId)

		gameState.players[targetId].DealDamage(-10)
	default:
	}
}
