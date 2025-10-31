package entities

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"

	"backend/pkg/utils"
)

type GameState struct {
	playerIds             map[*websocket.Conn]string
	players               map[string]*Character
	projectiles           map[string]*Projectile
	areaEffects           map[string]*AoE
	spawners              map[string]*Spawner
	npcs                  map[string]*Npc
	obstacles             [][]utils.Vector2
	runningMechanics      map[string]*Mechanic
	mechanicRemainingTime map[string]time.Duration
}

func StartGameState() *GameState {
	gs := &GameState{
		playerIds:             make(map[*websocket.Conn]string),
		players:               make(map[string]*Character),
		projectiles:           make(map[string]*Projectile),
		areaEffects:           make(map[string]*AoE),
		spawners:              make(map[string]*Spawner),
		npcs:                  make(map[string]*Npc),
		runningMechanics:      make(map[string]*Mechanic),
		mechanicRemainingTime: make(map[string]time.Duration),
	}

	RegisterMechanicHandler("heal", HealMechanic)
	RegisterMechanicHandler("damage", DamageMechanic)
	RegisterMechanicHandler("create_projectile", CreateProjectileMechanic)
	RegisterMechanicHandler("delay", DelayMechanic)
	RegisterMechanicHandler("create_AoE", AoEMechanic)
	RegisterMechanicHandler("buff_stat", BuffStatMechanic)

	return gs
}

func (gs *GameState) SetUpSkeletonEnemies(skeletonEnemiesAbilities map[string]*Ability) {
	stats := map[string]int64{"damage": 10, "defense": 10}
	skeletonNPCTemplate := NewNPCTemplate("0", "skeleton", 50, stats, 12, skeletonEnemiesAbilities)

	gs.spawners["skeleton_spawner_0"] = NewSpawner(
		*utils.NewVector2(0, 15),
		2.5,
		4*time.Second,
		skeletonNPCTemplate,
	)
}

func (gs *GameState) SetUpObstacles(obstacleColliders [][]utils.Vector2) {
	gs.obstacles = obstacleColliders
}

func (gs *GameState) AddPlayer(conn *websocket.Conn, player *Character) {
	playerId := player.GetId()
	gs.playerIds[conn] = playerId
	gs.players[playerId] = player
	player.AddEffect(
		Effect{
			id:      "spell_vampirism",
			name:    "Spell Vampirism",
			trigger: EffectOnDamageDealt,
			mechanics: []Mechanic{
				{
					MechanicType: "heal",
					Params: map[string]interface{}{
						"targeting_strategy": "caster",
						"base_amount":        10.0,
					},
				},
			},
			parameters: map[string]interface{}{
				"base_amount": 10.0,
				// "scaling_from_event": map[string]interface{}{
				// 	"event_field": "damage",
				// 	"factor":      0.1,
				// },
			},
		},
	)
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

func (gs *GameState) AreaEffects() map[string]*AoE {
	return gs.areaEffects
}

func (gs *GameState) AddAreaEffect(AoE *AoE) {
	gs.areaEffects[AoE.id] = AoE
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

func (gs GameState) GetObstacleColliders() [][]utils.Vector2 {
	return gs.obstacles
}

func (gs *GameState) DelayMechanics(delayMechanics []Mechanic, delayMs int, casterId string) {
	for _, mechanic := range delayMechanics {
		mechanic.Params["caster_id"] = casterId
		mechanicId := uuid.New().String()
		gs.runningMechanics[mechanicId] = &mechanic
		gs.mechanicRemainingTime[mechanicId] = time.Duration(delayMs) * time.Millisecond
	}
}

func (gs *GameState) UpdateMechanics(deltaTime float64) {
	var completedMechanicIds []string
	for mechanicId, mechanic := range gs.runningMechanics {
		gs.mechanicRemainingTime[mechanicId] -= time.Duration(deltaTime * float64(time.Second))
		if gs.mechanicRemainingTime[mechanicId] <= 0 {
			resolveMechanics(mechanic.Params["caster_id"].(string), gs, []Mechanic{*mechanic})
			completedMechanicIds = append(completedMechanicIds, mechanicId)
		}
	}

	for _, mechanicToDelete := range completedMechanicIds {
		delete(gs.mechanicRemainingTime, mechanicToDelete)
		delete(gs.runningMechanics, mechanicToDelete)
	}
}
