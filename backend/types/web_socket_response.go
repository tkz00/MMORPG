package types

type WebSocketResponse struct {
	PlayerID     string   `json:"PlayerID,omitempty"`
	GameStateDTO *GameDTO `json:"GameStateDTO,omitempty"`
}
