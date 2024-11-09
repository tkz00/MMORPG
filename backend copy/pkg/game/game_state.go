package game

import "tkz00/backend/pkg/utils"

type GameState struct {
	players map[string]Character
}

func NewGameState() GameState {
	return GameState{
		players: make(map[string]Character),
	}
}

func (gs GameState) AddPlayer(id string) GameState {
	gs.players[id] = Character{position: utils.Vector2{}}
	return gs
}

func (gs GameState) RemovePlayer(id string) GameState {
	delete(gs.players, id)
	return gs
}

func (gs GameState) Players() map[string]Character {
	return gs.players
}

func (gs GameState) MovePlayer(id string, position utils.Vector2) GameState {
	player := gs.players[id]
	player.position = position
	gs.players[id] = player
	return gs
}
