package character

import (
	"time"
	"unnamed-mmo/backend/pkg/utils"
)

type NpcState int

const (
	Pacific NpcState = iota
	Aggressive
)

type Npc struct {
	*Character
	lastActionDecided time.Time
	spawnerPosition   utils.Vector2
	state             NpcState
	target            *Character
}

func (npc *Npc) Update(deltaTime float64) {
	hasNoActionsInQueue := len(npc.Character.actionsQueue) == 0
	if hasNoActionsInQueue {
		switch npc.state {
		case Pacific:
			npc.TakePacificAction()
		case Aggressive:
			npc.TakeAggressiveAction()
		}
	}
	npc.Character.ExecuteNextAction()
	if npc.IsMoving() {
		npc.UpdatePosition(deltaTime)
	}
}

func (npc *Npc) TakePacificAction() {
	timeSinceLastDecision := time.Since(npc.lastActionDecided).Milliseconds()
	if timeSinceLastDecision > 5000 {
		targetPosition := utils.RandomCoordinatesInRadius(npc.spawnerPosition, 10)
		moveAction := &MoveAction{
			TargetPosition: targetPosition,
		}
		npc.EnqueueAction(moveAction)
		npc.lastActionDecided = time.Now()
	}
}

func (npc *Npc) TakeAggressiveAction() {
	// var ability Ability
	// for _, value := range npc.abilities {
	// 	ability = *value
	// 	break
	// }

	ability := npc.abilities["0"]

	targetPositionCallback := func(targetId string) (utils.Vector2, error) {
		if npc.GetPosition().Distance(npc.target.GetPosition()) > ability.Range() {
			return npc.ClosestPositionInRange(npc.target.GetPosition(), ability.Range()), nil
		} else {
			return npc.GetPosition(), nil
		}
	}

	var abilityParams AbilityParameters
	switch ability.targeting {
	case Target:
		targetId := npc.target.id
		abilityParams = TargetIdAbilityParams{
			TargetId:               targetId,
			TargetPositionCallback: targetPositionCallback,
		}
	case Coordinates:
		targetPosition := npc.target.position
		abilityParams = CoordinateAbilityParams{
			Target: targetPosition,
		}
	}
	npc.EnqueueAbilityCast(
		CastAbilityAction{
			ability: *ability,
			params:  abilityParams,
		},
	)
}

func (npc *Npc) BecomeAggressive(target *Character) {
	npc.target = target
	npc.state = Aggressive
}
