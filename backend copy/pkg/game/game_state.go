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

type Character struct {
	position utils.Vector2
}
