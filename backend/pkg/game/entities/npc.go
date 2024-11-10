package entities

import (
	"time"
	"tkz00/backend/pkg/utils"
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
				npc.TakeAggressiveAction()
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

func (npc *Npc) TakeAggressiveAction() {

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
