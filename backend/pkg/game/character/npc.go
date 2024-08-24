package character

import (
	"time"
	"unnamed-mmo/backend/pkg/utils"
)

type Npc struct {
	*Character
	lastActionDecided	time.Time
	spawnerPosition		utils.Vector2
}

func (npc *Npc) Update(deltaTime float64) {
	hasNoActionsInQueue := len(npc.Character.actionsQueue) == 0
	if hasNoActionsInQueue {
		timeSinceLastDecision := time.Since(npc.lastActionDecided).Milliseconds()
		if timeSinceLastDecision > 5000 {
			targetPosition := utils.RandomCoordinatesInRadius(npc.spawnerPosition, 25)
			moveAction := &MoveAction{
				TargetPosition: targetPosition,
			}
			npc.EnqueueAction(moveAction)
			npc.lastActionDecided = time.Now()
		}
	}
	npc.Character.ExecuteNextAction()
	if(npc.IsMoving()) {
		npc.UpdatePosition(deltaTime)
	}
}
