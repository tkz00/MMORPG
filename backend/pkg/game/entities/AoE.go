package entities

import (
	"backend/pkg/utils"
	"slices"

	"github.com/google/uuid"
)

type AoE struct {
	id                  string
	caster              string
	position            utils.Vector2 // For simplicity AoE will just be circular shape, in future it should also be at least polygonal don't know
	radius              float64       // if the best would be some kind of abstraction or just have CircularAoE and PolygonalAoE as different types
	remainingDurationMs int
	onHitMechanics      []Mechanic
	hitCharacterIds     []string
}

func InstantiateAoE(
	position utils.Vector2,
	radius float64,
	durationMs int,
	casterId string,
	onHitMechanics []Mechanic,
) *AoE {
	return &AoE{
		id:                  uuid.New().String(),
		caster:              casterId,
		position:            position,
		radius:              radius,
		remainingDurationMs: durationMs,
		onHitMechanics:      onHitMechanics,
	}
}

func (AoE AoE) Id() string {
	return AoE.id
}

func (AoE AoE) CasterId() string {
	return AoE.caster
}

func (AoE AoE) Position() utils.Vector2 {
	return AoE.position
}

func (AoE AoE) Radius() float64 {
	return AoE.radius
}

func (AoE *AoE) Tick(gs *GameState, deltaTimeS float64) {
	for _, player := range gs.players {
		if player.id != AoE.caster && areColliding(player, AoE) {
			if !slices.Contains(AoE.hitCharacterIds, player.id) {
				for _, mechanic := range AoE.onHitMechanics {
					mechanic.Params["target_id"] = player.id
				}
				resolveMechanics(AoE.caster, gs, AoE.onHitMechanics)
				AoE.hitCharacterIds = append(AoE.hitCharacterIds, player.id)
			}
		}
	}

	for _, npc := range gs.npcs {
		if npc.id != AoE.caster && areColliding(npc.Character, AoE) {
			if !slices.Contains(AoE.hitCharacterIds, npc.id) {
				for _, mechanic := range AoE.onHitMechanics {
					mechanic.Params["target_id"] = npc.id
				}
				resolveMechanics(AoE.caster, gs, AoE.onHitMechanics)
				AoE.hitCharacterIds = append(AoE.hitCharacterIds, npc.id)
			}
		}
	}

	AoE.remainingDurationMs -= int(deltaTimeS * 1000)
}

func (AoE *AoE) RemainingDurationMs() int {
	return AoE.remainingDurationMs
}

func areColliding(character *Character, AoE *AoE) bool {
	characterPosition := character.GetPosition()
	AoEPosition := AoE.position
	distance := characterPosition.Distance(AoEPosition)
	return distance < (character.GetRadius() + AoE.radius)
}
