package entities

import (
	"backend/pkg/utils"
	"fmt"
	"time"
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
	aggroRange        float64
}

func (npc *Npc) UpdateBehaviour() {
	hasNoActionsInQueue := len(npc.Character.actionsQueue) == 0
	if hasNoActionsInQueue {
		switch npc.state {
		case Pacific:
			npc.takePacificAction()
		case Aggressive:
			if npc.target.IsAlive() && npc.targetIsInRange(npc.target.position) {
				npc.takeAggressiveAction()
			} else {
				npc.BecomePacific()
			}
		}
	}
}

func (npc *Npc) takePacificAction() {
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

func (npc *Npc) takeAggressiveAction() {
	ability := npc.abilities["0"]
	if npc.IsInCooldown(ability.id) {
		const epsilon = 1e-9
		if (npc.position.Distance(npc.target.position) - ability.rangeValue) > epsilon {
			targetPosition := utils.ClosestPositionInRange(
				npc.position,
				npc.target.GetPosition(),
				(ability.Range() - 0.01),
			)
			moveAction := &MoveAction{
				TargetPosition: targetPosition,
			}
			npc.EnqueueAction(moveAction)
		}
		return
	}

	castParameters := make(map[Targeting]interface{})

	switch ability.targeting {
	case Target:
		castParameters[Target] = npc.target.id
	case Coordinates:
		fmt.Println("No npc logic for coordinate targetted abilities")
	}

	npc.EnqueueAbilityCastAction(ability.id, castParameters)
}

func (npc *Npc) BecomeAggressive(target *Character) {
	npc.target = target
	npc.state = Aggressive
	target.SubscribeToRemoval(func() {
		npc.BecomePacific()
	})
}

func (npc *Npc) BecomePacific() {
	npc.target = nil
	npc.state = Pacific
}

func (npc *Npc) targetIsInRange(targetPosition utils.Vector2) bool {
	return npc.spawnerPosition.Distance(targetPosition) < npc.aggroRange
}
