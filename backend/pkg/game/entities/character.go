package entities

import (
	"backend/pkg/utils"
	"fmt"
	"math"
	"time"
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
	actionType               Action
	direction                utils.Vector2
	durationInRemainingTicks *int64
}

func (action ExecutingAction) ActionType() Action {
	return action.actionType
}

func (action ExecutingAction) Direction() utils.Vector2 {
	return action.direction
}

func (action ExecutingAction) DurationInRemainingTicks() *int64 {
	return action.durationInRemainingTicks
}

type CharacterLastTickState struct {
	Position                   *utils.Vector2
	Radius                     *float64
	MaxHealth                  *int
	CurrentHealth              *int
	Stats                      map[string]int64
	Action                     *Action
	Direction                  *utils.Vector2
	Abilities                  []string
	AbilitiesRemainingCooldows map[string]int64
	Items                      map[string]int64
}

type Character struct {
	id string
	Health
	stats           map[string]int64
	executingAction ExecutingAction
	position        utils.Vector2
	to              utils.Vector2
	direction       utils.Vector2
	actionsQueue    []CharacterAction

	*Inventory

	CharacterLastTickState *CharacterLastTickState

	// should this be here?
	abilities map[string]*Ability
	lastUsed  map[string]time.Time

	onRemoved []func()
}

func CreateCharacter(
	id string,
	x, z float64,
	stats map[string]int64,
	abilities map[string]*Ability,
) *Character {
	initialPosition := *utils.NewVector2(x, z)

	lastUsed := make(map[string]time.Time, len(abilities))
	for _, ability := range abilities {
		lastUsed[ability.id] = time.Time{}
	}

	return &Character{
		id:              id,
		position:        initialPosition,
		to:              initialPosition,
		Health:          NewHealth(BASE_MAX_HEALTH),
		stats:           stats,
		executingAction: ExecutingAction{Idle, *utils.NewVector2(0, 0), nil},
		actionsQueue:    []CharacterAction{},
		Inventory:       NewInventory(),
		abilities:       abilities,
		lastUsed:        lastUsed,
		CharacterLastTickState: &CharacterLastTickState{
			AbilitiesRemainingCooldows: make(map[string]int64),
			Items:                      make(map[string]int64),
			Stats:                      make(map[string]int64),
		},
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
func (p Character) GetHealth() Health {
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

func (c *Character) ExecuteNextAction(gs *GameState) {
	if c.executingAction.durationInRemainingTicks != nil {
		*c.executingAction.durationInRemainingTicks--
		if *c.executingAction.durationInRemainingTicks > 0 {
			return
		}
	}

	if len(c.actionsQueue) == 0 {
		if c.executingAction.actionType != Idle {
			c.executingAction = ExecutingAction{Idle, c.executingAction.direction, nil}
		}
		return
	}

	currentAction := c.actionsQueue[0]
	if !currentAction.IsExecuted() {
		err := currentAction.Execute(c, gs)
		if err != nil {
			fmt.Println("Error executing action:", err)
			return
		}
	}

	if currentAction.IsComplete(c) {
		c.actionsQueue = c.actionsQueue[1:]
	}
}

func (p *Character) MoveTowards(to utils.Vector2) {
	p.to = to
	normalizedDirection := utils.Normalize(p.position, p.to)
	p.direction = normalizedDirection.Scale(CHARACTER_SPEED)
	// Wirty wirty
	if p.position == p.to {
		p.executingAction = ExecutingAction{Idle, p.executingAction.direction, nil}
	} else {
		p.executingAction = ExecutingAction{Moving, normalizedDirection, nil}
	}
}

func (p Character) IsMoving() bool {
	return !p.position.Equals(p.to)
}

func (p *Character) StopMovement() {
	p.to = p.position
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
	remainingCooldown := player.abilities[ability.id].cooldown - now.Sub(player.lastUsed[ability.id]).
		Milliseconds()
	if remainingCooldown > 0 {
		return remainingCooldown
	}
	return 0
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

// Everything below this should be in a functionality module, not in the entities package?

func (character *Character) EnqueueMovementAction(position utils.Vector2) {
	moveAction := &MoveAction{
		TargetPosition: position,
	}
	character.ClearActionsQueue()
	character.EnqueueAction(moveAction)
}

func (character *Character) UseItem(itemId string, targetId string, gs *GameState) {
	if !character.IsAlive() {
		return
	}
	if !character.Inventory.CanConsume(itemId) {
		err := fmt.Errorf(
			"character %s tried to use item %s, but is unable to",
			character.id,
			itemId,
		)
		fmt.Println(err.Error())
		return
	}

	item := GetItem(itemId)
	item.ExecuteMechanics(character, targetId, gs)
	character.Inventory.AddItem(itemId, -1)
}

func (caster *Character) EnqueueAbilityCastAction(
	abilityId string,
	abilityCastParameters map[Targeting]interface{},
) {
	if caster.IsInCooldown(abilityId) {
		return
	}

	ability := caster.abilities[abilityId]
	caster.ClearActionsQueue()
	castAbilityAction := &AbilityCastAction{
		ability:        *ability.Clone(),
		castParameters: abilityCastParameters,
	}
	caster.EnqueueAction(castAbilityAction)
}

func (c *Character) TakeDamage(d int) {
	damage := d - int(c.stats["defense"])
	c.HealthVariation(-int(math.Max(0, float64(damage))))
}

func (c *Character) Heal(a int) {
	c.HealthVariation(a)
}

func (c *Character) GetStat(s string) int64 {
	return c.stats[s]
}

func (c *Character) GetStats() map[string]int64 {
	return c.stats
}

func (c *Character) SetStat(stat string, value int64) {
	c.stats[stat] = value
}
