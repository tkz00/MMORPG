package entities

import (
	"fmt"
	"time"
	"tkz00/backend/pkg/game/stats"
	"tkz00/backend/pkg/utils"
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

type ExecutingAction struct {
	actionType Action
	direction  utils.Vector2
}

func (action ExecutingAction) ActionType() Action {
	return action.actionType
}

func (action ExecutingAction) Direction() utils.Vector2 {
	return action.direction
}

type Character struct {
	id string
	stats.Health
	executingAction ExecutingAction
	position        utils.Vector2
	to              utils.Vector2
	direction       utils.Vector2
	actionsQueue    []CharacterAction

	*Inventory

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
		executingAction: ExecutingAction{Idle, *utils.NewVector2(0, 0)},
		actionsQueue:    []CharacterAction{},
		Inventory:       NewInventory(),
		abilities:       abilities,
		lastUsed:        lastUsed,
	}
}

func (p Character) GetId() string {
	return p.id
}

func (p *Character) SetPosition(position utils.Vector2) {
	p.position.Teleport(position)
}

func (p Character) GetPosition() utils.Vector2 {
	return p.position
}

// this shouldn't be here
func (p Character) GetHealth() stats.Health {
	return p.Health
}

func (p Character) GetExecutingAction() ExecutingAction {
	return p.executingAction
}

func (p Character) GetAbilities() map[string]*Ability {
	return p.abilities
}

func (character *Character) EnqueueAction(action CharacterAction) {
	if character.IsAlive() {
		character.actionsQueue = append(character.actionsQueue, action)
	}
}

func (p *Character) ClearActionsQueue() {
	p.actionsQueue = nil
}

func (p *Character) PrependAction(action CharacterAction) {
	p.actionsQueue = append([]CharacterAction{action}, p.actionsQueue...)
}

func (c *Character) ExecuteNextAction() {
	if len(c.actionsQueue) == 0 {
		if c.executingAction.actionType != Idle {
			c.executingAction = ExecutingAction{Idle, c.executingAction.direction}
		}
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
	normalizedDirection := utils.Normalize(p.position, p.to)
	p.direction = normalizedDirection.Scale(CHARACTER_SPEED)
	// Wirty wirty
	if p.position == p.to {
		p.executingAction = ExecutingAction{Moving, p.executingAction.direction}
	} else {
		p.executingAction = ExecutingAction{Moving, normalizedDirection}
	}
}

func (p Character) IsMoving() bool {
	return !p.position.Equals(p.to)
}

// returns if the character has been able to complete the movement without colliding with any obstacle
func (p *Character) UpdatePosition(deltaTime float64, obstacles [][]utils.Vector2) bool {
	distanceToTarget := p.position.Distance(p.to)
	if distanceToTarget < (CHARACTER_SPEED * deltaTime) {
		p.position.Teleport(p.to)
	} else {
		nextPosition := p.position.Add(p.direction.Scale(deltaTime))

		// Define the character's circle center and radius
		circleCenter := nextPosition
		circleRadius := CHARACTER_BOUNDS_RADIUS

		for _, obstacle := range obstacles {
			if utils.CirclePolygonIntersect(circleCenter, circleRadius, obstacle) {
				// Should this be null and IsMoving check for null?
				p.to = p.position
				return false
			}
		}
		// No collision detected, update position
		p.position = nextPosition
	}
	return true
}

func (p Character) GetRadius() float64 {
	return CHARACTER_BOUNDS_RADIUS
}

func (player *Character) IsInCooldown(abilityId string) bool {
	ability := player.abilities[abilityId]
	if player.RemainingCooldown(ability) <= 0 {
		return false
	} else {
		return true
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

func (player *Character) CastAbility(ability Ability, params AbilityParameters) {
	player.lastUsed[ability.id] = time.Now()
	targetCoordinates, _ := params.GetTargetCoordinates()
	playerPosition := player.position
	if playerPosition != targetCoordinates {
		normalizedCastAbilityVector := utils.Normalize(playerPosition, targetCoordinates)
		player.executingAction = ExecutingAction{ability.characterState, normalizedCastAbilityVector}
	} else {
		player.executingAction = ExecutingAction{ability.characterState, player.executingAction.direction}
	}
	ability.Cast(*player, params)
}

// removal => disconnection & death, for now, but for players nothing subscribes to removal.
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

// Everything below this should be in a functionality module, not in the entities package?

func (character *Character) EnqueueMovementAction(position utils.Vector2) {
	moveAction := &MoveAction{
		TargetPosition: position,
	}
	character.ClearActionsQueue()
	character.EnqueueAction(moveAction)
}

func (character *Character) UseItem(itemId string) {
	if !character.Inventory.CanConsume(itemId) {
		err := fmt.Errorf("character %s tried to use item %s, but is unable to", character.id, itemId)
		fmt.Println(err.Error())
		return
	}

	item := character.GetItem(itemId)
	item.ExecuteMechanics(character)
	character.Inventory.AddItem(item, -1)
}
