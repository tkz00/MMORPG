package character

import (
	"fmt"
	"time"
	"unnamed-mmo/backend/pkg/game/stats"
	"unnamed-mmo/backend/pkg/utils"
)

const BASE_MAX_HEALTH = 100
const CHARACTER_SPEED float64 = 10
const CHARACTER_BOUNDS_RADIUS float64 = 0.5

// need to rename this
type Action int

const (
	Idle Action = iota
	Moving
	Attacking
	CastingHeal
)

type Character struct {
	id string
	stats.Health
	ExecutingAction Action
	position        utils.Vector2
	to              utils.Vector2
	direction       utils.Vector2
	actionsQueue    []CharacterAction

	// should this be here?
	abilities map[string]*Ability
	lastUsed  map[string]time.Time

	onRemoved []func()
}

func CreateCharacter(id string, x, z float64, abilities map[string]*Ability) *Character {
	initialPosition := *utils.NewVector2(x, z)

	lastUsed := make(map[string]time.Time, len(abilities))
	for _, ability := range abilities {
		lastUsed[ability.id] = time.Time{}
	}

	return &Character{
		id:              id,
		position:        initialPosition,
		to:              initialPosition,
		Health:          stats.NewHealth(BASE_MAX_HEALTH),
		ExecutingAction: Idle,
		actionsQueue:    []CharacterAction{},
		abilities:       abilities,
		lastUsed:        lastUsed,
	}
}

func (p Character) GetId() string {
	return p.id
}

func (p *Character) SetPosition(position utils.Vector2) {
	p.position = position
}

func (p Character) GetPosition() utils.Vector2 {
	return p.position
}

// this shouldn't be here
func (p Character) GetHealth() stats.Health {
	return p.Health
}

func (p Character) GetExecutingAction() Action {
	return p.ExecutingAction
}

func (p Character) GetAbilities() map[string]*Ability {
	return p.abilities
}

func (p *Character) EnqueueAction(action CharacterAction) {
	p.actionsQueue = append(p.actionsQueue, action)
}

func (p *Character) ClearActionsQueue() {
	p.actionsQueue = nil
}

func (p *Character) PrependAction(action CharacterAction) {
	p.actionsQueue = append([]CharacterAction{action}, p.actionsQueue...)
}

func (c *Character) ExecuteNextAction() {
	if len(c.actionsQueue) == 0 {
		c.ExecutingAction = Idle
		return
	}

	// this logic should be different, bc now every action is executed on each update, this has no effect on the
	// behavior of the system (at least for now) but it could be confusing. Maybe a way to fix is to have the
	// queue and a single executingAction field of type character action, when the last executingAction is
	// completed (IsComplete()), the next action in the queue is moved to this single field and it's Execute
	// function called, probably can be done with just the queue and changing some of the logic
	currentAction := c.actionsQueue[0]
	err := currentAction.Execute(c)
	if err != nil {
		fmt.Println("Error executing action:", err)
		return
	}

	if currentAction.IsComplete() {
		c.actionsQueue = c.actionsQueue[1:]
	}
}

func (p *Character) MoveTowards(to utils.Vector2) {
	p.to = to
	p.direction = utils.Normalize(p.position, p.to).Scale(CHARACTER_SPEED)
	p.ExecutingAction = Moving
}

func (p Character) IsMoving() bool {
	return !p.position.Equals(p.to)
}

func (p *Character) UpdatePosition(deltaTime float64) {
	distanceToTarget := p.position.Distance(p.to)
	if distanceToTarget < (CHARACTER_SPEED * deltaTime) {
		p.position.Teleport(p.to)
	} else {
		p.position = p.position.Add(p.direction.Scale(deltaTime))
	}
}

func (p Character) GetRadius() float64 {
	return CHARACTER_BOUNDS_RADIUS
}

func (character *Character) EnqueueAbilityCast(castAbilityAction CastAbilityAction) {
	character.EnqueueAction(&castAbilityAction)
}

func (player *Character) CheckCooldown(abilityId string) bool {
	ability := player.abilities[abilityId]
	if player.RemainingCooldown(ability) <= 0 {
		return true
	} else {
		return false
	}
}

func (player Character) RemainingCooldown(ability *Ability) int64 {
	now := time.Now()
	remainingCooldown := player.abilities[ability.id].cooldown - now.Sub(player.lastUsed[ability.id]).Milliseconds()
	if remainingCooldown > 0 {
		return remainingCooldown
	}
	return 0
}

func (player *Character) ClosestPositionInRange(target utils.Vector2, rangeValue float64) utils.Vector2 {
	totalDistance := player.position.Distance(target)
	normalizedMovementVector := utils.Normalize(player.position, target)
	movementVector := normalizedMovementVector.Scale(totalDistance - rangeValue)
	targetPosition := player.position.Add(movementVector)
	return targetPosition
}

func (player *Character) CastAbility(ability Ability, params AbilityParameters) {
	player.lastUsed[ability.id] = time.Now()
	player.ExecutingAction = ability.characterState
	ability.Cast(*player, params)
}

// removal => disconnection, for now
func (c *Character) SubscribeToRemoval(callback func()) {
	c.onRemoved = append(c.onRemoved, callback)
}

func (c *Character) Remove() {
	for _, callback := range c.onRemoved {
		callback()
	}
}

// where should this be?
type AbilityInfo interface {
	GetId() string
	GetTargetPosition() (utils.Vector2, error)
	GetTargetId() (string, error)
}
