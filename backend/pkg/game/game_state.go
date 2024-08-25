package game

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"

	"unnamed-mmo/backend/pkg/game/character"
	"unnamed-mmo/backend/pkg/game/spawner"
	"unnamed-mmo/backend/pkg/utils"
)

type GameState struct {
	playerIds   map[*websocket.Conn]string
	players     map[string]*character.Character
	projectiles map[string]*Projectile
	spawners    map[string]*spawner.Spawner
	npcs        map[string]*character.Npc
}

func StartGameState() GameState {
	gs := GameState{
		playerIds:   make(map[*websocket.Conn]string),
		players:     make(map[string]*character.Character),
		projectiles: make(map[string]*Projectile),
		spawners:    make(map[string]*spawner.Spawner),
		npcs:        make(map[string]*character.Npc),
	}

	skeletonEnemiesAbilities := map[string]*character.Ability{
		"0": character.NewAbility("0", "health burn", 7, 1500, character.Target, character.Attacking,
			func(caster character.Character, params character.AbilityParameters) {
				targetId := params.(character.TargetIdAbilityParams).TargetId

				target, err := gs.getCharacterById(targetId)
				if err != nil {
					fmt.Println(err)
					return
				}

				target.HealthVariation(-10)
			}),
	}

	skeletonNPCTemplate := character.NewNPCTemplate("0", "skeleton", 25, skeletonEnemiesAbilities)

	gs.spawners["skeleton_spawner_0"] = spawner.NewSpawner(
		*utils.NewVector2(0, 15),
		2.5,
		4*time.Second,
		skeletonNPCTemplate,
	)

	return gs
}

func (gs *GameState) AddPlayer(conn *websocket.Conn) character.Character {
	id := uuid.New()
	playerId := id.String()
	abilities := map[string]*character.Ability{
		"0": character.NewAbility("0", "projectile", 5, 2000, character.Coordinates, character.Attacking,
			func(caster character.Character, params character.AbilityParameters) {
				projectileId := uuid.NewString()
				gs.projectiles[projectileId] = CreateProjectile(uuid.NewString(), caster.GetPosition(), params.(character.CoordinateAbilityParams).Target, 5, caster.GetId())
			}),
		"1": character.NewAbility("1", "heal", 7, 3000, character.Target, character.CastingHeal,
			func(caster character.Character, params character.AbilityParameters) {
				targetId := params.(character.TargetIdAbilityParams).TargetId

				target, err := gs.getCharacterById(targetId)
				if err != nil {
					fmt.Println(err)
					return
				}

				target.HealthVariation(10)
			}),
	}
	player := character.CreateCharacter(playerId, 0, 0, abilities)
	gs.playerIds[conn] = playerId
	gs.players[playerId] = player
	return *player
}

func (gs *GameState) DeletePlayer(conn *websocket.Conn) {
	playerId := gs.playerIds[conn]
	delete(gs.players, playerId)
	delete(gs.playerIds, conn)
}

func (gs GameState) GetPlayerCount() int {
	return len(gs.players)
}

func (gs GameState) GetPlayers() []character.Character {
	// WTF is this shit scoob????
	playersSlice := make([]character.Character, 0, len(gs.players))
	for _, player := range gs.players {
		playersSlice = append(playersSlice, *player)
	}
	return playersSlice
}

func (gs GameState) GetProjectiles() []Projectile {
	projectilesSlice := make([]Projectile, 0, len(gs.projectiles))
	for _, projectile := range gs.projectiles {
		projectilesSlice = append(projectilesSlice, *projectile)
	}
	return projectilesSlice
}

func (gs GameState) GetNPCs() []character.Npc {
	npcsSlice := make([]character.Npc, 0, len(gs.npcs))
	for _, player := range gs.npcs {
		npcsSlice = append(npcsSlice, *player)
	}
	return npcsSlice
}

func (gs GameState) MovePlayer(conn *websocket.Conn, position utils.Vector2) {
	playerId := gs.playerIds[conn]
	moveAction := &character.MoveAction{
		TargetPosition: position,
	}
	gs.players[playerId].ClearActionsQueue()
	gs.players[playerId].EnqueueAction(moveAction)
}

func (gs GameState) UpdateState(deltaTime float64) {
	gs.updatePlayers(deltaTime)
	gs.updateProjectiles(deltaTime)
	gs.updateSpawners()
	gs.updateNpcs(deltaTime)
}

func (gs *GameState) updatePlayers(deltaTime float64) {
	for _, player := range gs.players {
		player.ExecuteNextAction()
		if player.IsMoving() {
			player.UpdatePosition(deltaTime)
		}
	}
}

func (gs GameState) updateProjectiles(deltaTime float64) {
	for key, projectile := range gs.projectiles {
		if projectile.state == Hit {
			delete(gs.projectiles, key)
			continue
		}

		if projectile.state == Active && gs.checkCollision(*projectile) {
			projectile.state = Hit
			continue
		}

		isAtMaxRange := projectile.UpdatePosition(deltaTime)

		if isAtMaxRange {
			delete(gs.projectiles, key)
		}
	}
}

func (gs GameState) checkCollision(projectile Projectile) bool {
	projectileHit := false
	for _, player := range gs.players {
		if player.GetId() != projectile.caster && gs.AreColliding(*player, projectile) {
			player.HealthVariation(-projectile.damage)
			projectileHit = true
		}
	}

	for _, npc := range gs.npcs {
		if npc.GetId() != projectile.caster && gs.AreColliding(*npc.Character, projectile) {
			npc.HealthVariation(-projectile.damage)
			if caster, err := gs.getCharacterById(projectile.caster); err == nil {
				npc.BecomeAggressive(caster)
			}
			projectileHit = true
		}
	}

	return projectileHit
}

func (gs *GameState) EnqueueAbilityCast(conn *websocket.Conn, abilityInfo character.AbilityInfo) {
	casterId := gs.playerIds[conn]
	caster := gs.players[casterId]
	ability := caster.GetAbilities()[abilityInfo.GetId()]
	abilityAction := ability.CreateAction(abilityInfo,
		func(targetId string) (utils.Vector2, error) {

			target, err := gs.getCharacterById(targetId)
			if err != nil {
				fmt.Println(err)
				return utils.Vector2{}, err
			}

			if caster.GetPosition().Distance(target.GetPosition()) > ability.Range() {
				return caster.ClosestPositionInRange(target.GetPosition(), ability.Range()), nil
			} else {
				return caster.GetPosition(), nil
			}
		})
	caster.EnqueueAbilityCast(abilityAction)
}

func (gs GameState) AreColliding(player character.Character, projectile Projectile) bool {
	playerPosition := player.GetPosition()
	projectilePosition := projectile.GetPosition()
	distance := playerPosition.Distance(projectilePosition)
	return distance < (player.GetRadius() + projectile.GetRadius())
}

func (gs *GameState) updateSpawners() {
	for _, spawner := range gs.spawners {
		newNPCs := spawner.GetNewNPCs()
		for _, newNPC := range newNPCs {
			gs.npcs[newNPC.GetId()] = newNPC
		}
	}
}

func (gs *GameState) updateNpcs(deltaTime float64) {
	for _, npc := range gs.npcs {
		npc.Update(deltaTime)
	}
}

func (gs GameState) getCharacterById(targetId string) (*character.Character, error) {
	target, exists := gs.players[targetId]
	if !exists {
		targetNpc, npcExists := gs.npcs[targetId]
		if !npcExists {
			fmt.Println("ERROR: no target found with id")
			return nil, errors.New("character not found")
		}
		target = targetNpc.Character
	}
	return target, nil
}
