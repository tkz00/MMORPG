package entities

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/net/websocket"

	"tkz00/backend/pkg/utils"
)

type GameState struct {
	playerIds   map[*websocket.Conn]string
	players     map[string]*Character
	projectiles map[string]*Projectile
	spawners    map[string]*Spawner
	npcs        map[string]*Npc
}

func StartGameState() GameState {
	gs := GameState{
		playerIds:   make(map[*websocket.Conn]string),
		players:     make(map[string]*Character),
		projectiles: make(map[string]*Projectile),
		spawners:    make(map[string]*Spawner),
		npcs:        make(map[string]*Npc),
	}

	return gs
}

func (gs *GameState) SetUpSkeletonEnemies(skeletonEnemiesAbilities map[string]*Ability) {
	skeletonNPCTemplate := NewNPCTemplate("0", "skeleton", 25, 12, skeletonEnemiesAbilities)

	gs.spawners["skeleton_spawner_0"] = NewSpawner(
		*utils.NewVector2(0, 15),
		2.5,
		4*time.Second,
		skeletonNPCTemplate,
	)
}

func (gs *GameState) AddPlayer(conn *websocket.Conn, player *Character) {
	playerId := player.GetId()
	gs.playerIds[conn] = playerId
	gs.players[playerId] = player
}

func (gs *GameState) DeletePlayer(conn *websocket.Conn) {
	playerId := gs.playerIds[conn]
	gs.players[playerId].Remove()
	delete(gs.players, playerId)
	delete(gs.playerIds, conn)
}

func (gs GameState) GetPlayerCount() int {
	return len(gs.players)
}

func (gs GameState) Players() map[string]*Character {
	return gs.players
}

func (gs GameState) GetPlayers() []Character {
	// WTF is this shit scoob????
	playersSlice := make([]Character, 0, len(gs.players))
	for _, player := range gs.players {
		playersSlice = append(playersSlice, *player)
	}
	return playersSlice
}

// I don't know if this is OK, probably not, probably native server should hold the relationship between connections and character/player ids
func (gs GameState) GetPlayerByConn(conn *websocket.Conn) *Character {
	playerId := gs.playerIds[conn]
	return gs.players[playerId]
}

func (gs *GameState) AddProjectile(projectile *Projectile) {
	gs.projectiles[projectile.GetId()] = projectile
}

func (gs *GameState) Projectiles() map[string]*Projectile {
	return gs.projectiles
}

func (gs GameState) GetProjectiles() []Projectile {
	projectilesSlice := make([]Projectile, 0, len(gs.projectiles))
	for _, projectile := range gs.projectiles {
		projectilesSlice = append(projectilesSlice, *projectile)
	}
	return projectilesSlice
}

func (gs *GameState) RemoveProjectile(key string) {
	delete(gs.projectiles, key)
}

func (gs GameState) GetNPCs() []Npc {
	npcsSlice := make([]Npc, 0, len(gs.npcs))
	for _, player := range gs.npcs {
		npcsSlice = append(npcsSlice, *player)
	}
	return npcsSlice
}

func (gs GameState) NPCs() map[string]*Npc {
	return gs.npcs
}

func (gs GameState) GetCharacterById(targetId string) (*Character, error) {
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

func (gs GameState) Spawners() map[string]*Spawner {
	return gs.spawners
}

// Everything below this should be in a functionality module, not in the entities package

func EnqueueAbilityCast(gs GameState, caster *Character, abilityInfo AbilityInfo) {
	abilityId := abilityInfo.GetId()
	if caster.IsInCooldown(abilityId) {
		return
	}

	ability := caster.GetAbilities()[abilityId]
	targetCoordinatesCallback := func(targetId string) (utils.Vector2, error) {
		target, err := gs.GetCharacterById(targetId)
		if err != nil {
			fmt.Println(err)
			return utils.Vector2{}, err
		}
		return target.GetPosition(), nil
	}
	castingCoordinatesCallback := func(targetId string) (utils.Vector2, error) {
		target, err := gs.GetCharacterById(targetId)
		if err != nil {
			fmt.Println(err)
			return utils.Vector2{}, err
		}

		if caster.GetPosition().Distance(target.GetPosition()) > ability.Range() {
			return utils.ClosestPositionInRange(caster.position, target.GetPosition(), ability.Range()), nil
		} else {
			return caster.GetPosition(), nil
		}
	}
	abilityAction := ability.CreateAction(abilityInfo, targetCoordinatesCallback, castingCoordinatesCallback)
	caster.ClearActionsQueue()
	caster.EnqueueAction(&abilityAction)
}
