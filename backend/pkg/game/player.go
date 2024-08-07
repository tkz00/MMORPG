package game

import (
	"fmt"
	"time"
	"unnamed-mmo/backend/pkg/utils"
)

const BASE_MAX_HEALTH = 100
const PLAYER_SPEED float64 = 10
const PLAYER_BOUNDS_RADIUS float64 = 0.5

// need to rename this
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
	actionsQueue	[]CharacterAction

	// should this be here?
	abilities		map[string]*Ability
	lastUsed   		map[string]time.Time
}

func CreatePlayer(id string, x, z float64, abilities map[string]*Ability) *Player {
	initPosition := Position{
		x: x,
		z: z,
	}

	lastUsed := make(map[string]time.Time, len(abilities))
	for _, ability := range abilities {
        lastUsed[ability.id] = time.Time{}
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
		actionsQueue: []CharacterAction{},
		abilities: abilities,
		lastUsed: lastUsed,
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

func (p Player) GetAbilities() map[string]*Ability {
	return p.abilities
}

func (p *Player) EnqueueAction(action CharacterAction) {
    p.actionsQueue = append(p.actionsQueue, action)
}

func (p *Player) ExecuteNextAction(gameState *GameState) {
    if len(p.actionsQueue) == 0 {
        return
    }

    currentAction := p.actionsQueue[0]
    err := currentAction.Execute(p, gameState)
    if err != nil {
        fmt.Println("Error executing action:", err)
        return
    }

    if currentAction.IsComplete() {
        p.actionsQueue = p.actionsQueue[1:]
    }
}


func (p *Player) MoveTowards(to Position) {
	p.to = to

	diffX, diffZ := utils.GetDiff(p.position.x, p.position.z, p.to.x, p.to.z)
	distance := p.position.Distance(p.to)

	p.direction = Position{
		x: diffX * PLAYER_SPEED / distance,
		z: diffZ * PLAYER_SPEED / distance,
	}
}

func (p Player) IsMoving() bool {
	return p.position.x != p.to.x || p.position.z != p.to.z
}

func (p *Player) UpdatePosition(deltaTime float64) {
	distanceToTarget := p.position.Distance(p.to)
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
	if player.CheckCooldown(abilityInfo) {
		ability := player.abilities[abilityInfo.GetId()]

		// following behaviour should be inside ability, not player
		switch ability.name {
			case "projectile":
				targetPosition, err := abilityInfo.GetTargetPosition()
				if err != nil {
					fmt.Println("Error:", err)
					return
				}

				player.actionsQueue = player.actionsQueue[:0]
				projectileAction := &ProjectileAction{
					targetPosition: targetPosition,
					ability: *ability,
				}
				player.EnqueueAction(projectileAction)
			case "heal":
				targetId, err := abilityInfo.GetTargetId()
				if err != nil {
					fmt.Println("Error:", err)
					return
				}
				target := gameState.players[targetId]
				player.actionsQueue = player.actionsQueue[:0]
				if player.position.Distance(target.position) > ability.rangeValue {
					targetPosition := player.closestPositionInRange(target, ability)
					moveAction := &MoveAction{
						targetPosition: targetPosition,
					}
					player.EnqueueAction(moveAction)
				}
				healAction := &HealAction{
					targetId: targetId,
					ability: *ability,
				}
				player.EnqueueAction(healAction)
			default:
		}
	}
}

func (player *Player) closestPositionInRange(target *Player, ability *Ability) Position {
	diffX, diffZ := utils.GetDiff(player.position.x, player.position.z, target.position.x, target.position.z)
	totalDistance := player.position.Distance(target.position)
	normalizedMovementVector := Position{
		x: diffX / totalDistance,
		z: diffZ / totalDistance,
	}
	movementVector := Position{
		x: normalizedMovementVector.x * (totalDistance - ability.rangeValue),
		z: normalizedMovementVector.z * (totalDistance - ability.rangeValue),
	}
	targetPosition := Position{
		x: player.position.x + movementVector.x,
		z: player.position.z + movementVector.z,
	}
	return targetPosition
}

func (player *Player) CheckCooldown(abilityInfo AbilityInfo) bool {
	ability := player.abilities[abilityInfo.GetId()]
	if player.RemainingCooldown(ability) <= 0 {
		return true
	} else {
		return false
	}
}

func (player Player) RemainingCooldown(ability *Ability) int64 {
	now := time.Now()
	remainingCooldown := player.abilities[ability.id].cooldown - now.Sub(player.lastUsed[ability.id]).Milliseconds()
	if remainingCooldown > 0 {
		return remainingCooldown
	}
	return 0
}

type AbilityInfo interface {
	GetId() string
	GetTargetPosition() (Position, error)
	GetTargetId() (string, error)
}
