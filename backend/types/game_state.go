
package types
import(
	"golang.org/x/net/websocket"
)
type GameState struct {
	players map[*websocket.Conn]Player
}

func StartGameState() GameState{
	return GameState{
		players: make(map[*websocket.Conn]Player),
	}
}

func (gs *GameState) AddPlayer(conn *websocket.Conn){
	var player Player
	gs.players[conn] = player.CreatePlayer(20.0,10.0)
}

func (gs GameState) GetPlayerCount() int{
	return len(gs.players)
}