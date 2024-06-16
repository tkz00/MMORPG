package types

import (
	"fmt"
	"math"
	"math/rand"
	"unnamed-mmo/backend/utils"

	"github.com/google/uuid"
)

const BASE_MAX_HEALTH = 100
const SPEED float64 = 1
const BOUNDS_RADIUS float64 = 0.5

type Player struct {
	id 		  string
	stats	  PlayerStats
	position  Position
	to        Position
	direccion Position
}

func CreatePlayer(x, z float32, id string) *Player {
	initPosition := Position{
		x: x,
		z: z,
	}

	return &Player{
		id: id,
		position: initPosition,
		to:       initPosition,
		stats: PlayerStats{
			currentHealth: rand.Intn(BASE_MAX_HEALTH) + 1,
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

	p.direccion = Position{
		x: float32(diffX * SPEED / distanceMagnitude),
		z: float32(diffZ * SPEED / distanceMagnitude),
	}
}

func (p Player) IsMoving() bool {
	return p.position.x != p.to.x || p.position.z != p.to.z
}

func (p *Player) UpdatePosition() {
	diffX, diffZ := utils.GetDiff(p.position.x, p.position.z, p.to.x, p.to.z)
	distanceToTarget := utils.GetDistance(diffX, diffZ)
	if distanceToTarget < SPEED {
		p.position.Teleport(p.to)
	} else {
		p.position.Move(p.direccion)
	}
}

func (p *Player) DealDamage(damagePoints int) {
	p.stats.currentHealth -= damagePoints
}

func (p Player) GetRadius() float64 {
	return BOUNDS_RADIUS
}

func (player *Player) CastAbility(gameState *GameState, abilityInfo AbilityCastDTO) {
	switch abilityInfo.Name {
	case "projectile":
		id := uuid.New().String()
		// Extract the target position map
		targetPositionMap, ok := abilityInfo.AbilityParameters[TargetPosition].(map[string]interface{})
		if !ok {
			fmt.Println("Error: Unable to cast TargetPosition to map[string]interface{}")
			return
		}

		// Extract the x and z values, converting them to float64
		xValue, xOk := targetPositionMap["x"]
		zValue, zOk := targetPositionMap["z"]

		if !xOk || !zOk {
			fmt.Println("Error: x or z value not found in the target position map")
			return
		}

		// Convert xValue and zValue to float64
		x, xConvOk := xValue.(float64)
		z, zConvOk := zValue.(float64)

		if !xConvOk || !zConvOk {
			fmt.Println("Error: x or z value could not be converted to float64")
			return
		}

		// Create a new Position instance
		targetPosition := Position{
			x: float32(x),
			z: float32(z),
		}
	
		gameState.projectiles[id] = CreateProjectile(id, player.position, targetPosition, player.id)
	case "heal":
	default:
	}
}
