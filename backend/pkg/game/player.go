package game

import (
	"fmt"
	"time"
	"unnamed-mmo/backend/pkg/game/stats"
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
	health	  		stats.Health
	executingAction Action
	position  		utils.Vector2
	to        		utils.Vector2
	direction 		utils.Vector2
	actionsQueue	[]CharacterAction

	// should this be here?
	abilities		map[string]*Ability
	lastUsed   		map[string]time.Time
}

func CreatePlayer(id string, x, z float64, abilities map[string]*Ability) *Player {
	initialPosition := *utils.NewVector2(x, z)

	lastUsed := make(map[string]time.Time, len(abilities))
	for _, ability := range abilities {
        lastUsed[ability.id] = time.Time{}
    }

	return &Player{
		id: id,
		position: initialPosition,
		to:       initialPosition,
		health: stats.NewHealth(BASE_MAX_HEALTH),
		executingAction: Idle,
		actionsQueue: []CharacterAction{},
		abilities: abilities,
		lastUsed: lastUsed,
	}
}

func (p Player) GetId() string {
	return p.id
}

func (p *Player) SetPosition(position utils.Vector2) {
	p.position = position
}

func (p Player) GetPosition() utils.Vector2 {
	return p.position
}

func (p Player) GetHealth() stats.Health {
	return p.health
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

func (p *Player) ClearActionsQueue() {
    p.actionsQueue = nil
}

func (p *Player) ExecuteNextAction(gameState *GameState) {
    if len(p.actionsQueue) == 0 {
        return
    }

    currentAction := p.actionsQueue[0]
    err := currentAction.Execute(p)
    if err != nil {
        fmt.Println("Error executing action:", err)
        return
    }

    if currentAction.IsComplete() {
        p.actionsQueue = p.actionsQueue[1:]
    }
}


func (p *Player) MoveTowards(to utils.Vector2) {
	p.to = to
	p.direction = utils.Normalize(p.position, p.to).Scale(PLAYER_SPEED)
}

func (p Player) IsMoving() bool {
	return !p.position.Equals(p.to)
}

func (p *Player) UpdatePosition(deltaTime float64) {
	distanceToTarget := p.position.Distance(p.to)
	if distanceToTarget < (PLAYER_SPEED * deltaTime){
		p.position.Teleport(p.to)
	} else {
		p.position = p.position.Add(p.direction.Scale(deltaTime))
	}
}

func (p Player) GetRadius() float64 {
	return PLAYER_BOUNDS_RADIUS
}

func (player *Player) EnqueueAbilityCast(gameState *GameState, abilityInfo AbilityInfo) {
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

				player.ClearActionsQueue()
				coordinateAbilityParams := CoordinateAbilityParams{
					target: targetPosition,
				}
				castProjectileAction := &CastAbilityAction{
					ability: *ability,
					params: coordinateAbilityParams,
				}
				player.EnqueueAction(castProjectileAction)
			case "heal":
				targetId, err := abilityInfo.GetTargetId()
				if err != nil {
					fmt.Println("Error:", err)
					return
				}
				player.ClearActionsQueue()
				target := gameState.players[targetId]
				player.actionsQueue = player.actionsQueue[:0]
				if player.position.Distance(target.position) > ability.rangeValue {
					targetPosition := player.closestPositionInRange(target, ability)
					moveAction := &MoveAction{
						targetPosition: targetPosition,
					}
					player.EnqueueAction(moveAction)
				}
				targetAbilityParams := TargetIdAbilityParams{
					targetId: targetId,
				}
				castHealAction := &CastAbilityAction{
					ability: *ability,
					params: targetAbilityParams,
				}
				player.EnqueueAction(castHealAction)
			default:
		}
	}
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

func (player *Player) closestPositionInRange(target *Player, ability *Ability) utils.Vector2 {
	totalDistance := player.position.Distance(target.position)
	normalizedMovementVector := utils.Normalize(player.position, target.position)
	movementVector := normalizedMovementVector.Scale(totalDistance - ability.rangeValue)
	targetPosition := player.position.Add(movementVector)
	return targetPosition
}

func (player *Player) CastAbility(ability Ability, params AbilityParameters) {
	player.lastUsed[ability.id] = time.Now()
	ability.Cast(*player, params)
}

type AbilityInfo interface {
	GetId() string
	GetTargetPosition() (utils.Vector2, error)
	GetTargetId() (string, error)
}
